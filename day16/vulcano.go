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
	pressure int
	timeleft [2]int
	path     [2]Path
	is_open  []bool
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
	seen[name] = true
	for len(walkers) > 0 {
		w := walkers[0]
		walkers = walkers[1:]
		newdist := w.dist + 1
		for _, to := range vl.valves[w.pos].tunnel {
			_, beenthere := seen[to]
			if beenthere {
				// we've already been here, skip
				continue
			}
			seen[to] = true
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

// Structure needed to sort valves with their distance
type ValveDist struct {
	flowrate, dist int
}

type ValveDists []ValveDist

func (v ValveDists) Len() int { return len(v) }

func (v ValveDists) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v ValveDists) Less(i, j int) bool { return v[j].flowrate < v[i].flowrate }

func max_extra_pressure(sol Solution, rv ReducedVulcano) int {
	valvedists := make(ValveDists, 0, len(rv.valves))
	pos := sol.path[len(sol.path)-1] &^ opened
	// we can assume current position has tunnels to all relevant valves
	for _, tun := range rv.valves[pos].tunnel {
		if !sol.is_open[tun.valve] && rv.valves[tun.valve].flowrate > 0 {
			valvedists = append(valvedists, ValveDist{flowrate: rv.valves[tun.valve].flowrate, dist: tun.dist})
		}
	}
	sort.Sort(valvedists)

	pressure := 0
	// assume the theoretical optimal that each path to the best valve is this projected distance
	timeleft := sol.timeleft
	for _, vd := range valvedists {
		// it takes 1 minute to open a valve
		if timeleft <= vd.dist {
			// just skip this valve if it is too far away
			continue
		}
		timeleft--
		if timeleft <= 0 {
			return pressure
		}
		pressure += vd.flowrate * (timeleft - vd.dist)
	}
	return pressure
}

// Structure needed to sort valves with their distance.
type ValveNr struct {
	flowrate int
	valvenr  Valvenr
}

type ValveNrs []ValveNr

func (v ValveNrs) Len() int { return len(v) }

func (v ValveNrs) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v ValveNrs) Less(i, j int) bool { return v[j].flowrate < v[i].flowrate }

