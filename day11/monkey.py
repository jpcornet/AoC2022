#!/usr/bin/env python3

import sys
import re
import time
import math

def main():
    if not len(sys.argv) in (2,3):
        print("Give input file", file=sys.stderr)
        exit(-1)
    panic = False
    if len(sys.argv) == 3:
        if sys.argv[2] != 'panic':
            print("Just add 'panic' to increase panic level", file=sys.stderr)
            exit(-1)
        panic = True
    monkeys = parse_monkeys(sys.argv[1], panic=panic)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    for roundnum in range(1, 21 if not panic else 10001):
        maxmonkey = max([ max(m["items"]) if m["items"] else 0 for m in monkeys ])
        print(f"== Starting round {roundnum}. Max item value is {maxmonkey}")
        do_one_round(monkeys, panic)
        if roundnum in (1, 20) or roundnum in range(1000, 10000, 1000):
            print(f"== After round {roundnum} ==")
            for i in range(0, len(monkeys)):
                print(f"Monkey {i} inspected items {monkeys[i]['inspected']} times")
    endtime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print("== after all rounds")
    inspected = [ m['inspected'] for m in monkeys ]
    for i in range(0, len(monkeys)):
        print(f"Monkey {i} inspected items {inspected[i]} items")
    inspected.sort(reverse=True)
    print(f"Most active monkeys had {inspected[0]} and {inspected[1]} items. Monkey business is {inspected[0] * inspected[1]}")
    print(f"Rounds took: {(endtime - starttime) / 1e6}ms")

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
    plusop = re.fullmatch(r'old\s*\+\s*(\d+)', evalstr)
    if plusop:
        term = int(plusop[1])
        if moditem:
            return lambda old: (old + term) % moditem
        else:
            return lambda old: old + term
    multop = re.fullmatch(r'old\s*\*\s*(\d+)', evalstr)
    if multop:
        fact = int(multop[1])
        if moditem:
            return lambda old: (old * fact) % moditem
        else:
            return lambda old: old * fact
    sqop = re.fullmatch(r'old\s*\*\s*old', evalstr)
    if sqop:
        if moditem:
            return lambda old: (old * old) % moditem
        else:
            return lambda old: old * old
    dblop = re.fullmatch(r'old\s*\+\s*old', evalstr)
    if dblop:
        if moditem:
            return lambda old: (old + old) % moditem
        else:
            return lambda old: old + old
    raise ValueError(f"Parsing op {evalstr} not implemented")

def create_test(divisible, throwtrue, throwfalse):
    return lambda item: throwtrue if item % divisible == 0 else throwfalse

def parse_monkeys(filename: str, panic=False) -> list:
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
    # init at 3 because... in part 1 we still divide by 3. Which is impossible to do in a non-multiple-of-3 modulo.
    moditem = 1 if panic else 3
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
        moditem = math.lcm(moditem, monkeys[monkey_num]["divisible"])
    print(f"Product of all divisible tests is {moditem}")
    # now create callbacks for each monkey
    for m in monkeys:
        m["operation"] = create_eval(m["opstr"], moditem if panic else None)
        m["test"] = create_test(m["divisible"], m["truedest"], m["falsedest"])
    return monkeys

if __name__ == "__main__":
    main()