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
	valvenr map[string]Valvenr
	valves  []RValve
}

type Solution struct {
	pressure int
	timeleft int
	path     Path
	is_open  []bool
}

type DuoSolution struct {
	pressure  int
	timeleft1 int
	timeleft2 int
	path1     Path
	path2     Path
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

type TreeWalker struct {
	pos  string
	dist int
}

// map the tunnels in the original vulcano to the tunnels in the reduced vulcano.
// valves that are not in the reduced vulcano do not appear, and are only represented in the distance between valves.
// This actually maps distances to all other valves, independent of any intermediate valves
func map_tunnels(name string, vl Vulcano, rv ReducedVulcano) []TunnelElem {
	// walk the vulcano using simplified Dijkstra
	var walkers []TreeWalker = []TreeWalker{{pos: name, dist: 0}}
	seen := make(map[string]bool)
	result := make([]TunnelElem, 0, len(rv.valves))
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
			}
			// keep walking
			walkers = append(walkers, TreeWalker{pos: to, dist: newdist})
		}
	}
	return result
}

func reduce_vulcano(vl Vulcano, start string) ReducedVulcano {
	var rv ReducedVulcano
	rv.valvenr = make(map[string]Valvenr)
	rv.valves = make([]RValve, 0, 20)
	// put all valves with non-zero flowrate in the reduced vulcano, plus the start valve
	for name, v := range vl.valves {
		if v.flowrate != 0 || name == start {
			rv.valvenr[name] = Valvenr(len(rv.valves))
			rv.valves = append(rv.valves, RValve{name: name, flowrate: v.flowrate})
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
	pressures := make([]int, 0, len(rv.valves))
	for nr, v := range rv.valves {
		is_open, _ := already_open[Valvenr(nr)]
		if !is_open && v.flowrate > 0 {
			pressures = append(pressures, v.flowrate)
		}
	}
	sort.Ints(pressures)
	pressure := 0
	// assume the theoretical near-optimal that we are 1 step from the best valve, and that the next valve is only 1 step away
	timeleft := sol.timeleft - 2
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
	return fmt.Sprintf("pressure=%d, timeleft=%d, path=[%s]",
		s.pressure, s.timeleft, strings.Join(symbolicpath, " "))
}

func possible_next_steps(candidate Solution, rv ReducedVulcano, best *Solution) []Solution {
	// if we cannot open more valves, this is a final solution
	if candidate.timeleft <= 0 {
		return nil
	}
	// calculate maximum pressure we could achieve by opening all remaining valves in order
	max_possible := candidate.pressure + max_extra_pressure(candidate, rv)
	if max_possible < best.pressure {
		// no point continuing with this solution
		return nil
	}
	pos := candidate.path[len(candidate.path)-1]
	pos = pos &^ opened
	valve := rv.valves[pos]
	new_solutions := make([]Solution, 0, 10)
	// a solution is opening this valve. Except if the flowrate is zero
	if valve.flowrate > 0 && len(candidate.path) == 1 {
		panic("Cannot handle starting at a non-zero valve")
		// This is simply not implemented
	}
	// try all tunnels from this position
	for _, tunnel := range valve.tunnel {
		remote := tunnel.valve
		// no point going there unless we need to open this
		if rv.valves[remote].flowrate > 0 && !candidate.is_open[remote] {
			new_path := make(Path, len(candidate.path), len(candidate.path)+2)
			copy(new_path, candidate.path)
			// go there and open it
			new_path = append(new_path, remote)
			new_path = append(new_path, remote|opened)
			new_open := make([]bool, len(candidate.is_open))
			copy(new_open, candidate.is_open)
			new_open[remote] = true
			rvalve := rv.valves[remote]
			new_solutions = append(new_solutions, Solution{
				pressure: candidate.pressure + rvalve.flowrate*(candidate.timeleft-tunnel.dist-1),
				timeleft: candidate.timeleft - tunnel.dist - 1,
				path:     new_path,
				is_open:  new_open,
			})
		}
	}
	return new_solutions
}

// find max flow using breadth-first parallel search
func findmaxflow1(rv ReducedVulcano, start string, initial_timeleft int) (int, string) {
	// collect all possible partial solutions here, sorted by pressure
	partial_solutions := make([]Solution, 0, 20)
	// calculate total flow of all closed valve
	partial_solutions = append(partial_solutions, Solution{
		pressure: 0,
		timeleft: initial_timeleft,
		path:     []Valvenr{rv.valvenr[start]},
		is_open:  make([]bool, len(rv.valves)),
	})
	// whenever we have a better solution, store it in best_solution
	best_solution := partial_solutions[0]
	fmt.Printf("Starting with solution: %s\n", solution_str(rv, best_solution))
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
	return best_solution.pressure, solution_str(rv, best_solution)
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
	flow, solution := findmaxflow1(rvulcano, "AA", 30)
	part1time := time.Now()
	fmt.Printf("maxflow pressure=%d: %s\n", flow, solution)
	fmt.Printf("part1 took: %s\n", part1time.Sub(parsetime))
}
