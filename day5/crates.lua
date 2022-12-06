#!/usr/bin/env lua

-- uses the lpeg and chronos extensions
lpeg = require "lpeg"
chronos = require "chronos"

function main (cmdline)
    local starttime = chronos.nanotime()
    assert(#cmdline == 1, "Provide input file")
    local f = assert(io.open(cmdline[1]), "Cannot read file " .. cmdline[1])
    local parsed = parseinput(f)
    local parsetime = chronos.nanotime()
    print("Parsing took: " .. (parsetime - starttime))
    local stk = parsed.stacks
    do_moves_part2(parsed.moves, stk)
    local endtime = chronos.nanotime()
    print("after moves:")
    local message = ""
    for _, stack in ipairs(stk) do
        message = message .. stack:sub(-1)
    end
    print("Message: " .. message)
    print("Processing took: " .. (endtime - starttime) * 1000 .. "ms")
end

function parseinput (f)
    local stacks = {}
    -- augment lpeg with convenient matchers
    lpeg.locale(lpeg)
    -- build the expression to match a line of crates. either blanks or [letter]
    -- blanks return nil as the captured value, letters return themselves
    local item = lpeg.P("   ") * lpeg.Cc(nil) + lpeg.P("[") * lpeg.C(lpeg.alpha) * lpeg.P("]")
    local crateline = lpeg.Ct((item * lpeg.P(" "))^0 * item^-1) * -1
    local capnumber = lpeg.digit^1 / tonumber
    local countline = lpeg.Ct(lpeg.space^0 * (capnumber * lpeg.space^0)^1) * -1
    for l in f:lines() do
        crates = crateline:match(l)
        if crates then
            -- put crates in stacks.
            for i, v in pairs(crates) do
                if not stacks[i] then
                    stacks[i] = v
                else
                    stacks[i] = stacks[i] .. v
                end
            end
        else
            -- should now match a line with stack counters
            nums = countline:match(l)
            if not nums then
                error("Cannot parse number line " .. l)
            end
            for i, v in pairs(nums) do
                if i ~= v then
                    error(("Unexpected number %s != %s on number line: %s"):format(i, v, l))
                end
            end
            break
        end
    end
    -- the stacks are now reverse-way around, flip them
    for i, stack in pairs(stacks) do
        stacks[i] = stacks[i]:reverse()
    end
    moves = {}
    moveline = lpeg.Ct(lpeg.P("move ") * lpeg.Cg(capnumber, "n")* lpeg.P(" from ") * lpeg.Cg(capnumber, "from") * lpeg.P(" to ") * lpeg.Cg(capnumber, "to")) * -1
    for l in f:lines() do
        if #l > 0 then
            move = moveline:match(l)
            if not move then
                error("Cannot parse move line: " .. moveline)
            end
            table.insert(moves, move)
        end
    end
    return { stacks=stacks, moves=moves }
end

function do_moves_part1 (moves, stacks)
    for _, move in ipairs(moves) do
        -- print(("Moving %s from %s to %s"):format(move.n, move.from, move.to))
        fromlen = #stacks[move.from]
        stacks[move.to] = stacks[move.to] .. stacks[move.from]:sub(fromlen - move.n + 1):reverse()
        stacks[move.from] = stacks[move.from]:sub(1, fromlen - move.n)
    end
end

function do_moves_part2 (moves, stacks)
    for _, move in ipairs(moves) do
        -- print(("Moving %s from %s to %s"):format(move.n, move.from, move.to))
        fromlen = #stacks[move.from]
        stacks[move.to] = stacks[move.to] .. stacks[move.from]:sub(fromlen - move.n + 1)
        stacks[move.from] = stacks[move.from]:sub(1, fromlen - move.n)
    end
end

function showstacks (stacks)
    for i, v in pairs(stacks) do
        print("Stack["..i.."] = " .. v)
    end
end

function showtable (t)
    for i, v in pairs(t) do
        print(("[%s]=%s"):format(i,v))
    end
end

main(arg)
