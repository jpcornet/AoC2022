#!/usr/bin/env python3

import sys
import time
import re

class MonkeySheet:
    def __init__(self, filename):
        self.eq = {}
        self.orig = {}
        self.uses = {}
        self.value = {}
        # keep a list of reverse equations, to be used if necessary
        self.reverse = {}
        self.allow_reverse = False
        self.busy = {}

        for line in open(filename, "r"):
            name, sep, content = line.rstrip().partition(':')
            if not sep:
                print(f"Bad input: {line}", file=sys.stderr)
                exit(-1)
            content = content.lstrip()
            if content.isdigit():
                self.set_value(name, int(content))
            else:
                matches = re.match(r'(\w+)\s*([-+*/])\s*(\w+)', content)
                if not matches:
                    print(f"Bad input: {line}.", file=sys.stderr)
                    exit(-1)
                v1 = matches[1]
                v2 = matches[3]
                self.uses[name] = [ v1, v2 ]
                self.eq[name] = self.create_eval(v1, matches[2], v2)
                e1, e2 = self.get_reverses(name, v1, matches[2], v2)
                if v1 not in self.reverse:
                    self.reverse[v1] = []
                if v2 not in self.reverse:
                    self.reverse[v2] = []
                self.reverse[v1].append((name, v2, self.create_eval(*e1)))
                self.reverse[v2].append((name, v1, self.create_eval(*e2)))

    def set_value(self, name, value):
        # plain content
        self.eq[name] = lambda: value
        self.value[name] = value
        self.orig[name] = True

    def get_reverses(self, name, v1, op, v2):
        if op == "+":
            return (name, '-', v2), (name, '-', v1)
        if op == '-':
            return (name, '+', v2), (v1, '-', name)
        if op == '*':
            return (name, '/', v2), (name, '/', v1)
        if op == '/':
            return (name, '*', v2), (v1, '/', name)

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

    def solve(self, name):
        print(f"In solve({name})")
        if name in self.busy:
            print(f"Already trying to solve {name}, giving up")
            return None
        self.busy[name] = True
        if name in self.value:
            print(f"Value for {name} already known: {self.value[name]}")
            del self.busy[name]
            return self.value[name]
        direct_solve = True
        if name in self.eq and name in self.uses:
            for dep in self.uses[name]:
                if dep not in self.value:
                    x = self.solve(dep)
                    if x == None:
                        direct_solve = False
                        break
            if direct_solve:
                self.value[name] = self.eq[name]()
                print(f"Solved {name} = {self.value[name]}")
                del self.busy[name]
                return self.value[name]

        if self.allow_reverse and name in self.reverse:
            ans = self.try_reverse(name)
            if ans != None:
                del self.busy[name]
                return ans
        else:
            print(f"Cannot solve {name}")
            del self.busy[name]
            return None

        del self.busy[name]
        if name not in self.eq:
            print(f"Trying to solve {name}, but that has no equation")
            return None
        if name not in self.uses:
            print(f"Trying to solve {name}, but that lists no used names")
            return None
        return None

    def try_reverse(self, name):
        print(f"Trying reverse lookup for {name}, {len(self.reverse[name])} available")
        for dep1, dep2, ev in self.reverse[name]:
            v1 = self.solve(dep1)
            if v1 == None:
                print(f"This reverse depending on {dep1} gives no answer")
                continue
            v2 = self.solve(dep2)
            if v2 == None:
                print(f"This reverse depending on {dep2} gives no answer")
                continue
            self.value[name] = ev()
            print(f"Solved {name} via reverse on {dep1} and {dep2} = {self.value[name]}")
            return self.value[name]
        print(f"None of the reverses for {name} worked out")
        return None

    def reset(self):
        test_these = [ x for x in self.value ]
        for n in test_these:
            if n not in self.orig:
                print(f"Deleted value for {n}, not original in spec")
                del self.value[n]

    def delete_eq(self, name):
        ok = False
        if name in self.eq:
            del self.eq[name]
            ok = True
        if name in self.value:
            del self.value[name]
            ok = True
        if not ok:
            print(f"Deleting {name} but no equation or value for it")
            exit(-1)

    def get_deps(self, name):
        return self.uses[name]

    def enable_reverse(self):
        self.allow_reverse = True

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    sheet = MonkeySheet(sys.argv[1])
    parsetime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    root = sheet.solve("root")
    part1time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"part1: root={root}")
    print(f"parsing input took: {(parsetime-starttime)/1000}µs")
    print(f"solving sheet took: {(part1time-parsetime)/1000}µs")
    # reset it so we can fix it for part 2
    sheet.reset()
    # remove the equation for us, "humn"
    sheet.delete_eq("humn")
    # get the dependencies for "root"
    rdeps = sheet.get_deps("root")
    sheet.delete_eq("root")
    # try solving each dependency. One of these should fail
    failed = None
    succeeded = None
    answer = None
    for name in rdeps:
        answer = sheet.solve(name)
        if answer != None:
            succeeded = name
        elif failed != None:
            print(f"Argh, trying to solve root dependencies but both {failed} and {name} failed!", file=sys.stderr)
            exit(-1)
        else:
            failed = name
    if failed == None:
        print("Odd, none of the root deps failed?", file=sys.stderr)
    # just set the failed symbol equal to the answer of the succeeded one
    sheet.set_value(failed, answer)
    # turn on deductive reasoning
    sheet.enable_reverse()
    # now solve for us, humn
    humn = sheet.solve("humn")
    part2time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"part2: humn={humn}")
    print(f"solving part2 took: {(part2time-part1time)/1000}µs")

if __name__ == "__main__":
    main()
