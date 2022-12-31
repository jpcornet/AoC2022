package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Blueprint struct {
	nr             int
	ore_cost       int
	clay_cost      int
	obs_ore_cost   int
	obs_clay_cost  int
	geode_ore_cost int
	geode_obs_cost int
	//max_ore_robot  int // not sure if this is needed
}

type State struct {
	timeleft    int
	ore         int
	clay        int
	obsidian    int
	geode       int
	ore_robot   int
	clay_robot  int
	obs_robot   int
	geode_robot int
}

func parse_input(filename string) []Blueprint {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	bp_re := regexp.MustCompile(`^Blueprint (\d+):\s+` +
		`Each ore robot costs (\d+) ore.\s+` +
		`Each clay robot costs (\d+) ore.\s+` +
		`Each obsidian robot costs (\d+) ore and (\d+) clay.\s+` +
		`Each geode robot costs (\d+) ore and (\d+) obsidian.\s+`)
	var results []Blueprint
	for len(inbuf) > 0 {
		match := bp_re.FindSubmatchIndex(inbuf)
		if match == nil {
			panic(fmt.Sprintf("Cannot parse input at: [%s...]\n", string(inbuf[:40])))
		}
		if match[0] != 0 {
			panic(fmt.Sprintf("Logic error, expected match start at 0, not at %d\n", match[0]))
		}
		var bp Blueprint
		// Assume Atoi goes fine, the regex already makes sure it looks numeric
		bp.nr, _ = strconv.Atoi(string(inbuf[match[2]:match[3]]))
		bp.ore_cost, _ = strconv.Atoi(string(inbuf[match[4]:match[5]]))
		bp.clay_cost, _ = strconv.Atoi(string(inbuf[match[6]:match[7]]))
		bp.obs_ore_cost, _ = strconv.Atoi(string(inbuf[match[8]:match[9]]))
		bp.obs_clay_cost, _ = strconv.Atoi(string(inbuf[match[10]:match[11]]))
		bp.geode_ore_cost, _ = strconv.Atoi(string(inbuf[match[12]:match[13]]))
		bp.geode_obs_cost, _ = strconv.Atoi(string(inbuf[match[14]:match[15]]))
		// skip the matched part in the buffer
		inbuf = inbuf[match[1]:]
		results = append(results, bp)
	}
	return results
}

// Path is a collection of states. Current state is the final one.
type Path struct {
	states []State
	score  int
}

