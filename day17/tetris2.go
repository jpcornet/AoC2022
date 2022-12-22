package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// each bit is one rock, bit 6 is leftmost, bit 0 is rightmost. 0 = air, 1 = rock
type Line uint8

const width = 7

type Stack struct {
	offset int
	lines  []Line
}

type Rock struct {
	x, y  int
	lines []Line
}

var rocks []Rock

var rocknr, streamnr int
var stack Stack
var stream string

func setup_field() {
	rocks = []Rock{
		{x: 4, y: 1, lines: []Line{0x78}},
		{x: 3, y: 3, lines: []Line{0x20, 0x70, 0x20}},
		{x: 3, y: 3, lines: []Line{0x70, 0x10, 0x10}}, // first row is bottom row of rock
		{x: 1, y: 4, lines: []Line{0x40, 0x40, 0x40, 0x40}},
		{x: 2, y: 2, lines: []Line{0x60, 0x60}},
	}
	stack.lines = make([]Line, 0, 10000)
	rocknr = 0
	streamnr = 0
}

func next_rock() Rock {
	result := rocks[rocknr]
	rocknr++
	if rocknr == len(rocks) {
		rocknr = 0
	}
	return result
}

func next_stream() int {
	schar := stream[streamnr]
	streamnr++
	if streamnr == len(stream) {
		streamnr = 0
	}
	switch schar {
	case '<':
		return -1
	case '>':
		return 1
	default:
		panic(fmt.Sprintf("Invalid stream char [%c] in input", schar))
	}
}

func has_overlap(rock Rock, x, y int) bool {
	for ry, l := range rock.lines {
		if y+ry >= len(stack.lines)+stack.offset {
			return false
		}
		if stack.lines[y+ry-stack.offset]&(l>>x) != 0 {
			return true
		}
	}
	return false
}

// adds rock to the stack. Returns true if the stack has just been pruned.
func add_rock_to_stack(rock Rock, x, y int) bool {
	hbar := Line((1 << width) - 1)
	barfound := -1
	for ry, l := range rock.lines {
		for y+ry >= len(stack.lines)+stack.offset {
			stack.lines = append(stack.lines, Line(0))
		}
		stack.lines[y+ry-stack.offset] |= l >> x
		if stack.lines[y+ry-stack.offset] == hbar {
			barfound = y + ry - stack.offset
		}
	}
	if barfound != -1 {
		stack.lines = stack.lines[barfound:]
		stack.offset += barfound
		return true
	} else {
		// find two adjacent lines that together form an hbar. Which is also impenetrable for any rock with size > 1
		for dy := len(stack.lines) - 2; dy >= y-stack.offset; dy-- {
			if (stack.lines[dy] | stack.lines[dy+1]) == hbar {
				stack.lines = stack.lines[dy:]
				stack.offset += dy
				show_stack()
				return true
			}
		}
	}
	return false
}

func drop_one_rock() bool {
	rock := next_rock()
	x := 2
	y := len(stack.lines) + stack.offset + 3
	for true {
		dx := next_stream()
		if x+dx >= 0 && x+dx+rock.x-1 < width && !has_overlap(rock, x+dx, y) {
			x += dx
		}
		if y == 0 || has_overlap(rock, x, y-1) {
			return add_rock_to_stack(rock, x, y)
		}
		y--
	}
	//NOTREACHED
	return false
}

func show_stack() {
	if len(stack.lines) > 10000 {
		fmt.Printf("Before showing stack, pruning to 10000. Real length was: %d\n", len(stack.lines))
		stack.offset += len(stack.lines) - 10000
		stack.lines = stack.lines[len(stack.lines)-10000:]
	}
	for y := len(stack.lines) - 1; y >= 0; y-- {
		strline := ""
		var leftmost Line = 1 << (width - 1)
		for x := 0; x < width; x++ {
			var char byte
			if stack.lines[y]&(leftmost>>x) != 0 {
				char = '#'
			} else {
				char = '.'
			}
			strline += string(char)
		}
		fmt.Printf("|%s|\n", strline)
	}
	if stack.offset > 0 {
		fmt.Printf("(..%d..)\n", stack.offset)
	}
	fmt.Printf("|%s|\n", strings.Join(make([]string, width+1), "-"))
	fmt.Printf("Stack size: %d\n", len(stack.lines)+stack.offset)
}

const snaplen = 10

type GamePos struct {
	rocknr, streamnr int
	stacksnap        [snaplen]Line
}

type GameProgress struct {
	num_rocks, stack_size int
}

type MemRepeat map[GamePos]GameProgress

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	bstream, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	stream = strings.TrimRight(string(bstream[:]), "\r\n")
	setup_field()
	loadtime := time.Now()
	fmt.Printf("Loading took: %s\n", loadtime.Sub(starttime))
	target_rocks := 1_000_000_000_000
	total_rocks := 0
	game_repeat := make(MemRepeat)
	for total_rocks < target_rocks {
		total_rocks++
		if drop_one_rock() {
			// the rock we just dropped pruned the stack.
			// try to see if we can find a repeat. Only if the stacksize fits in the GamePos struct
			if len(stack.lines) <= snaplen {
				this_pos := GamePos{
					rocknr:   rocknr,
					streamnr: streamnr,
				}
				copy(this_pos.stacksnap[:], stack.lines)
				progress, found := game_repeat[this_pos]
				if found {
					delta_rocks := total_rocks - progress.num_rocks
					delta_stack := len(stack.lines) + stack.offset - progress.stack_size
					repeats := (target_rocks - total_rocks) / delta_rocks
					total_rocks += repeats * delta_rocks
					stack.offset += repeats * delta_stack
				} else {
					game_repeat[this_pos] = GameProgress{
						num_rocks:  total_rocks,
						stack_size: len(stack.lines) + stack.offset,
					}
				}
			}
		}
	}
	part1time := time.Now()
	show_stack()
	fmt.Printf("Part 1 took: %s\n", part1time.Sub(starttime))
}
