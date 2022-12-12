#!/usr/bin/env python3

import sys
import re

def main():
    if len(sys.argv) != 2:
        print("Give input file", file=sys.stderr)
        exit(-1)
    monkeys = parse_monkeys(sys.argv[1])
    for roundnum in range(1, 21):
        do_one_round(monkeys)
        print(f"After round {roundnum}, monkeys hold items with worry levels:")
        for i in range(0, len(monkeys)):
            print(f"Monkey {i}: {', '.join([str(x) for x in monkeys[i]['items']])}")
    inspected = [ m['inspected'] for m in monkeys ]
    for i in range(0, len(monkeys)):
        print(f"Monkey {i} inspected {inspected[i]} items")
    inspected.sort(reverse=True)
    print(f"Most active monkeys had {inspected[0]} and {inspected[1]} items. Monkey business is {inspected[0] * inspected[1]}")

def do_one_round(monkeys):
    print("== Start of round")
    for monkeynum in range(0, len(monkeys)):
        m = monkeys[monkeynum]
        while m["items"]:
            item = m["items"].pop(0)
            m["inspected"] += 1
            item = m["operation"](item)
            item = item // 3
            throw = m["test"](item)
            monkeys[throw]["items"].append(item)

def create_eval(evalstr):
    return lambda old: eval(evalstr, None, {"old": old})

def create_test(divisible, throwtrue, throwfalse):
    return lambda item: throwtrue if item % divisible == 0 else throwfalse

def parse_monkeys(filename: str) -> list:
    input = open(filename, "r").read()
    # regex to parse one monkey
    monkey_re = re.compile(r'''
        Monkey\s(?P<monkey>\d+): *\n
        \s+Starting\ items:\s*(?P<items>[\d\ ,]+)\n
        \s+Operation:\ new\ =\s*(?P<opstr>(old|[*+\ \d])+)\n
        \s+Test:\ divisible\ by\ (?P<divisible>\d+)\ *\n
        \s+If\ true:\ throw\ to\ monkey\ (?P<truedest>\d+)\ *\n
        \s+If\ false:\ throw\ to\ monkey\ (?P<falsedest>\d+)\ *\n+
    ''', re.VERBOSE)
    monkeys = []
    for m_monkey in monkey_re.finditer(input):
        monkey_num = int(m_monkey["monkey"])
        while len(monkeys) <= monkey_num:
            monkeys.append(None)
        print(f"Parsed monkey {monkey_num}: item={m_monkey['items']}, op={m_monkey['opstr']}, divtest={m_monkey['divisible']}, true={m_monkey['truedest']}, false={m_monkey['falsedest']}")
        monkeys[monkey_num] = {
            "items": [ int(i) for i in re.split(r',\s*', m_monkey["items"]) ],
            "operation": create_eval(m_monkey["opstr"]),
            "test": create_test(int(m_monkey["divisible"]), int(m_monkey["truedest"]), int(m_monkey["falsedest"])),
            "inspected": 0,
        }
    return monkeys

if __name__ == "__main__":
    main()