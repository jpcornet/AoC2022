#!/usr/bin/env python3

import sys
import time
import re

class MonkeySheet:
    def __init__(self, filename):
        self.eq = {}
        self.uses = {}
        self.used = {}
        self.value = {}

        for line in open(filename, "r"):
            name, sep, content = line.rstrip().partition(':')
            if not sep:
                print(f"Bad input: {line}", file=sys.stderr)
                exit(-1)
            content = content.lstrip()
            if content.isdigit():
                # plain content
                self.eq[name] = lambda: int(content)
                self.value[name] = int(content)
            else:
                matches = re.match(r'(\w+)\s*([-+*/])\s*(\w+)', content)
                if not matches:
                    print(f"Bad input: {line}.", file=sys.stderr)
                    exit(-1)
                self.uses[name] = [ matches[1], matches[3] ]
                for i in 1, 3:
                    if matches[i] not in self.used:
                        self.used[matches[i]] = [name]
                    else:
                        self.used[matches[i]].append(name)
                self.eq[name] = self.create_eval(matches[1], matches[2], matches[3])

    def create_eval(self, v1, op, v2):
        if op == '-':
            return lambda: self.value[v1] - self.value[v2]
        elif op == '+':
            return lambda: self.value[v1] + self.value[v2]
        elif op == '*':
            return lambda: self.value[v1] * self.value[v2]
        elif op == '/':
            return lambda: self.value[v1] / self.value[v2]
        else:
            print("logic error", file=sys.stderr)
            exit(-1)

    def try_solve(self, name):
        print(f"In try_solve({name})")
        if name in self.value:
            return True
        if name not in self.eq:
            print(f"Trying to solve {name}, but that has no equation")
            return False
        if name not in self.uses:
            print(f"Trying to solve {name}, but that lists no used names")
            return False
        for dep in self.uses[name]:
            if dep not in self.value:
                return False
        self.value[name] = self.eq[name]()
        print(f"Solved {name} = {self.value[name]}")
        return True

    def solve_all(self):
        # loop along everything that has a value, check which equations use those, and try solving those
        solved = list(self.value)
        for n in solved:
            print(f"{n} has been solved, finding used equations...")
            if n in self.used:
                for other in self.used[n]:
                    print(f"{n} is used by {other}")
                    if other not in self.value:
                        if self.try_solve(other):
                            solved.append(other)

    def get_val(self, name):
        return self.value[name]

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    sheet = MonkeySheet(sys.argv[1])
    parsetime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    sheet.solve_all()
    root = sheet.get_val("root")
    part1time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"part1: root={root}")
    print(f"parsing input took: {(parsetime-starttime)/1000}µs")
    print(f"solving sheet took: {(part1time-parsetime)/1000}µs")

if __name__ == "__main__":
    main()