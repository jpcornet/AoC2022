package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Field []string

type PathElem struct {
	dist   int
	rotate int
}

type Path []PathElem

type Pos [2]int

var directions [4]Pos

type PosDir struct {
	pos Pos
	dir int
}

func init_vars() {
	directions = [4]Pos{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}
}

func parse_input(filename string) (Field, Path) {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fieldstr, pathstr, ok := strings.Cut(string(inbuf[:]), "\n\n")
	if !ok {
		panic("Invalid input")
	}
	field := strings.Split(fieldstr, "\n")
	if len(field[len(field)-1]) == 0 {
		field = field[:len(field)-1]
	}
	path := make(Path, 0, 10)
	pathstr = strings.TrimRight(pathstr, "\n")
	for len(pathstr) > 0 {
		var pe PathElem
		lroffset := strings.IndexAny(pathstr, "LR")
		if lroffset == -1 {
			lroffset = len(pathstr)
		} else if pathstr[lroffset] == 'L' {
			pe.rotate = -1
		} else if pathstr[lroffset] == 'R' {
			pe.rotate = 1
		}
		if lroffset != 0 {
			if i, err := strconv.Atoi(pathstr[:lroffset]); err == nil {
				pe.dist = i
			} else {
				panic(fmt.Sprintf("Invalid input in path, not a number: %s (%v)", pathstr[:lroffset], err))
			}
		}
		path = append(path, pe)
		if lroffset == len(pathstr) {
			break
		} else {
			pathstr = pathstr[lroffset+1:]
		}
	}
	return field, path
}

func get_startpos(field Field) PosDir {
	x := strings.Index(field[0], ".")
	if x == -1 {
		panic("Invalid input")
	}
	return PosDir{
		pos: Pos{x, 0},
		dir: 0,
	}
}

func walk_path(field Field, path Path, pos PosDir) PosDir {
	for _, p := range path {
		for i := 0; i < p.dist; i++ {
			newpos := Pos{pos.pos[0] + directions[pos.dir][0], pos.pos[1] + directions[pos.dir][1]}
			if newpos[0] >= 0 && newpos[1] >= 0 && newpos[1] < len(field) && newpos[0] < len(field[newpos[1]]) {
				chr := field[newpos[1]][newpos[0]]
				if chr == '.' {
					pos.pos = newpos
					continue
				} else if chr == '#' {
					break
				} else if chr != ' ' {
					panic(fmt.Sprintf("Invalid character [%c] in field at %d,%d", chr, newpos[0], newpos[1]))
				}
			}
			// We get here if we need to wrap around the field.
			wrapstart := pos.pos
			newpos = pos.pos
			for newpos[0] >= 0 && newpos[1] >= 0 && newpos[1] < len(field) && newpos[0] < len(field[newpos[1]]) && (field[newpos[1]][newpos[0]] == '.' || field[newpos[1]][newpos[0]] == '#') {
				pos.pos = newpos
				newpos[0] -= directions[pos.dir][0]
				newpos[1] -= directions[pos.dir][1]
			}
			if field[pos.pos[1]][pos.pos[0]] == '#' {
				pos.pos = wrapstart
				break
			}
		}
		pos.dir = (pos.dir + p.rotate + 4) % 4
	}
	return pos
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	field, path := parse_input(os.Args[1])
	init_vars()
	parsetime := time.Now()
	startpos := get_startpos(field)
	endpos := walk_path(field, path, startpos)
	walktime := time.Now()
	fmt.Printf("endpos: %v. Password: %d\n", endpos, 1000*(endpos.pos[1]+1)+4*(endpos.pos[0]+1)+endpos.dir)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", walktime.Sub(parsetime))
}
