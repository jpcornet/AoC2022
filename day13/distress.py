#!/usr/bin/env python3

import sys
import re

class Mylist:
    def __init__(self, input):
        self.l = []
        self.as_string = repr(input)
        for elem in input:
            if type(elem) == list:
                self.l.append(Mylist(elem))
            else:
                self.l.append(elem)

    def __eq__(self, other):
        if type(other) != Mylist:
            other = Mylist([other])
        if len(self.l) != len(other.l):
            return False
        for i in range(0, len(self.l)):
            if self.l[i] != other.l[i]:
                return False
        return True

    def __lt__(self, other):
        if type(other) == int:
            other = Mylist([other])
        for i in range(0, len(self.l)):
            if i >= len(other.l):
                return False
            if self.l[i] == other.l[i]:
                continue
            if type(self.l[i]) == type(other.l[i]):
                return self.l[i] < other.l[i]
            if type(self.l[i]) == int:
                return Mylist([self.l[i]]) < other.l[i]
            else:
                return self.l[i] < Mylist([other.l[i]])
        return len(self.l) < len(other.l)

    def __str__(self):
        return self.as_string

def parse_input(filename):
    input = open(filename, "r").read()
    # make sure the input is not rogue
    if not re.fullmatch(r'[\[\]\d,\s]+', input):
        raise ValueError("Input contains rogue characters")
    pairs = []
    for pair in input.split("\n\n"):
        pairs.append([ Mylist(eval(line)) for line in pair.rstrip().split("\n")])
    return pairs

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    pairs = parse_input(sys.argv[1])

    sum = 0
    for i in range(0, len(pairs)):
        if pairs[i][0] < pairs[i][1]:
            sum += i+1
    print(f"Sum of correct orders: {sum}")

    # convert to flat list and sort it
    flat = []
    for p in pairs:
        flat += p
    
    # add divider signals
    dividers = [ Mylist([[2]]), Mylist([[6]]) ]
    flat += dividers
    print("Sorting flat list")
    flat.sort()

    decoder = 1
    for i in range(0, len(flat)):
        print(f"{i+1}: {str(flat[i])}")
        if flat[i] in dividers:
            print("^^^ divider")
            decoder *= i+1
    print(f"Decoder: {decoder}")

if __name__ == "__main__":
    main()
