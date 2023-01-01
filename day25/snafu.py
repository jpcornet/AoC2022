#!/usr/bin/env python3

import sys
import time

class Snafu:
    numchars = "=-012"

    def __init__(self, init):
        if type(init) == str:
            self.val = self.from_snafu(init)
        elif type(init) == int:
            self.val = init
        else:
            print(f"Cannot handle type {type(init)}", file=sys.stderr)
            exit(-1)

    def from_snafu(self, input):
        num = 0
        power = 1
        while len(input) > 0:
            char = input[-1]
            val = Snafu.numchars.find(char)
            if val == -1:
                print(f"Invalid character in SNAFU number: {char}")
                return
            num += (val - 2) * power
            power *= 5
            input = input[:-1]
        return num

    def __index__(self):
        return self.val

    def __str__(self):
        string = ""
        val = self.val
        while True:
            digit = val % 5
            if digit >= 3:
                digit -= 5
            string = Snafu.numchars[digit+2] + string
            val -= digit
            val //= 5
            if val == 0:
                return string

    def __add__(self, other):
        return Snafu(self.val + int(other))

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    total = Snafu(0)
    for line in open(sys.argv[1], "r"):
        line = line.rstrip()
        num = Snafu(line)
        #print(f"input [{line}] num: {int(num)} str: {str(num)}")
        total += num
    endtime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"Total num: {int(total)} str: {str(total)}")
    print(f"Took: {(endtime - starttime)/1e3}µs")

if __name__ == "__main__":
    main()
