package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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
}

type Coord struct {
	x int
	y int
}

type Line []Coord

func parse_input(filename string, pyramid_top int) Field {
	inputstr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	var lines []Line
	minx := math.MaxInt
	maxx := math.MinInt
	maxy := math.MinInt
	for _, linestr := range inputlines {
		if len(linestr) == 0 {
			continue
		}
		var newline Line
		coords := strings.Split(linestr, " -> ")
		for _, coord := range coords {
			if len(coord) == 0 {
				continue
			}
			strxy := strings.Split(coord, ",")
			var c Coord
			var err error
			c.x, err = strconv.Atoi(strxy[0])
			if err != nil {
				panic(err)
			}
			c.y, err = strconv.Atoi(strxy[1])
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
		// calculate new "worst case" minx, maxx around pyramid top, and add 10 more for good measure
		wc_minx := pyramid_top - maxy - 10
		wc_maxx := pyramid_top + maxy + 10
		if wc_minx < minx {
			minx = wc_minx
		}
		if wc_maxx > maxx {
			maxx = wc_maxx
		}
	}

	var field Field
	field.xoffset = minx
	field.yoffset = 0
	field.xsize = maxx - minx + 1
	field.ysize = maxy + 1
	field.space = make([][]Item, field.ysize)
	for y := range field.space {
		field.space[y] = make([]Item, field.xsize)
		for x := range field.space[y] {
			field.space[y][x] = Air
		}
	}
	draw_rocks(field, lines)
	return field
}

func draw_rocks(field Field, lines []Line) {
	fmt.Printf("Drawing on field xoffset=%d, yoffset=%d, xsize=%d, ysize=%d\n", field.xoffset, field.yoffset, field.xsize, field.ysize)
	for _, l := range lines {
		fmt.Printf("Drawing line %v\n", l)
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
func drop_sand(f Field, s Coord) bool {
	s.x -= f.xoffset
	s.y -= f.yoffset
	for f.space[s.y][s.x] == Air {
		if s.y+1 == f.ysize {
			return false
		}
		if f.space[s.y+1][s.x] == Air {
			s.y += 1
		} else if f.space[s.y+1][s.x-1] == Air {
			s.y += 1
			s.x -= 1
		} else if f.space[s.y+1][s.x+1] == Air {
			s.y += 1
			s.x += 1
		} else {
			f.space[s.y][s.x] = Sand
			return true
		}
	}
	return false
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
	field := parse_input(os.Args[1], 0)
	sand := 0
	for drop_sand(field, Coord{500, 0}) {
		sand += 1
	}
	fmt.Printf("Number of sand dropped: %d\n", sand)
	show_field(field)

	// for part 2, read input again this time adding the extra space for the entire pyramid
	field = parse_input(os.Args[1], 500)
	// draw the extra bottom line
	bottom := []Coord{{field.xoffset, field.yoffset + field.ysize - 1}, {field.xoffset + field.xsize - 1, field.yoffset + field.ysize - 1}}
	bottomlines := []Line{bottom}
	draw_rocks(field, bottomlines)
	sand = 0
	for drop_sand(field, Coord{500, 0}) {
		sand += 1
	}
	show_field(field)
	fmt.Printf("Number of sand part 2: %d\n", sand)
}
