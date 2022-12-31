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
	max_ore_robot  int // not sure if this is needed
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
	// geodes are obviously best, score them good.
	p.score = s.geode * 1_000_000_000
	// geode cracking robots are second best, score according to how many we could have built
	p.score += intmin(s.ore*1_000_000/bp.geode_ore_cost, s.obsidian*1_000_000/bp.geode_obs_cost)
	// obsidian collecting robots are third, again score to how many could have been built
	p.score += intmin(s.ore*1000/bp.obs_ore_cost, s.clay*1000/bp.obs_clay_cost)
	// clay collecting robots are next, scoring even less
	p.score += s.ore / bp.clay_cost
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
	// try adding all robot types, as much as possible
	// start with ore robots. Include adding 0 robots
	// XXX no ore robots for now
	//for ore_robot := 0; state.ore >= ore_robot*bp.ore_cost; ore_robot++ {
	for ore_robot := 0; ore_robot < 1; ore_robot++ {
		new_state := BuildState{s: state, ore_robot: ore_robot}
		new_state.s.ore -= ore_robot * bp.ore_cost
		next_step = append(next_step, new_state)
	}
	// try adding clay robots, at least one. Add to more_next_step.
	var clay_next_step []BuildState
	for clay_robot := 1; state.ore >= clay_robot*bp.clay_cost; clay_robot++ {
		// try adding to all the new states, stop if there are not enough materials
		for _, ns := range next_step {
			if ns.s.ore < clay_robot*bp.clay_cost {
				break
			}
			new_state := BuildState{s: ns.s, clay_robot: clay_robot}
			new_state.s.ore -= clay_robot * bp.clay_cost
			clay_next_step = append(clay_next_step, new_state)
		}
	}
	next_step = append(next_step, clay_next_step...)
	// try adding obsidian robots, at least one
	var obs_next_step []BuildState
	for obs_robot := 1; state.ore >= obs_robot*bp.obs_ore_cost && state.clay >= obs_robot*bp.obs_clay_cost; obs_robot++ {
		for _, ns := range next_step {
			if ns.s.ore < obs_robot*bp.obs_ore_cost || ns.s.clay < obs_robot*bp.obs_clay_cost {
				continue
			}
			new_state := BuildState{s: ns.s, obs_robot: obs_robot}
			new_state.s.ore -= obs_robot * bp.obs_ore_cost
			new_state.s.clay -= obs_robot * bp.obs_clay_cost
			obs_next_step = append(obs_next_step, new_state)
		}
	}
	next_step = append(next_step, obs_next_step...)
	// try adding geode cracking robots, at least one
	var geo_next_step []BuildState
	for geode_robot := 1; state.ore >= geode_robot*bp.geode_ore_cost && state.obsidian >= geode_robot*bp.geode_obs_cost; geode_robot++ {
		for _, ns := range next_step {
			if ns.s.ore < geode_robot*bp.geode_ore_cost || ns.s.obsidian < geode_robot*bp.geode_obs_cost {
				continue
			}
			new_state := BuildState{s: ns.s, geode_robot: geode_robot}
			new_state.s.ore -= geode_robot * bp.geode_ore_cost
			new_state.s.obsidian -= geode_robot * bp.geode_obs_cost
			geo_next_step = append(geo_next_step, new_state)
		}
	}
	next_step = append(next_step, geo_next_step...)

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

func get_max_geodes(bp Blueprint, state State) Path {
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
		solutions = new_solutions
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
		result := get_max_geodes(bp, state)
		final := result.states[len(result.states)-1]
		fmt.Printf("Blueprint #%d produces max %d geodes, with: %v\n", bp.nr, final.geode, result)
		total_quality += bp.nr * final.geode
	}
	part1time := time.Now()
	fmt.Printf("Total quality: %d\n", total_quality)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", part1time.Sub(parsetime))
}
