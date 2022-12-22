package main

import (
	"fmt"
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

func calc_surface(boulder []Cube) int {
	cubes := make(map[Cube]bool)
	for _, c := range boulder {
		cubes[c] = true
	}
	surface := 0
	directions := []Cube{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}
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

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	boulder := parse_input(os.Args[1])
	parsetime := time.Now()
	surface := calc_surface(boulder)
	calctime := time.Now()
	fmt.Printf("Surface=%d\n", surface)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("Surface calc took: %s\n", calctime.Sub(parsetime))
}
