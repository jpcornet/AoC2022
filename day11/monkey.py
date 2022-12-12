#!/usr/bin/env python3

import sys
import re

def main():
    if not len(sys.argv) in (2,3):
        print("Give input file", file=sys.stderr)
        exit(-1)
    monkeys = parse_monkeys(sys.argv[1])
    panic = False
    if len(sys.argv) == 3:
        if sys.argv[2] != 'panic':
            print("Just add 'panic' to increase panic level", file=sys.stderr)
            exit(-1)
        panic = True
    for roundnum in range(0, 20 if not panic else 10000):
        do_one_round(monkeys, panic)
        if roundnum in (0, 19) or roundnum in range(999, 9999, 1000):
            print(f"== After round {roundnum + 1} ==")
            for i in range(0, len(monkeys)):
                print(f"Monkey {i} inspected items {monkeys[i]['inspected']} times")
    print("== after all rounds")
    inspected = [ m['inspected'] for m in monkeys ]
    for i in range(0, len(monkeys)):
        print(f"Monkey {i} inspected items {inspected[i]} items")
    inspected.sort(reverse=True)
    print(f"Most active monkeys had {inspected[0]} and {inspected[1]} items. Monkey business is {inspected[0] * inspected[1]}")

def do_one_round(monkeys, is_panic):
    for monkeynum in range(0, len(monkeys)):
        m = monkeys[monkeynum]
        while m["items"]:
            item = m["items"].pop(0)
            m["inspected"] += 1
            item = m["operation"](item)
            if not is_panic:
                item = item // 3
            throw = m["test"](item)
            monkeys[throw]["items"].append(item)

def create_eval(evalstr, moditem):
    return lambda old: eval(evalstr, None, {"old": old}) % moditem

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
    moditem = 1
    for m_monkey in monkey_re.finditer(input):
        monkey_num = int(m_monkey["monkey"])
        while len(monkeys) <= monkey_num:
            monkeys.append(None)
        print(f"Parsed monkey {monkey_num}: item={m_monkey['items']}, op={m_monkey['opstr']}, divtest={m_monkey['divisible']}, true={m_monkey['truedest']}, false={m_monkey['falsedest']}")
        monkeys[monkey_num] = {
            "items": [ int(i) for i in re.split(r',\s*', m_monkey["items"]) ],
            "opstr": m_monkey["opstr"],
            "divisible": int(m_monkey["divisible"]),
            "truedest": int(m_monkey["truedest"]),
            "falsedest": int(m_monkey["falsedest"]),
            "inspected": 0,
        }
        # keep track of product of all divisible tests
        moditem *= monkeys[monkey_num]["divisible"]
    print(f"Product of all divisible tests is {moditem}")
    # now create callbacks for each monkey
    for m in monkeys:
        m["operation"] = create_eval(m["opstr"], moditem)
        m["test"] = create_test(m["divisible"], m["truedest"], m["falsedest"])
    return monkeys

if __name__ == "__main__":
    main()