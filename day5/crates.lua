#!/usr/bin/env lua

-- uses the lpeg and chronos extensions
lpeg = require "lpeg"
chronos = require "chronos"

function main (cmdline)
    local starttime = chronos.nanotime()
    assert(#cmdline == 1, "Provide input file")
    local f = assert(io.open(cmdline[1]), "Cannot read file " .. cmdline[1])
    local parsed = parseinput(f)
    print("stacks after parsing:")
    local stk = parsed.stacks
    showstacks(stk)
    do_moves(parsed.moves, stk)
    print("after moves:")
    showstacks(stk)
    local message = ""
    for _, stack in ipairs(stk) do
        message = message .. stack[#stack]
    end
    print("Message: " .. message)
end

function parseinput (f)
    local stacks = {}
    -- augment lpeg with convenient matchers
    lpeg.locale(lpeg)
    -- build the expression to match a line of crates. either blanks or [letter]
    -- blanks return nil as the captured value, letters return themselves
    local item = lpeg.P("   ") * lpeg.Cc(nil) + lpeg.P("[") * lpeg.C(lpeg.alpha) * lpeg.P("]")
    local crateline = lpeg.Ct((item * lpeg.P(" "))^0 * item) * -1
    local capnumber = lpeg.digit^1 / tonumber
    local countline = lpeg.Ct(lpeg.space^0 * (capnumber * lpeg.space^0)^1) * -1
    for l in f:lines() do
        crates = crateline:match(l)
        if crates then
            -- put crates in stacks. Note stacks are reversed.
            for i, v in pairs(crates) do
                if not stacks[i] then
                    stacks[i] = {}
                end
                table.insert(stacks[i], 1, v)
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

function do_moves (moves, stacks)
    for _, move in ipairs(moves) do
        for _ = 1, move.n do
            table.insert(stacks[move.to], table.remove(stacks[move.from]))
        end
    end
end

function showstacks (stacks)
    for i, v in pairs(stacks) do
        print("Stack["..i.."] = " .. table.concat(v, ","))
    end
end

function showtable (t)
    for i, v in pairs(t) do
        print(("[%s]=%s"):format(i,v))
    end
end

main(arg)
