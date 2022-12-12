#!/usr/bin/env python

import sys
import time

class Hill:
    def __init__(self, filename):
        self.height = []
        self.xsize = None
        self.ysize = None
        self.start = None
        self.target = None
        self.path = None

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
        as_string = "\n".join(retstr) + "\n"
        if self.path:
            # prepare a blank canvas for the route
            walk_pos = [ ["."] * self.xsize for _ in range(0, self.ysize) ]
            # what to use as the direction markers
            trans = {
                (-1, 0): ">",
                (1, 0): "<",
                (0, -1): "v",
                (0, 1): "^",
            }
            # walk back freom target to start, following the path while colouring that in
            x, y = self.target
            walk_pos[y][x] = 'E'
            while self.path[y][x] and (x != self.start[0] or y != self.start[1]):
                newx, newy = self.path[y][x]
                dx = newx - x
                dy = newy - y
                walk_pos[newy][newx] = trans[(dx, dy)]
                x, y = newx, newy
            as_string += "\n" + "\n".join([ "".join(l) for l in walk_pos]) + "\n"
        return as_string

    def walk(self):
        # start walking from start
        walkers = [ self.start ]
        # initialize all disntances to infinity
        self.distance = [ [sys.maxsize] * self.xsize for _ in range(0, self.ysize) ]
        # start position has distance 0 obviously
        self.distance[self.start[1]][self.start[0]] = 0
        # record the path taken
        self.path = [ [None] * self.xsize for _ in range(0, self.ysize) ]
        while walkers:
            new_walkers = set()
            for w in walkers:
                for neww in self.walk_from(w):
                    new_walkers.add(neww)
            walkers = new_walkers
            if self.target in walkers:
                return self.distance[self.target[1]][self.target[0]]
        return None

    def walk_from(self, pos, back=False):
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
            if back:
                if ord(self.height[y + dy][x + dx]) < ord(curheight) - 1:
                    # do not walk back a path that is too steep
                    continue
            elif ord(self.height[y + dy][x + dx]) > ord(curheight) + 1:
                # do not walk a path that is too steep
                continue
            # record (new) shortest distance to this position
            self.distance[y + dy][x + dx] = curdist + 1
            # record (new) previous location from here
            self.path[y + dy][x + dx] = (x, y)
            # in the next round, walk from that spot
            new_walkers.append( (x + dx, y + dy) )
        return new_walkers

    def walkback(self):
        # walk from target back to any starting position
        walkers = [ self.target ]
        # initialize all disntances to infinity
        self.distance = [ [sys.maxsize] * self.xsize for _ in range(0, self.ysize) ]
        # start position has distance 0 obviously
        self.distance[self.target[1]][self.target[0]] = 0
        # record the path taken
        self.path = [ [None] * self.xsize for _ in range(0, self.ysize) ]
        while walkers:
            new_walkers = set()
            for w in walkers:
                for neww in self.walk_from(w, back=True):
                    new_walkers.add(neww)
            walkers = new_walkers
            done = [ w for w in walkers if self.height[w[1]][w[0]] == 'a' ]
            if done:
                if len(done) > 1:
                    print(f"Multiple solutions to walk back: {done}")
                # make sure it prints nicely
                self.start = self.target
                self.target = done[0]
                return self.distance[done[0][1]][done[0][0]]
        return None

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    starttime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    hill = Hill(sys.argv[1])
    parsetime = time.clock_gettime_ns(time.CLOCK_REALTIME)
    steps = hill.walk()
    part1time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"Number of steps to target: {steps}")
    path1str = str(hill)
    part2start = time.clock_gettime_ns(time.CLOCK_REALTIME)
    steps2 = hill.walkback()
    part2time = time.clock_gettime_ns(time.CLOCK_REALTIME)
    print(f"Number of steps back from target to any start: {steps2}")
    print(f"\nPath from start to target:\n{path1str}\nPart from target to closest starting pos:\n{str(hill)}")
    print(f"Loading data took: {(parsetime - starttime) / 1000}µs")
    print(f"Start to target path took: {(part1time - parsetime) / 1000}µs")
    print(f"Target to any start took: {(part2time - part2start) / 1000}µs")

if __name__ == "__main__":
    main()
