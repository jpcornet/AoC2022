package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Pos struct{ x, y int }

type Direction int8

const (
	RIGHT Direction = iota
	DOWN
	LEFT
	UP
)

type Blizzard struct {
	pos Pos
	dir Direction
}

type Field struct {
	width, height int
	in, out       Pos
	blizzards     []Blizzard
}

func parse_input(filename string) Field {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(inbuf[:]), "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	var field Field
	field.height = len(lines)
	field.width = len(lines[0])
	field.blizzards = make([]Blizzard, 0, 10)
	for y, line := range lines {
		if y == 0 || y == field.height-1 {
			// expect the wall and the entry/exit
			entrypos := strings.IndexByte(line, '.')
			if entrypos < 0 {
				panic("Invalid input")
			}
			// make sure we read a valid string
			should_be := strings.Repeat("#", entrypos) + "." + strings.Repeat("#", field.width-entrypos-1)
			if should_be != line {
				panic(fmt.Sprintf("Invalid input. At line %d expected [%s] got [%s]", y, should_be, line))
			}
			if y == 0 {
				field.in = Pos{entrypos, y}
			} else {
				field.out = Pos{entrypos, y}
			}
		} else {
			for x, char := range line {
				if x == 0 || x == field.width-1 {
					if char != '#' {
						panic(fmt.Sprintf("Expected wall at %d,%d got [%c]", x, y, char))
					}
				} else if char == '#' {
					panic(fmt.Sprintf("Unexpected wall at %d,%d", x, y))
				} else if char != '.' {
					dir := Direction(strings.IndexByte(">v<^", byte(char)))
					if dir == Direction(-1) {
						panic(fmt.Sprintf("Unexpected character at %d,%d: [%c]", x, y, char))
					}
					field.blizzards = append(field.blizzards, Blizzard{pos: Pos{x, y}, dir: dir})
				}
			}
		}
	}
	return field
}

type Dir struct{ dx, dy int }

func (f Field) MoveBlizzards() {
	directions := []Dir{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	for nr, b := range f.blizzards {
		d := directions[b.dir]
		newpos := Pos{b.pos.x + d.dx, b.pos.y + d.dy}
		if newpos.x <= 0 {
			newpos.x = f.width - 2
		} else if newpos.x >= f.width-1 {
			newpos.x = 1
		}
		if newpos.y <= 0 {
			newpos.y = f.height - 2
		} else if newpos.y >= f.height-1 {
			newpos.y = 1
		}
		f.blizzards[nr] = Blizzard{pos: newpos, dir: b.dir}
	}
}

// Path is a collection of posistions. Current position is the tail
type Path []Pos

func walk_path(f Field, start, finish Pos) (int, Path) {
	current_pos := make([]Path, 1)
	current_pos[0] = []Pos{start}
	// possible directions we can walk. Standing still is also an option.
	directions := []Dir{{0, 0}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	step := 1
	for {
		//fmt.Printf("Start step %d, positions to consider: %d\n", step, len(current_pos))
		// start by moving the blizzards to the next position. From there we can figure out where we can go.
		f.MoveBlizzards()
		// determine where the blizzards are, in a map
		blizzards := make(map[Pos]bool)
		for _, b := range f.blizzards {
			blizzards[b.pos] = true
		}
		// keep track of positions we've already seen in this move, to prevent duplicates
		seen := make(map[Pos]bool)
		// collect a list of all possible next moves
		next_pos := make([]Path, 0, 10)
		for _, cur := range current_pos {
			at := cur[len(cur)-1]
			for _, d := range directions {
				newpos := Pos{at.x + d.dx, at.y + d.dy}
				if newpos == finish {
					// found the exit!
					return step, cur
				}
				if newpos != start && (newpos.x <= 0 || newpos.x >= f.width-1 || newpos.y <= 0 || newpos.y >= f.height-1) {
					// we cannot walk into the wall or off the field
					continue
				}
				if _, has_blizzard := blizzards[newpos]; has_blizzard {
					// we cannot walk into a blizzard
					continue
				}
				if _, is_seen := seen[newpos]; is_seen {
					// we've already seen this position this round
					continue
				}
				// we can go here. Append newpos to the path and store as a possible solution
				newpath := make(Path, len(cur)+1)
				copy(newpath, cur)
				newpath[len(cur)] = newpos
				next_pos = append(next_pos, newpath)
				seen[newpos] = true
			}
		}
		if len(next_pos) == 0 {
			panic("No solution possible")
		}
		// prepare for the next step
		step++
		current_pos = next_pos
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	field := parse_input(os.Args[1])
	parsetime := time.Now()
	steps, path := walk_path(field, field.in, field.out)
	part1time := time.Now()
	fmt.Printf("part 1 Steps: %d , Path taken: %v\n", steps, path)
	steps2, path2 := walk_path(field, field.out, field.in)
	fmt.Printf("part 2 going back, steps: %d, path taken: %v\n", steps2, path2)
	steps3, path3 := walk_path(field, field.in, field.out)
	part2time := time.Now()
	fmt.Printf("part 2 return, steps: %d, path taken: %v\n", steps3, path3)
	fmt.Printf("part 2 total time: %d\n", steps+steps2+steps3)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", part1time.Sub(parsetime))
	fmt.Printf("part 2 took: %s\n", part2time.Sub(part1time))
}
