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

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    elist = EncryptedList(sys.argv[1])
    i = 0
    while i < len(elist):
        if not elist.move(i):
            i += 1
    print("Result:")
    print(str(elist))
    part1 = 0
    for x in (1000,2000,3000):
        num = elist.offset0(x)
        print(f"Number {x} after 0 is {num}")
        part1 += num
    print(f"part1={part1}")

if __name__ == "__main__":
    main()