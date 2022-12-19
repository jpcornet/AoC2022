package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
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

// for a simplified vulcano
type Valvenr uint8

const opened Valvenr = ^Valvenr(0) ^ (^Valvenr(0) >> 1) // 0x80

type Path []Valvenr

type TunnelElem struct {
	dist  int
	valve Valvenr
}

type RValve struct {
	name     string
	flowrate int
	tunnel   []TunnelElem
}
type ReducedVulcano struct {
	valvenr       map[string]Valvenr
	valves        []RValve
	workingvalves Valvenr
}

type Solution struct {
	pressure int
	timeleft int
	path     Path
}

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

type TreeWalker struct {
	pos  string
	dist int
}

// map the tunnels in the original vulcano to the tunnels in the reduced vulcano.
// valves that are not in the reduced vulcano do not appear, and are only represented in the distance between valves.
func map_tunnels(name string, vl Vulcano, rv ReducedVulcano) []TunnelElem {
	// walk the vulcano using simplified Dijkstra
	var walkers []TreeWalker = []TreeWalker{{pos: name, dist: 0}}
	seen := make(map[string]bool)
	result := make([]TunnelElem, 0, len(vl.valves[name].tunnel))
	for len(walkers) > 0 {
		w := walkers[0]
		walkers = walkers[1:]
		seen[w.pos] = true
		newdist := w.dist + 1
		for _, to := range vl.valves[w.pos].tunnel {
			_, beenthere := seen[to]
			if beenthere {
				// we've already been here, skip
				continue
			}
			rvalvenr, in_reduced := rv.valvenr[to]
			if in_reduced {
				// it is in the reduced tunnel valve, we found the distance
				result = append(result, TunnelElem{newdist, rvalvenr})
			} else {
				// not a target valve, keep walking
				walkers = append(walkers, TreeWalker{pos: to, dist: newdist})
			}
		}
	}
	return result
}

func reduce_vulcano(vl Vulcano, start string) ReducedVulcano {
	var rv ReducedVulcano
	rv.valvenr = make(map[string]Valvenr)
	rv.valves = make([]RValve, 0, 20)
	rv.workingvalves = 0
	// put all valves with non-zero flowrate in the reduced vulcano, plus the start valve
	for name, v := range vl.valves {
		if v.flowrate != 0 || name == start {
			rv.valvenr[name] = Valvenr(len(rv.valves))
			rv.valves = append(rv.valves, RValve{name: name, flowrate: v.flowrate})
		}
		if v.flowrate != 0 {
			rv.workingvalves++
		}
	}
	// now map the tunnels between these valves
	for vn, v := range rv.valves {
		rv.valves[vn].tunnel = map_tunnels(v.name, vl, rv)
	}
	return rv
}

func max_extra_pressure(sol Solution, rv ReducedVulcano) int {
	already_open := make(map[Valvenr]bool)
	// record valves that are already open
	for _, rv := range sol.path {
		if rv&opened == opened {
			already_open[rv&^opened] = true
		}
	}
	pressures := make([]int, 0, rv.workingvalves)
	for nr, v := range rv.valves {
		is_open, _ := already_open[Valvenr(nr)]
		if !is_open && v.flowrate > 0 {
			pressures = append(pressures, v.flowrate)
		}
	}
	sort.Ints(pressures)
	pressure := 0
	// assume the theoretical optimal that we are on the best valve, and that the next valve is only 1 step away
	timeleft := sol.timeleft - 1
	for i := len(pressures) - 1; i >= 0; i-- {
		if timeleft <= 0 {
			return pressure
		}
		pressure += pressures[i] * timeleft
		// one step to reach the next valve, 1 step to open it
		timeleft -= 2
	}
	return pressure
}

// test if valve v has been opened in path p
func is_opened(v Valvenr, p Path) bool {
	for _, item := range p {
		if item == v|opened {
			return true
		}
	}
	return false
}

func solution_str(rv ReducedVulcano, s Solution) string {
	symbolicpath := make([]string, 0, len(s.path))
	for i, p := range s.path {
		valvenr := p &^ opened
		if p&opened == opened {
			symbolicpath = append(symbolicpath, fmt.Sprintf("Open-%s", rv.valves[valvenr].name))
		} else {
			if i > 0 {
				for _, tun := range rv.valves[s.path[i-1]&^opened].tunnel {
					if tun.valve == p {
						if tun.dist > 1 {
							symbolicpath = append(symbolicpath, fmt.Sprintf("(%d)", tun.dist-1))
						}
						break
					}
				}
			}
			symbolicpath = append(symbolicpath, rv.valves[valvenr].name)
		}
	}
	return fmt.Sprintf("pressure=%d, timeleft=%d, path=[%s]", s.pressure, s.timeleft, strings.Join(symbolicpath, " "))
}