func max_extra_pressure2(sol DuoSolution, rv ReducedVulcano) int {
	// sort the valves that need to be opened based on valvenr
	valvenrs := make(ValveNrs, 0, len(rv.valves))
	for vnr, valve := range rv.valves {
		if !sol.is_open[vnr] && valve.flowrate > 0 {
			valvenrs = append(valvenrs, ValveNr{flowrate: valve.flowrate, valvenr: Valvenr(vnr)})
		}
	}
	sort.Sort(valvenrs)

	pressure := 0
	timeleft := make([]int, 2)
	copy(timeleft, sol.timeleft[:])
	var pos [2]Valvenr
	pos[0] = sol.path[0][len(sol.path[0])-1] &^ opened
	pos[1] = sol.path[1][len(sol.path[1])-1] &^ opened
	for _, vnr := range valvenrs {
		// take the one which is closest. Or rather, with the most time left when going to the valve
		which := -1
		maxtimeleft := -1
		for try := 0; try <= 1; try++ {
			var this_dist int
			if pos[try] == vnr.valvenr {
				this_dist = 0
			} else {
				for _, tun := range rv.valves[pos[try]].tunnel {
					if tun.valve == vnr.valvenr {
						this_dist = tun.dist
						break
					}
				}
			}
			if timeleft[try]-this_dist > maxtimeleft {
				maxtimeleft = timeleft[try] - this_dist
				which = try
			}
		}
		if which == -1 || maxtimeleft <= 0 {
			// nothing found in range, just skip this valve
			continue
		}
		timeleft[which]--
		if timeleft[which] <= 0 {
			continue
		}
		pressure += vnr.flowrate * (maxtimeleft - 1)
	}
	return pressure
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

func solution_str2(rv ReducedVulcano, s DuoSolution) string {
	var symbolicpaths [2][]string
	symbolicpaths[0] = make([]string, 0, len(s.path[0]))
	symbolicpaths[1] = make([]string, 0, len(s.path[1]))
	for which, subp := range s.path {
		for i, p := range subp {
			valvenr := p &^ opened
			if p&opened == opened {
				symbolicpaths[which] = append(symbolicpaths[which], fmt.Sprintf("Open-%s", rv.valves[valvenr].name))
			} else {
				if i > 0 {
					for _, tun := range rv.valves[subp[i-1]&^opened].tunnel {
						if tun.valve == p {
							if tun.dist > 1 {
								symbolicpaths[which] = append(symbolicpaths[which], fmt.Sprintf("(%d)", tun.dist-1))
							}
							break
						}
					}
				}
				symbolicpaths[which] = append(symbolicpaths[which], rv.valves[valvenr].name)
			}
		}
	}
	return fmt.Sprintf("pressure=%d, path1=[%s] path2=[%s] timeleft1=%d timeleft2=%d",
		s.pressure, strings.Join(symbolicpaths[0], " "), strings.Join(symbolicpaths[1], " "), s.timeleft[0], s.timeleft[1])
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
	pos := candidate.path[len(candidate.path)-1] &^ opened
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

func possible_next_steps2(candidate DuoSolution, rv ReducedVulcano, best *DuoSolution) []DuoSolution {
	// take the one that has the most time left, and step that one
	var which int
	if candidate.timeleft[0] >= candidate.timeleft[1] {
		which = 0
	} else {
		which = 1
	}
	if candidate.timeleft[which] <= 0 {
		return nil
	}
	// calculate maximum pressure we could achieve by opening all remaining valves in order
	max_possible := candidate.pressure + max_extra_pressure2(candidate, rv)
	if max_possible < best.pressure {
		// no point continuing with this solution
		return nil
	}
	pos := candidate.path[which][len(candidate.path[which])-1] &^ opened
	valve := rv.valves[pos]
	new_solutions := make([]DuoSolution, 0, 10)
	// XXX again, opening current valve is not implemented
	if valve.flowrate > 0 && len(candidate.path[which]) == 1 {
		panic("Also not implemented")
	}
	for _, tunnel := range valve.tunnel {
		remote := tunnel.valve
		// to prevent symmetric identical solutions, walker #1 is limited in the first steps by what walker #0 does
		if which == 1 && len(candidate.path[1]) == 1 && remote < candidate.path[0][1]&^opened {
			continue
		}
		if rv.valves[remote].flowrate > 0 && !candidate.is_open[remote] {
			new_path := make(Path, len(candidate.path[which]), len(candidate.path[which])+2)
			copy(new_path, candidate.path[which])
			// go there and open it
			new_path = append(new_path, remote)
			new_path = append(new_path, remote|opened)
			new_open := make([]bool, len(candidate.is_open))
			copy(new_open, candidate.is_open)
			new_open[remote] = true
			rvalve := rv.valves[remote]
			var timeleft [2]int
			timeleft[which] = candidate.timeleft[which] - tunnel.dist - 1
			timeleft[1-which] = candidate.timeleft[1-which]
			var paths [2]Path
			paths[which] = new_path
			paths[1-which] = candidate.path[1-which]
			new_solutions = append(new_solutions, DuoSolution{
				pressure: candidate.pressure + rvalve.flowrate*timeleft[which],
				timeleft: timeleft,
				path:     paths,
				is_open:  new_open,
			})
		}
	}
	return new_solutions
}

// find max flow using depth-first parallel search, expaning on the best path first
func findmaxflow1(rv ReducedVulcano, start string, initial_timeleft int) (int, string) {
	// collect all possible partial solutions here, sorted by pressure
	partial_solutions := make([]Solution, 0, 20)
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
	fmt.Printf("max solutions held: %d\n", cap(partial_solutions))
	return best_solution.pressure, solution_str(rv, best_solution)
}

func findmaxflow2(rv ReducedVulcano, start string, initial_timeleft int) (int, string) {
	// collect possible partial solutions, sorted by pressure
	partial_solutions := make([]DuoSolution, 0, 20)
	partial_solutions = append(partial_solutions, DuoSolution{
		pressure: 0,
		timeleft: [2]int{initial_timeleft, initial_timeleft},
		path:     [2]Path{[]Valvenr{rv.valvenr[start]}, []Valvenr{rv.valvenr[start]}},
		is_open:  make([]bool, len(rv.valves)),
	})
	// store the best solution here
	best_solution := partial_solutions[0]
	for len(partial_solutions) > 0 {
		candidate := partial_solutions[len(partial_solutions)-1]
		partial_solutions = partial_solutions[:len(partial_solutions)-1]
		new_solutions := possible_next_steps2(candidate, rv, &best_solution)
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
			partial_solutions = append(partial_solutions, DuoSolution{})
			copy(partial_solutions[newpos+1:], partial_solutions[newpos:])
			partial_solutions[newpos] = news
		}
		if len(partial_solutions) > 0 && partial_solutions[len(partial_solutions)-1].pressure > best_solution.pressure {
			best_solution = partial_solutions[len(partial_solutions)-1]
		}
	}
	fmt.Printf("max solutions held: %d\n", cap(partial_solutions))
	return best_solution.pressure, solution_str2(rv, best_solution)
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
	flow2, solution2 := findmaxflow2(rvulcano, "AA", 26)
	part2time := time.Now()
	fmt.Printf("maxflow with elephant pressure=%d: %s\n", flow2, solution2)
	fmt.Printf("part2 took: %s\n", part2time.Sub(part1time))
}
