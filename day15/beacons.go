package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
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
	sensor_line_re := regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)
	for _, linestr := range inputlines {
		match := sensor_line_re.FindStringSubmatchIndex(linestr)
		if match == nil {
			panic(fmt.Sprintf("Cannot parse line: %s\n", linestr))
		}
		sx, _ := strconv.Atoi(linestr[match[2]:match[3]])
		sy, _ := strconv.Atoi(linestr[match[4]:match[5]])
		bx, _ := strconv.Atoi(linestr[match[6]:match[7]])
		by, _ := strconv.Atoi(linestr[match[8]:match[9]])
		sensors = append(sensors, Sensor{x: sx, y: sy, beaconx: bx, beacony: by, dist: intabs(sx-bx) + intabs(sy-by)})
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
			fmt.Printf("Hm, skipping sensor %v on line %d\n", s, y)
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
			fmt.Printf("Best sensor so far: %v, starting at %d\n", s, my_xmin)
			found = &sensors[i]
			influence = my_xmin
			end = my_xmax
		} else {
			fmt.Printf("Sensor %v starting at %d is more to the right\n", s, my_xmin)
		}
	}
	return found, influence, end
}

func part1(sensors []Sensor, y int) {
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
		fmt.Printf("Found sensor %v after %d between %d-%d\n", s, leftpos, start, end)
		if start < leftpos {
			fmt.Printf("This sensor overlaps with previous found sensor, not counting range %d-%d\n", start, leftpos)
			start = leftpos
		}
		excludepos += end - start + 1
		// is the current sensor's beacon on this line and in this range?
		if s.beacony == y && s.beaconx >= start && s.beaconx <= end {
			fmt.Printf("  The current sensor's beacon at %d,%d is in this range\n", s.beaconx, s.beacony)
			excludepos--
		}
		leftpos = end + 1
	}
	fmt.Printf("total exclude: %d\n", excludepos)
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
	part1(sensors, y)
	fmt.Printf("Took: %s\n", parsetime.Sub(starttime))
}