func possible_next_steps(candidate Solution, rv ReducedVulcano, best *Solution) []Solution {
	// calculate max extra pressure
	max_extra := max_extra_pressure(candidate, rv)
	// if we cannot open more valves, this is a final solution
	if max_extra == 0 {
		if candidate.pressure > best.pressure {
			*best = candidate
		}
		return nil
	}
	// calculate maximum pressure we could achieve by opening all remaining valves in order
	max_possible := candidate.pressure + max_extra
	if max_possible < best.pressure {
		// no point continuing with this solution
		return nil
	}
	pos := candidate.path[len(candidate.path)-1]
	pos = pos &^ opened
	valve := rv.valves[pos]
	new_solutions := make([]Solution, 0, 10)
	// a solution is opening this valve. Except if the flowrate is zero or we already opened it
	if valve.flowrate > 0 && !is_opened(pos, candidate.path) {
		new_path := make(Path, len(candidate.path), len(candidate.path)+1)
		copy(new_path, candidate.path)
		new_path = append(new_path, pos|opened)
		new_solutions = append(new_solutions, Solution{
			pressure: candidate.pressure + valve.flowrate*(candidate.timeleft-1),
			timeleft: candidate.timeleft - 1,
			path:     new_path,
		})
	}
	// try all tunnels from this position
	for _, tunnel := range valve.tunnel {
		remote := tunnel.valve
		// there is no point going to this valve if we've already been here without opening another one
		beenthere := false
		for i := len(candidate.path) - 1; i >= 0; i-- {
			if candidate.path[i]&^opened == remote {
				beenthere = true
				break
			}
			if candidate.path[i]&opened == opened {
				// this valve was opened, so we're not interested in previous valves
				break
			}
		}
		if !beenthere {
			new_path := make(Path, len(candidate.path), len(candidate.path)+1)
			copy(new_path, candidate.path)
			new_path = append(new_path, remote)
			new_solutions = append(new_solutions, Solution{
				pressure: candidate.pressure,
				timeleft: candidate.timeleft - tunnel.dist,
				path:     new_path,
			})
		}
	}
	return new_solutions
}

// find max flow using breadth-first parallel search
func findmaxflow1(rv ReducedVulcano, start string, initial_timeleft int) (int, Path) {
	// collect all possible partial solutions here, sorted by pressure
	partial_solutions := make([]Solution, 0, 20)
	partial_solutions = append(partial_solutions, Solution{
		pressure: 0,
		timeleft: initial_timeleft,
		path:     []Valvenr{rv.valvenr[start]},
	})
	// whenever we have a better solution, store it in best_solution
	best_solution := partial_solutions[0]
	for len(partial_solutions) > 0 {
		candidate := partial_solutions[len(partial_solutions)-1]
		partial_solutions = partial_solutions[:len(partial_solutions)-1]
		new_solutions := possible_next_steps(candidate, rv, &best_solution)

		// insert new solutions into partial solutions, in order
		for _, news := range new_solutions {
			newpos, _ := sort.Find(len(partial_solutions), func(i int) int {
				if news.pressure < partial_solutions[i].pressure {
					return -1
				} else if news.pressure > partial_solutions[i].pressure {
					return 1
				} else {
					return 0
				}
			})
			partial_solutions = append(partial_solutions, Solution{})
			copy(partial_solutions[newpos+1:], partial_solutions[newpos:])
			partial_solutions[newpos] = news
		}
		if len(partial_solutions) > 0 && partial_solutions[len(partial_solutions)-1].pressure > best_solution.pressure {
			best_solution = partial_solutions[len(partial_solutions)-1]
		}
	}
	fmt.Printf("Solution to part 1: %s\n", solution_str(rv, best_solution))
	fmt.Printf("partial solutions array: %d\n", cap(partial_solutions))
	return best_solution.pressure, best_solution.path
}

// very ugly and slow recursive solution, doing depth-first search
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
	rvulcano := reduce_vulcano(vulcano, "AA")
	parsetime := time.Now()
	fmt.Printf("Got: %v\n", rvulcano)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	flow, path := findmaxflow1(rvulcano, "AA", 30)
	part1time := time.Now()
	fmt.Printf("maxflow pressure=%d path=%v\n", flow, path)
	fmt.Printf("part1 took: %s\n", part1time.Sub(parsetime))
}
