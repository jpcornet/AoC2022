#!/usr/bin/env python

import sys

class Hill:
    def __init__(self, filename):
        self.height = []
        self.xsize = None
        self.ysize = None
        self.start = None
        self.target = None

        for line in open(filename, "r"):
            line = line.rstrip()
            if len(line) == 0 or self.xsize and self.xsize != len(line):
                raise ValueError(f"Invalid input line: {line}")
            elif not self.xsize:
                self.xsize = len(line)
            s = line.find("S")
            if s != -1:
                line = line.replace('S', 'a')
                self.start = (s, len(self.height))
            e = line.find('E')
            if e != -1:
                line = line.replace('E', 'z')
                self.target = (e, len(self.height))
            self.height.append(line)
        self.ysize = len(self.height)
        if not self.start or not self.target:
            raise ValueError("No Start/End in input")

    def __str__(self):
        retstr = []
        for y in range(0, self.ysize):
            line = self.height[y]
            if y == self.start[1]:
                x = self.start[0]
                line = line[0:x] + 'S' + line[x+1:]
            if y == self.target[1]:
                x = self.target[0]
                line = line[0:x] + 'E' + line[x+1:]
            retstr.append(line)
        return "\n".join(retstr) + "\n"

    def walk(self):
        # start walking from start
        self.walkers = [ self.start ]
        # initialize all disntances to infinity
        self.distance = [ [sys.maxsize] * self.xsize for _ in range(0, self.ysize) ]
        # start position has distance 0 obviously
        self.distance[self.start[1]][self.start[0]] = 0
        while self.walkers:
            new_walkers = set()
            for w in self.walkers:
                for neww in self.walk_from(w):
                    new_walkers.add(neww)
            self.walkers = new_walkers
            if self.target in new_walkers:
                return self.distance[self.target[1]][self.target[0]]
        return None

    def walk_from(self, pos):
        x, y = pos
        curheight = self.height[y][x]
        curdist = self.distance[y][x]
        new_walkers = []
        for dx, dy in ( (0, 1), (0, -1), (1, 0), (-1, 0) ):
            if x + dx not in range(0, self.xsize) or y + dy not in range (0, self.ysize):
                # do not walk off the edge
                continue
            if self.distance[y + dy][x + dx] <= curdist:
                # do not walk to a position that we already reached and is closer
                continue
            if ord(self.height[y + dy][x + dx]) > ord(curheight) + 1:
                # do not walk a path that is too steep
                continue
            print(f"Walking from {x},{y} to {x+dx},{y+dy} distance {curdist+1}, new height {self.height[y+dy][x+dx]}")
            self.distance[y + dy][x + dx] = curdist + 1
            # in the next round, walk from that spot
            new_walkers.append( (x + dx, y + dy) )
        return new_walkers

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    hill = Hill(sys.argv[1])
    steps = hill.walk()
    print(f"Number of steps to target: {steps}")

if __name__ == "__main__":
    main()
