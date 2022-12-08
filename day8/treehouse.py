#!/usr/bin/env python3

import sys

def read_input(infile: str) -> list[str]:
    print(f"Reading file {infile}")
    return open(infile, "r").read().splitlines()

def count_visible_trees(forest: list[str]) -> int:
    count = 0
    max_y = len(forest)
    max_x = len(forest[0])
    for y in range(0, max_y):
        assert len(forest[y]) == max_x, "Rows in forest not of equals size"
        for x in range(0, max_x):
            if is_visible_from_edge(forest, x, y):
                count += 1
    return count

def is_visible_from_edge(forest, x, y) -> bool:
    max_y = len(forest)
    max_x = len(forest[0])
    height = forest[y][x]
    # now look around in 4 directions
    for dx, dy in ((-1, 0), (1, 0), (0, -1), (0, 1)):
        x2, y2 = x, y
        while True:
            # have we reached the edge?
            if x2 in (0, max_x - 1) or y2 in (0, max_y - 1):
                return True
            # take 1 step
            x2 += dx
            y2 += dy
            # check if current tree is lower than our tree, otherwise try another direction
            if forest[y2][x2] >= height:
                break
    # tried all direction, it is not visible
    return False

def visible_trees(forest, x, y, dx, dy) -> int:
    count = 0
    max_y = len(forest)
    max_x = len(forest[0])
    height = forest[y][x]
    while True:
        if x in (0, max_x - 1) or y in (0, max_y - 1):
            return count
        count += 1
        x += dx
        y += dy
        if forest[y][x] >= height:
            return count

def best_scenic_tree(forest) -> (int, int, int):
    max_y = len(forest)
    max_x = len(forest[0])
    best_x, best_y, best_score = None, None, 0
    # all trees on the edge have a score of 0, because in 1 direction they see 0 trees, so skip those
    for y in range(1, max_y - 1):
        for x in range(1, max_x - 1):
            this_score = 1
            for dx, dy in ((-1, 0), (1, 0), (0, -1), (0, 1)):
                this_score *= visible_trees(forest, x, y, dx, dy)
            if this_score > best_score:
                best_score = this_score
                best_x, best_y = x, y
    return best_score, best_x, best_y

def main():
    if len(sys.argv) != 2:
        print("Specify input file", file=sys.stderr)
        exit(-1)
    forest = read_input(sys.argv[1])
    part1 = count_visible_trees(forest)
    print(f"Number of visible trees from the edge: {part1}")
    part2, x, y = best_scenic_tree(forest)
    print(f"Best scenic score is {part2} at tree {x},{y}")

if __name__ == "__main__":
    main()