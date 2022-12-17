package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Valve struct {
	flowrate int
	tunnel   []string
}

type Vulcano struct {
	valves  map[string]Valve
	is_open map[string]bool
}

type Visited map[string]bool

func parse_input(filename string) Vulcano {
	inputstr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	// remove last empty line
	if len(inputlines[len(inputlines)-1]) == 0 {
		inputlines = inputlines[0 : len(inputlines)-1]
	}
	valves := make(map[string]Valve)
	valveline_re := regexp.MustCompile(`^Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (\w+(?:, \w+)*)$`)
	for _, linestr := range inputlines {
		match := valveline_re.FindStringSubmatch(linestr)
		if match == nil {
			panic("Invalid input")
		}
		name := match[1]
		flowrate, err := strconv.Atoi(match[2])
		if err != nil {
			panic(err)
		}
		tunnels := strings.Split(match[3], ", ")
		fmt.Printf("Valve %s flow rate=%d, tunnels to %v\n", name, flowrate, tunnels)
		valves[name] = Valve{
			flowrate: flowrate,
			tunnel:   tunnels,
		}
	}
	return Vulcano{
		valves:  valves,
		is_open: make(map[string]bool),
	}
}

func copystrarray(in []string) []string {
	result := make([]string, 0, len(in))
	for _, s := range in {
		result = append(result, strings.Clone(s))
	}
	return result
}

func maxflow(vulcano Vulcano, pos string, timeleft int, path []string, pressure int, visited Visited) (int, []string) {
	cur_max := pressure
	cur_bestpath := path
	if timeleft <= 0 {
		// no time left, cannot do anything
		return pressure, path
	}
	// we never explicitly close valves, just reset them to current state when we're not explicitly opening it
	valve_state := vulcano.is_open[pos]
	defer func() {
		vulcano.is_open[pos] = valve_state
	}()
	// keep track of which nodes we visited between opening valves
	visited[pos] = true
	defer func() {
		visited[pos] = false
	}()
	// try going to every connected tunnel, with or without opening the valve
	for i := 0; i <= 1; i++ {
		openit := i == 0
		extratime := 0
		my_visited := visited
		my_path := append(copystrarray(path), pos)
		my_pressure := pressure
		if openit {
			if vulcano.is_open[pos] {
				// valve is already open, so we visited this node already. Try just walking past it...
				continue
			}
			if vulcano.valves[pos].flowrate == 0 {
				// no point opening this valve, flow rate is 0
				continue
			}
			// open the valve. This costs time
			vulcano.is_open[pos] = true
			my_path = append(my_path, fmt.Sprintf("Open-%s", pos))
			extratime = 1
			// calculate pressure released by opening this valve
			my_pressure = pressure + (timeleft-extratime)*vulcano.valves[pos].flowrate
			// since we opened a valve, start a new "visited" map
			my_visited = make(Visited)
			my_visited[pos] = true
			// fmt.Printf("Time left: %d. Opened valve %s, extra pressure is %d\n", timeleft-extratime, pos, my_pressure-pressure)
		} else {
			vulcano.is_open[pos] = valve_state
		}
		// if all valves are open, we cannot do anything
		all_open := true
		for name, v := range vulcano.valves {
			if v.flowrate != 0 && !vulcano.is_open[name] {
				all_open = false
				break
			}
		}
		if all_open {
			// fmt.Printf("Time left: %d. Doing nothing at valve %s\n", timeleft-extratime, pos)
			// fmt.Printf("   all open, path=%v is_open=%v\n", my_path, vulcano.is_open)
			return my_pressure, append(my_path, "nothing")
		}
		for _, nextpos := range vulcano.valves[pos].tunnel {
			if my_visited[nextpos] {
				continue
			}
			// fmt.Printf("Time left: %d. Try going to pos %s, pressure is %d\n", timeleft-1-extratime, nextpos, my_pressure)
			this_pressure, this_path := maxflow(vulcano, nextpos, timeleft-1-extratime, my_path, my_pressure, my_visited)
			if this_pressure >= cur_max {
				cur_max = this_pressure
				cur_bestpath = copystrarray(this_path)
				fmt.Printf("Time left: %d. Solution: pressure=%d, path=%v\n", timeleft-extratime, cur_max, cur_bestpath)
			}
		}
	}
	return cur_max, cur_bestpath
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	vulcano := parse_input(os.Args[1])
	parsetime := time.Now()
	fmt.Printf("Got: %v\n", vulcano)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	flow, path := maxflow(vulcano, "AA", 30, nil, 0, make(Visited))
	part1time := time.Now()
	fmt.Printf("maxflow pressure=%d path=%v\n", flow, path)
	fmt.Printf("part1 took: %s\n", part1time.Sub(parsetime))
}
