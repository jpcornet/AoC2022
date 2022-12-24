#!/usr/bin/env python3

import sys
import time

class EncryptedList:
    def __init__(self, filename):
        self.l = []
        self.zeroat = None

        for line in open(filename, "r"):
            self.l.append((int(line.rstrip()), False))

    def __str__(self):
        result = ""
        for n, done in self.l:
            if len(result) > 0:
                result += " "
            if not done:
                result += ","
            result += str(n)
        return result

    def __len__(self):
        return len(self.l)

    def move(self, i):
        if self.l[i][1]:
            # already done
            return False
        n = self.l[i][0]
        newpos = (i + n) % (len(self.l)-1)
        if newpos == i:
            self.l[i] = (n, True)
            return True
        elif newpos < i:
            self.l = self.l[:newpos] + [(n, True)] + self.l[newpos:i] + self.l[i+1:]
            self.zeroat = None
        else:
            self.l = self.l[:i] + self.l[i+1:newpos+1] + [(n, True)] + self.l[newpos+1:]
            self.zeroat = None
        return True

    def offset0(self, i):
        if self.zeroat == None:
            zeros = [ i for i in range(0,len(self.l)) if self.l[i][0] == 0 ]
            if not zeros:
                print("No zeros in list", file=sys.stderr)
                return
            if len(zeros) > 1:
                print("Too many zeros in list, taking first")
            self.zeroat = zeros[0]
        return self.l[(self.zeroat + i) % len(self.l)][0]

# another implementation, allowing multiple move operations
class EncryptedList2:
    def __init__(self, filename):
        self.l = []
        self.keypos = None

        lineno = 0
        for line in open(filename, "r"):
            self.l.append([int(line.rstrip()), lineno])
            lineno += 1

    def __str__(self):
        sorted = self.l[:]
        list.sort(sorted, key=lambda x: x[1])
        return " ".join([ f"({x[1]})={x[0]}" for x in sorted ])

    def __len__(self):
        return len(self.l)

    def move(self, i):
        n = self.l[i][0]
        curpos = self.l[i][1]
        newpos = (curpos + n) % (len(self.l)-1)
        if newpos == curpos:
            return
        elif newpos < curpos:
            # everything between newpos and curpos moves 1 up
            for e in self.l:
                if newpos <= e[1] < curpos:
                    e[1] += 1
            self.l[i][1] = newpos
            self.keypos = None
        else:
            # curpos < newpos. everything between curpos and newpos moves 1 down
            for e in self.l:
                if curpos < e[1] <= newpos:
                    e[1] -= 1
            self.l[i][1] = newpos
            self.keypos = None

    def offset0(self, i):
        if self.keypos == None:
            sorted = [ [i, self.l[i][1]] for i in range(0,len(self.l)) ]
            list.sort(sorted, key=lambda x: x[1])
            zeros = [ self.l[i][1] for i in range(0,len(self.l)) if self.l[i][0] == 0 ]
            if not zeros:
                print("No zeros in list", file=sys.stderr)
                return
            if len(zeros) > 1:
                print("Too many zeros in list, taking first")
            self.keypos = [ x[0] for x in sorted ]
            self.zeroat = zeros[0]
        return self.l[self.keypos[(self.zeroat + i) % len(self.l)]][0]

    def reset(self):
        for i in range(0, len(self.l)):
            self.l[i][1] = i

    def decrypt(self, key):
        for i in range(0, len(self.l)):
            self.l[i][0] *= key

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    elist = EncryptedList2(sys.argv[1])
    parsetime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    for i in range(0, len(elist)):
        elist.move(i)
    print("Result:")
    print(str(elist))
    part1 = 0
    for x in (1000,2000,3000):
        num = elist.offset0(x)
        print(f"Number {x} after 0 is {num}")
        part1 += num
    part1time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"part1={part1}")
    print(f"parsing took: {(parsetime - starttime)/1000}µs")
    print(f"part1 took: {(part1time - parsetime)/1000}µs")
    elist.reset()
    elist.decrypt(811589153)
    for rounds in range(0,10):
        for i in range(0, len(elist)):
            elist.move(i)
        print(f"After round {rounds}, list is: {str(elist)}")
    part2 = 0
    for x in (1000,2000,3000):
        num = elist.offset0(x)
        print(f"Number {x} after 0 is {num}")
        part2 += num
    part2time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"part2={part2}")
    print(f"part 2 took: {(part2time - part1time)/1000000}ms")

if __name__ == "__main__":
    main()