// return the minimum of 2 integers
func intmin(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func (p Path) Score(bp Blueprint) int {
	if p.score != 0 {
		return p.score
	}
	s := p.states[len(p.states)-1]
	// include all the materials that the robots will build
	// geodes are obviously best, score them good.
	p.score = (s.geode + s.geode_robot*s.timeleft) * 10_000_000_000
	// geode cracking robots are second best, score according to how many we could have built
	p.score += intmin((s.ore+s.ore_robot*s.timeleft)*10_000_000/bp.geode_ore_cost, (s.obsidian+s.obs_robot*s.timeleft)*10_000_000/bp.geode_obs_cost)
	// obsidian collecting robots are third, again score to how many could have been built
	p.score += intmin((s.ore+s.ore_robot*s.timeleft)*10_000/bp.obs_ore_cost, (s.clay+s.clay_robot*s.timeleft)*10_000/bp.obs_clay_cost)
	// clay collecting robots are next, scoring even less
	p.score += 10 * (s.ore + s.ore_robot*s.timeleft) / bp.clay_cost
	return p.score
}

// an original state and the extra robots we will build here
type BuildState struct {
	s                                             State
	ore_robot, clay_robot, obs_robot, geode_robot int
}

func possible_next_steps(p Path, bp Blueprint) []Path {
	state := p.states[len(p.states)-1]
	// collect possible next steps, as BuildStates
	var next_step []BuildState
	// try adding all robot types, 0 or 1 ore robot because not adding robots is also an option
	for ore_robot := 0; ore_robot <= 1 && state.ore >= ore_robot*bp.ore_cost; ore_robot++ {
		new_state := BuildState{s: state, ore_robot: ore_robot}
		new_state.s.ore -= ore_robot * bp.ore_cost
		next_step = append(next_step, new_state)
	}
	// try adding a clay robot
	if state.ore >= bp.clay_cost {
		new_state := BuildState{s: state, clay_robot: 1}
		new_state.s.ore -= bp.clay_cost
		next_step = append(next_step, new_state)
	}
	// try adding an obsidian robot
	if state.ore >= bp.obs_ore_cost && state.clay >= bp.obs_clay_cost {
		new_state := BuildState{s: state, obs_robot: 1}
		new_state.s.ore -= bp.obs_ore_cost
		new_state.s.clay -= bp.obs_clay_cost
		next_step = append(next_step, new_state)
	}
	// try adding a geode cracking robot
	if state.ore >= bp.geode_ore_cost && state.obsidian >= bp.geode_obs_cost {
		new_state := BuildState{s: state, geode_robot: 1}
		new_state.s.ore -= bp.geode_ore_cost
		new_state.s.obsidian -= bp.geode_obs_cost
		next_step = append(next_step, new_state)
	}

	// now we actually go and build the updated paths
	next_path := make([]Path, len(next_step))
	for i, ns := range next_step {
		next_path[i].states = make([]State, len(p.states)+1)
		copy(next_path[i].states, p.states)
		// advance the state by one minute
		ns.s.timeleft--
		// materials are being delivered by the robots
		ns.s.ore += ns.s.ore_robot
		ns.s.clay += ns.s.clay_robot
		ns.s.obsidian += ns.s.obs_robot
		ns.s.geode += ns.s.geode_robot
		// and the new robots are delivered
		ns.s.ore_robot += ns.ore_robot
		ns.s.clay_robot += ns.clay_robot
		ns.s.obs_robot += ns.obs_robot
		ns.s.geode_robot += ns.geode_robot
		// store this new state at the end of the new path
		next_path[i].states[len(p.states)] = ns.s
	}
	return next_path
}

func get_max_geodes(bp Blueprint, state State, maxsolutions int) Path {
	// keep a list of possible states here
	solutions := make([]Path, 1)
	solutions[0] = Path{states: []State{state}}
	seen := make(map[State]bool)
	seen_hits := 0
	seen_miss := 0
	for len(solutions) > 0 {
		fmt.Printf("timeleft=%d, considering %d possible solutions\n", solutions[0].states[len(solutions[0].states)-1].timeleft, len(solutions))
		new_solutions := make([]Path, 0, len(solutions))
		for _, s := range solutions {
			next_step := possible_next_steps(s, bp)
			for _, ns := range next_step {
				// check if we haven't seen this state before, skip if we have
				// We check the state without the timeleft field.
				verify_state := ns.states[len(ns.states)-1]
				verify_state.timeleft = 0
				if _, already_seen := seen[verify_state]; already_seen {
					seen_hits++
					continue
				} else {
					seen[verify_state] = true
					seen_miss++
				}
				position, _ := sort.Find(len(new_solutions), func(i int) int {
					if ns.Score(bp) < new_solutions[i].Score(bp) {
						return 1
					} else if ns.Score(bp) > new_solutions[i].Score(bp) {
						return -1
					} else {
						return 0
					}
				})
				new_solutions = append(new_solutions, Path{})
				copy(new_solutions[position+1:], new_solutions[position:])
				new_solutions[position] = ns
			}
		}
		//fmt.Printf("seen cache hits=%d miss=%d\n", seen_hits, seen_miss)
		// iterate over new_solutions, and drop any that have fewer or same of everything
		solutions = make([]Path, 0)
		dropped := 0
	inspect_new_solutions:
		for _, ns := range new_solutions {
			new_final := ns.states[len(ns.states)-1]
			for _, s := range solutions {
				sol_final := s.states[len(s.states)-1]
				if sol_final.ore >= new_final.ore && sol_final.clay >= new_final.clay && sol_final.obsidian >= new_final.obsidian &&
					sol_final.geode >= new_final.geode && sol_final.ore_robot >= new_final.ore_robot && sol_final.clay_robot >= new_final.clay_robot &&
					sol_final.obs_robot >= new_final.obs_robot && sol_final.geode_robot >= new_final.geode_robot {
					dropped++
					continue inspect_new_solutions
				}
			}
			solutions = append(solutions, ns)
			if len(solutions) >= maxsolutions {
				// only take the top scoring solutions
				fmt.Printf("Cutting off solutions at %d. Bottom score=%d\n", maxsolutions, ns.Score(bp))
				break
			}
		}
		fmt.Printf("Dropped %d inferior solutions\n", dropped)
		if len(solutions) == 0 {
			fmt.Printf("No solutions found?\n")
			return Path{}
		}
		if solutions[0].states[len(solutions[0].states)-1].timeleft == 0 {
			break
		}
	}
	// the first one is always the best solution, since they are ordered
	return solutions[0]
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	blueprints := parse_input(os.Args[1])
	parsetime := time.Now()
	total_quality := 0
	for _, bp := range blueprints {
		var state State
		state.ore_robot = 1
		state.timeleft = 24
		// experimentally, consider max 150 solutions is more than enough
		result := get_max_geodes(bp, state, 150)
		final := result.states[len(result.states)-1]
		fmt.Printf("Blueprint #%d produces max %d geodes, with: %v\n", bp.nr, final.geode, result)
		total_quality += bp.nr * final.geode
	}
	part1time := time.Now()
	fmt.Printf("Total quality: %d\n", total_quality)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", part1time.Sub(parsetime))
	multiple := 1
	for _, bp := range blueprints[:intmin(3, len(blueprints))] {
		var state State
		state.ore_robot = 1
		state.timeleft = 32
		result := get_max_geodes(bp, state, 150)
		final := result.states[len(result.states)-1]
		fmt.Printf("Blueprint #%d produces max %d geodes, with: %v\n", bp.nr, final.geode, result)
		multiple *= final.geode
	}
	part2time := time.Now()
	fmt.Printf("Multiplied number of geodes: %d\n", multiple)
	fmt.Printf("part 2 took: %s\n", part2time.Sub(part1time))
}
