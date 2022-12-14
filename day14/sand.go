package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Item uint8

const (
	Air Item = iota
	Rock
	Sand
)

type Field struct {
	xoffset int
	yoffset int
	xsize   int
	ysize   int
	space   [][]Item
	grains  int
	part1   int
}

type Coord struct {
	x int
	y int
}

type Line []Coord

func parse_input(filename string, pyramid_top int) Field {
	starttime := time.Now()
	inputstr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	// remove last empty line
	if len(inputlines[len(inputlines)-1]) == 0 {
		inputlines = inputlines[0 : len(inputlines)-1]
	}
	linesplittime := time.Now()
	lines := make([]Line, 0, 500)
	minx := math.MaxInt
	maxx := math.MinInt
	maxy := math.MinInt
	for _, linestr := range inputlines {
		newline := make(Line, 0, 40)
		coords := strings.Split(linestr, " -> ")
		for _, coord := range coords {
			strx, stry, ok := strings.Cut(coord, ",")
			if !ok {
				panic(fmt.Sprintf("No , in coord %s", coord))
			}
			var c Coord
			var err error
			c.x, err = strconv.Atoi(strx)
			if err != nil {
				panic(err)
			}
			c.y, err = strconv.Atoi(stry)
			if err != nil {
				panic(err)
			}
			if c.x < minx {
				minx = c.x
			}
			if c.x > maxx {
				maxx = c.x
			}
			if c.y > maxy {
				maxy = c.y
			}
			newline = append(newline, c)
		}
		lines = append(lines, newline)
	}
	// Make sure enough room around the edges is free
	minx -= 1
	maxx += 1
	if pyramid_top != 0 {
		// assume the entire cave is going to fill up like a pyramid, and add 2 y lines
		maxy += 2
		// calculate new "worst case" minx, maxx around pyramid top, and add 1 more for good measure
		wc_minx := pyramid_top - maxy - 1
		wc_maxx := pyramid_top + maxy + 1
		if wc_minx < minx {
			minx = wc_minx
		}
		if wc_maxx > maxx {
			maxx = wc_maxx
		}
	}

	lineparsetime := time.Now()
	var field Field
	field.xoffset = minx
	field.yoffset = 0
	field.xsize = maxx - minx + 1
	field.ysize = maxy + 1
	field.space = make([][]Item, field.ysize)
	for y := range field.space {
		field.space[y] = make([]Item, field.xsize)
	}
	alloctime := time.Now()
	draw_rocks(&field, lines)
	field.grains = 0
	drawtime := time.Now()
	fmt.Printf("Reading to lines took: %s, parsing lines took: %s, alloc field took: %s, drawing rocks took: %s\n",
		linesplittime.Sub(starttime), lineparsetime.Sub(linesplittime), alloctime.Sub(lineparsetime), drawtime.Sub(alloctime))
	return field
}

func draw_rocks(field *Field, lines []Line) {
	for _, l := range lines {
		pos := l[0]
		field.space[pos.y-field.yoffset][pos.x-field.xoffset] = Rock
		for i := 1; i < len(l); i++ {
			newpos := l[i]
			dx := sign(newpos.x - pos.x)
			dy := sign(newpos.y - pos.y)
			for pos.x != newpos.x || pos.y != newpos.y {
				pos.x += dx
				pos.y += dy
				field.space[pos.y-field.yoffset][pos.x-field.xoffset] = Rock
			}
		}
	}
}

func sign(a int) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	} else {
		return 0
	}
}

// returns false if the sand falls through the field
func drop_sand(f *Field, s Coord) {
	s.x -= f.xoffset
	s.y -= f.yoffset
	drop_sand_off(f, s)
}

func drop_sand_off(f *Field, s Coord) {
	// if this point is already occupied, continue
	if f.space[s.y][s.x] != Air {
		return
	}

	// did we reach the part1 dept?
	if f.part1 == 0 && s.y == f.ysize-2 {
		f.part1 = f.grains
	}

	drop_sand_off(f, Coord{s.x, s.y + 1})     // try dropping below
	drop_sand_off(f, Coord{s.x - 1, s.y + 1}) // try dropping to the left
	drop_sand_off(f, Coord{s.x + 1, s.y + 1}) // try dropping to the right

	f.space[s.y][s.x] = Sand
	f.grains++
}

func show_field(f Field) {
	for _, l := range f.space {
		for _, i := range l {
			switch i {
			case Air:
				fmt.Print(".")
			case Rock:
				fmt.Print("#")
			case Sand:
				fmt.Print("o")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	field := parse_input(os.Args[1], 500)
	// draw the extra bottom line
	bottom := []Coord{{field.xoffset, field.yoffset + field.ysize - 1}, {field.xoffset + field.xsize - 1, field.yoffset + field.ysize - 1}}
	bottomlines := []Line{bottom}
	draw_rocks(&field, bottomlines)
	parsetime := time.Now()
	drop_sand(&field, Coord{500, 0})
	part2time := time.Now()
	show_field(field)

	fmt.Printf("Number of sand part 1: %d\n", field.part1)
	fmt.Printf("Number of sand part 2: %d\n", field.grains)
	fmt.Printf("took: parsing=%s pouring sand=%s\n", parsetime.Sub(starttime), part2time.Sub(parsetime))
}
