package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Sensor struct {
	x       int
	y       int
	beaconx int
	beacony int
	dist    int
}

func intabs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func parse_input(filename string) []Sensor {
	inputstr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	// remove last empty line
	if len(inputlines[len(inputlines)-1]) == 0 {
		inputlines = inputlines[0 : len(inputlines)-1]
	}
	sensors := make([]Sensor, 0, 500)
	for _, linestr := range inputlines {
		var err error
		var sx, sy, bx, by int
		startstring := "Sensor at x="
		if !strings.HasPrefix(linestr, startstring) {
			panic("Invalid input 1")
		}
		num_start := len(startstring)
		sep_at := strings.IndexRune(linestr[num_start:], ',')
		if sep_at == -1 {
			panic("Invalid input 1b")
		}
		sx, err = strconv.Atoi(linestr[num_start : num_start+sep_at])
		if err != nil {
			panic(err)
		}

		commaystring := ", y="
		num_start += sep_at
		if !strings.HasPrefix(linestr[num_start:], ", y=") {
			panic("Invalid input 2")
		}
		num_start += len(commaystring)
		sep_at = strings.IndexRune(linestr[num_start:], ':')
		if sep_at == -1 {
			panic("Invalid input 2b")
		}
		sy, err = strconv.Atoi(linestr[num_start : num_start+sep_at])
		if err != nil {
			panic(err)
		}

		beaconstring := ": closest beacon is at x="
		num_start += sep_at
		if !strings.HasPrefix(linestr[num_start:], beaconstring) {
			panic("Invalid input 3")
		}
		num_start += len(beaconstring)
		sep_at = strings.IndexRune(linestr[num_start:], ',')
		if sep_at == -1 {
			panic("Invalid input 3b")
		}
		bx, err = strconv.Atoi(linestr[num_start : num_start+sep_at])
		if err != nil {
			panic(err)
		}

		num_start += sep_at
		if !strings.HasPrefix(linestr[num_start:], commaystring) {
			panic("Invalid input 4")
		}
		num_start += len(commaystring)
		by, err = strconv.Atoi(linestr[num_start:])
		if err != nil {
			panic(err)
		}

		sensors = append(sensors, Sensor{x: sx, y: sy, beaconx: bx, beacony: by, dist: intabs(sx-bx) + intabs(sy-by)})
		//fmt.Printf("Parsed sensor: %v\n", sensors[len(sensors)-1])
	}
	return sensors
}

// Only return sensors that are in range of the given y coord
func filter_sensors(sensors []Sensor, y int) []Sensor {
	sfilter := make([]Sensor, 0, 500)
	for _, s := range sensors {
		if intabs(s.y-y) < s.dist {
			sfilter = append(sfilter, s)
		}
	}
	return sfilter
}

func leftmostsensor(sensors []Sensor, x int, y int) (*Sensor, int, int) {
	var found *Sensor
	var influence int
	var end int
	for i, s := range sensors {
		// size of exclusion range on this line
		ssize := s.dist - intabs(s.y-y)
		// if ssize is 0 or negative, this sensor does not help. Should not happen?
		if ssize <= 0 {
			continue
		}
		// calculate maxiomium x coord on y line this sensor excludes
		my_xmax := s.x + ssize
		// if this sensor does not reach beyong x, skip it
		if my_xmax < x {
			continue
		}
		// calculate minimum x coord on y line this sensor excludes
		my_xmin := s.x - ssize
		if found == nil || my_xmin < influence {
			found = &sensors[i]
			influence = my_xmin
			end = my_xmax
		}
	}
	return found, influence, end
}

func part1(sensors []Sensor, y int) int {
	snear := filter_sensors(sensors, y)
	// count number of exclude positions
	excludepos := 0
	// find sensors from left to right with a sphere of influence starting at least at leftpos
	leftpos := math.MinInt
	for true {
		s, start, end := leftmostsensor(snear, leftpos, y)
		if s == nil {
			break
		}
		if start < leftpos {
			start = leftpos
		}
		excludepos += end - start + 1
		// is the current sensor's beacon on this line and in this range?
		if s.beacony == y && s.beaconx >= start && s.beaconx <= end {
			excludepos--
		}
		leftpos = end + 1
	}
	return excludepos
}

func scan_line(sensors []Sensor, y int, maxx int) ([]int, int) {
	snear := filter_sensors(sensors, y)
	leftpos := math.MinInt
	var gap []int
	min_overlap := math.MaxInt
	for true {
		s, start, end := leftmostsensor(snear, leftpos, y)
		if s == nil {
			break
		}
		if start <= leftpos && leftpos-start < min_overlap {
			min_overlap = leftpos - start
		}
		gapstart := leftpos
		gapend := start - 1
		if gapstart < 0 {
			gapstart = 0
		}
		if gapend > maxx {
			gapend = maxx
		}
		if gapstart <= gapend {
			gap = append(gap, gapstart, gapend)
		}
		leftpos = end + 1
	}
	if min_overlap == math.MaxInt {
		min_overlap = -1
	}
	return gap, min_overlap
}

func part2(sensors []Sensor, maxy int) int {
	for y := 0; y <= maxy; {
		line_gap, min_overlap := scan_line(sensors, y, maxy)
		if line_gap != nil {
			if line_gap[0] == line_gap[1] {
				return line_gap[0]*4000000 + y
			} else {
				fmt.Printf("On line %d, gaps at: %v\n", y, line_gap)
			}
		}
		if min_overlap >= 2 {
			// skip the next overlap/2 lines
			y += min_overlap >> 1
		} else {
			// no significant overlap, just take next line
			y++
		}
	}
	return -1
}

func main() {
	if len(os.Args) <= 1 {
		panic("Provide input file")
	}
	starttime := time.Now()
	sensors := parse_input(os.Args[1])
	y := 2000000
	if len(os.Args) == 3 {
		var err error
		y, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	}
	parsetime := time.Now()
	answ := part1(sensors, y)
	part1time := time.Now()
	answ2 := part2(sensors, y*2)
	part2time := time.Now()
	fmt.Printf("Part 1, total excluded on line %d is: %d\n", y, answ)
	fmt.Printf("Part 2, frequency of single gap: %d\n", answ2)
	fmt.Printf("Parsing took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("Part1 took: %s\n", part1time.Sub(parsetime))
	fmt.Printf("Part2 took: %s\n", part2time.Sub(part1time))
}
