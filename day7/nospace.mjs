#!/usr/bin/env node

import { createInterface } from 'readline';
import { createReadStream } from 'fs';
import { syncBuiltinESMExports } from 'module';

if ( process.argv.length != 3 ) {
    console.error("Provide input file");
    process.exit(-1);
}

async function parseInput(path) {
    const rl = createInterface({
        input: createReadStream(path)
    });

    const filesys = { files: [], dirs: {} };
    let curpath;
    let curdir;
    let lsoutput = false;

    rl.on('error', (err) => {
        if ( err.code === 'ENOENT' ) {
            console.error(`File ${path} does not exist`)
            process.exit(1);
        } else {
            throw err;
        }
    })

    for await (const line of rl) {
        if ( line[0] == '$' ) {
            // handle commands
            lsoutput = false
            let cmdline = line.substring(1).trim().split(" ")
            switch ( cmdline[0] ) {
                case "cd":
                    if ( cmdline[1][0] == '/' ) {
                        curpath = cmdline[1]
                    } else if ( cmdline[1] == ".." ) {
                        if ( ! curpath ) {
                            console.error("Path not set yet, cannot cd ..")
                            process.exit(1)
                        }
                        let lastslash = curpath.lastIndexOf("/")
                        if ( lastslash <= 0 ) {
                            console.log("Warning: doing cd .. in root dir")
                        } else {
                            curpath = curpath.substring(0, lastslash)
                        }
                    } else {
                        curpath = curpath + "/" + cmdline[1]
                    }
                    // now walk to the specified directory, creating subdirs where necessary
                    curdir = filesys
                    for (const elem of curpath.split("/")) {
                        if ( elem ) {
                            if ( ! curdir.dirs[elem] ) {
                                curdir.dirs[elem] = { files: [], dirs: {} }
                            }
                            curdir = curdir.dirs[elem]
                        }
                    }
                    break
                case "ls":
                    lsoutput = true
                    break
                default:
                    console.error(`Unknown command in input ${cmdline[0]}`)
                    process.exit(3);
            }
        } else if ( lsoutput ) {
            // handle output of the "ls" command
            let spacepos = line.indexOf(" ")
            let size = line.substring(0, spacepos)
            let name = line.substring(spacepos + 1).trim()
            if ( size == "dir" ) {
                if ( curdir.dirs[name] ) {
                    console.log("Warning: duplicate directory %s found", name)
                } else {
                    curdir.dirs[name] = { files: [], dirs: {} }
                }
            } else {
                let intsize = Number.parseInt(size)
                if ( Number.isNaN(intsize) ) {
                    console.error(`Error parsing ls output, ${size} not dir or integer size`)
                    process.exit(4);
                }
                curdir.files.push( {name: name, size: intsize} )
            }
        } else {
            console.error(`Unexpected input line ${line}`)
            process.exit(5);
        }
    }
    return filesys
}

function walkDirs (fs, cb) {
    walkDirsRecursive("/", fs, cb)
}

function walkDirsRecursive (path, fs, cb) {
    for (let subdir in fs.dirs) {
        walkDirsRecursive(path + subdir + "/", fs.dirs[subdir], cb)
    }
    cb(path, fs)
}

function setdirsize(path, fs) {
    let totalsize = 0
    // walk our files
    for (let direntry of fs.files ) {
        totalsize += direntry.size
    }
    // console.log("Size just of files in directory %s is %d", path, totalsize)
    // add sizes of subdirectories
    for (let subdir in fs.dirs) {
        totalsize += fs.dirs[subdir].size
    }
    fs.size = totalsize
}

async function main() {
    const startTime = process.hrtime()

    const fs = await parseInput(process.argv[2])
    const parseTook = process.hrtime(startTime)
    console.log(`Parsing took ${parseTook[0] + parseTook[1] / 1e9} seconds`)
    walkDirs(fs, setdirsize)

    // part 1
    let sumsmall = 0
    walkDirs(fs, (path, dir) => {
        if ( dir.size <= 100000 ) {
            sumsmall += dir.size
        }
    })
    const part1Took = process.hrtime(startTime)
    console.log("Sum of small directories is %d", sumsmall)
    console.log(`Parsing and calculating part 1 took: ${part1Took[0] + part1Took[1] / 1e9} seconds`)

    // part 2
    const totalSize = 70000000
    const sizeNeeded = 30000000
    const nowFree = totalSize - fs.size
    if ( nowFree > sizeNeeded ) {
        console.log("We have %d bytes free, which is enough already", nowFree)
    } else {
        let bestMatch = Number.MAX_SAFE_INTEGER
        const toDelete = sizeNeeded - nowFree
        walkDirs(fs, (path, dir) => {
            if ( dir.size >= toDelete && dir.size < bestMatch ) {
                bestMatch = dir.size
            }
        })
        const part2Took = process.hrtime(startTime)
        console.log("Best directory to delete has size %d", bestMatch)
        console.log(`Parsing and calculating part 2 took: ${part2Took[0] + part2Took[1] / 1e9} seconds`)
    }
}

main()
