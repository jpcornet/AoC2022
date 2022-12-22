package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cube [3]int

func parse_input(filename string) []Cube {
	inputstr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	if inputlines[len(inputlines)-1] == "" {
		inputlines = inputlines[:len(inputlines)-1]
	}
	result := make([]Cube, 0, len(inputlines))
	for _, l := range inputlines {
		xyzstr := strings.Split(l, ",")
		if len(xyzstr) != 3 {
			panic(fmt.Sprintf("invalid input line [%s] converted to: %v\n", l, xyzstr))
		}
		var cube Cube
		for i, coord := range xyzstr {
			var err error
			cube[i], err = strconv.Atoi(coord)
			if err != nil {
				panic(err)
			}
		}
		result = append(result, cube)
	}
	return result
}

var directions = []Cube{
	{1, 0, 0},
	{-1, 0, 0},
	{0, 1, 0},
	{0, -1, 0},
	{0, 0, 1},
	{0, 0, -1},
}

func calc_surface(boulder []Cube) int {
	cubes := make(map[Cube]bool)
	for _, c := range boulder {
		cubes[c] = true
	}
	surface := 0
	for _, c := range boulder {
		for _, dir := range directions {
			testc := Cube{c[0] + dir[0], c[1] + dir[1], c[2] + dir[2]}
			_, ok := cubes[testc]
			if !ok {
				// count as a surface if there is no cube in this direction
				surface++
			}
		}
	}
	return surface
}

func check_interior(pos Cube, boulder map[Cube]bool, exterior map[Cube]bool) (bool, map[Cube]bool) {
	// keep a map of everywhere we've been
	seen := make(map[Cube]bool)
	// keep walkers
	walkers := make([]Cube, 0, 10)
	walkers = append(walkers, pos)
	seen[pos] = true
	for len(walkers) > 0 {
		w := walkers[0]
		walkers = walkers[1:]
		for _, dir := range directions {
			nextw := Cube{w[0] + dir[0], w[1] + dir[1], w[2] + dir[2]}
			_, in_boulder := boulder[nextw]
			if in_boulder {
				// we walked into the boulder
				continue
			}
			_, beenthere := seen[nextw]
			if beenthere {
				// we've already been here
				continue
			}
			_, in_ext := exterior[nextw]
			if in_ext {
				// it's connected to the exterior!
				return false, seen
			}
			seen[nextw] = true
			walkers = append(walkers, nextw)
		}
	}
	return true, seen
}

func add_interior_bubbles(boulder *[]Cube) {
	// again, need a map of cubes present
	cubes := make(map[Cube]bool)
	// also collect x, y, z ranges
	min := [3]int{math.MaxInt, math.MaxInt, math.MaxInt}
	max := [3]int{math.MinInt, math.MinInt, math.MinInt}
	for _, c := range *boulder {
		cubes[c] = true
		for i, coord := range c {
			if coord < min[i] {
				min[i] = coord
			}
			if coord > max[i] {
				max[i] = coord
			}
		}
	}
	// keep a map of all cubes that are connected to the exterior
	exterior := make(map[Cube]bool)
	// first, mark everything around the boulder as exterior
	// Top and bottom
	for x := min[0] - 1; x <= max[0]+1; x++ {
		for y := min[1] - 1; y <= max[1]+1; y++ {
			exterior[Cube{x, y, min[2] - 1}] = true
			exterior[Cube{x, y, max[2] + 1}] = true
		}
	}
	// front and back
	for x := min[0] - 1; x <= max[0]+1; x++ {
		for z := min[2] - 1; z <= max[2]+1; z++ {
			exterior[Cube{x, min[1] - 1, z}] = true
			exterior[Cube{x, max[1] + 1, z}] = true
		}
	}
	// left and right
	for y := min[1] - 1; y <= max[1]+1; y++ {
		for z := min[2] - 1; z <= max[2]+1; z++ {
			exterior[Cube{min[0] - 1, y, z}] = true
			exterior[Cube{max[0] + 1, y, z}] = true
		}
	}
	// now walk every possible cube within that area
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			for z := min[2]; z <= max[2]; z++ {
				this_c := Cube{x, y, z}
				_, isboulder := cubes[this_c]
				if isboulder {
					// part of the boulder, just continue
					continue
				}
				_, isext := exterior[this_c]
				if isext {
					// part of the exterior
					continue
				}
				is_int, seen := check_interior(this_c, cubes, exterior)
				if is_int {
					// if it is interior, add everything to the boulder itself
					for c := range seen {
						*boulder = append(*boulder, c)
						cubes[c] = true
					}
				} else {
					// if it isn't interior, it's exterior
					for c := range seen {
						exterior[c] = true
					}
				}
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	boulder := parse_input(os.Args[1])
	parsetime := time.Now()
	surface := calc_surface(boulder)
	calctime := time.Now()
	add_interior_bubbles(&boulder)
	surface2 := calc_surface(boulder)
	calctime2 := time.Now()
	fmt.Printf("Surface part1=%d\n", surface)
	fmt.Printf("Surface part2=%d\n", surface2)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("Surface calc took: %s\n", calctime.Sub(parsetime))
	fmt.Printf("part 2 took: %s\n", calctime2.Sub(calctime))
}
