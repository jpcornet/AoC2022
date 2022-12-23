package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type GeodesResult struct {
	state State
	path  []string
}

var cache_hit, cache_miss int

func prep_get_max_geodes(bp Blueprint) func(State) GeodesResult {
	cache := make(map[State]GeodesResult)
	var call_getmax func(State) GeodesResult
	call_getmax = func(state State) GeodesResult {
		if ret, ok := cache[state]; ok {
			cache_hit++
			return ret
		}
		cache_miss++
		value := raw_get_max_geodes(bp, state, call_getmax)
		cache[state] = value
		return value
	}
	return call_getmax
}

func raw_get_max_geodes(bp Blueprint, state State, recurse func(State) GeodesResult) GeodesResult {
	if state.timeleft <= 0 {
		// no time left, so just return the input state
		return GeodesResult{state: state, path: nil}
	}
	// calculate the stuff all robots produce this minute
	ore_produced := state.ore_robot
	clay_produced := state.clay_robot
	obs_produced := state.obs_robot
	geode_cracked := state.geode_robot
	// keep a list of current and new states. We keep the current state for easier looping over states
	newstates := make([]State, 1, 5)
	newstates[0] = state
	// for each new state, also keep track of what we do as a path
	path := make([]string, 1, 5)
	path[0] = ""
	// first, try constructing one or more geode cracking robots
	gcr := 1
	for state.ore >= gcr*bp.geode_ore_cost && state.obsidian >= gcr*bp.geode_obs_cost {
		state2 := state
		state2.ore -= gcr * bp.geode_ore_cost
		state2.obsidian -= gcr * bp.geode_obs_cost
		state2.geode_robot += gcr
		newstates = append(newstates, state2)
		path = append(path, fmt.Sprintf("Spend %d ore and %d obsidian to start building %d geode-cracking robot.\n", gcr*bp.geode_ore_cost, gcr*bp.geode_obs_cost, gcr))
		gcr++
	}
	// Try constructing one or more obsidian collecting robots
	// Do this for both the current state and any states with geode cracking robots (if we have any materials left)
	obsrobots := make([]State, 0, 1)
	for _, st := range newstates {
		obs := 1
		for st.ore >= obs*bp.obs_ore_cost && st.clay >= obs*bp.obs_clay_cost {
			state2 := st
			state2.ore -= obs * bp.obs_ore_cost
			state2.clay -= obs * bp.obs_clay_cost
			state2.obs_robot += obs
			obsrobots = append(obsrobots, state2)
			path = append(path, fmt.Sprintf("Spend %d ore and %d clay to start building %d obsidian-collecting robot.\n", obs*bp.obs_ore_cost, obs*bp.obs_clay_cost, obs))
			obs++
		}
	}
	// and append obsrobots to newstates
	for _, or := range obsrobots {
		newstates = append(newstates, or)
	}
	// try constructing one or more clay-collecting robots.
	clayrobots := make([]State, 0, 1)
	for _, st := range newstates {
		clay := 1
		for st.ore >= clay*bp.clay_cost {
			state2 := st
			state2.ore -= clay * bp.clay_cost
			state2.clay_robot += clay
			clayrobots = append(clayrobots, state2)
			path = append(path, fmt.Sprintf("Spend %d ore to start building %d clay-collecting robot.\n", clay*bp.clay_cost, clay))
			clay++
		}
	}
	// append clayrobots to newstates
	for _, cl := range clayrobots {
		newstates = append(newstates, cl)
	}
	// try constructing one or more ore-collecting robots.
	orerobots := make([]State, 0, 1)
	for _, st := range newstates {
		ore := 1
		for st.ore >= ore*bp.ore_cost {
			state2 := st
			state2.ore -= ore * bp.ore_cost
			state2.ore_robot += ore
			orerobots = append(orerobots, state2)
			path = append(path, fmt.Sprintf("Spend %d ore to start building %d ore-collecting robot.\n", ore*bp.ore_cost, ore))
			ore++
		}
	}
	// append ore robots to newstates
	for _, or := range orerobots {
		newstates = append(newstates, or)
	}
	// move initial state last
	newstates = append(newstates[1:], state)
	path = append(path[1:], path[0])
	best := state
	var bestpath []string
	// update the materials we have based on the robots we have, and recurse
	for i, st := range newstates {
		st.ore += ore_produced
		if ore_produced > 0 {
			path[i] += fmt.Sprintf("ore robots collect %d ore. You now have %d ore.\n", ore_produced, st.ore)
		}
		st.clay += clay_produced
		if clay_produced > 0 {
			path[i] += fmt.Sprintf("clay robots collect %d clay. You now have %d clay.\n", clay_produced, st.clay)
		}
		st.obsidian += obs_produced
		if obs_produced > 0 {
			path[i] += fmt.Sprintf("obsidian robots collect %d obsidian. You now have %d obsidian.\n", obs_produced, st.obsidian)
		}
		st.geode += geode_cracked
		if geode_cracked > 0 {
			path[i] += fmt.Sprintf("geode cracking robots crack %d geodes. You now have %d open geodes.\n", geode_cracked, st.geode)
		}
		st.timeleft--
		result := recurse(st)
		if result.state.geode > best.geode {
			best = result.state
			bestpath = append(result.path, path[i])
			//fmt.Printf("Best solution at timeleft=%d has %d geodes. state=%v\n", best.timeleft, best.geode, best)
		}
	}
	return GeodesResult{state: best, path: bestpath}
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
		// init caching recursive function for this blueprint
		get_max_geodes := prep_get_max_geodes(bp)
		// calculate max geodes for this blueprint
		var state State
		state.ore_robot = 1
		state.timeleft = 24
		result := get_max_geodes(state)
		fmt.Printf("Blueprint #%d produces max %d geodes, with: %s\n", bp.nr, result.state.geode, strings.Join(result.path, "===\n"))
		fmt.Printf("Cache hits=%d miss=%d\n", cache_hit, cache_miss)
		total_quality += bp.nr * result.state.geode
	}
	part1time := time.Now()
	fmt.Printf("Total quality: %d\n", total_quality)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", part1time.Sub(parsetime))
}
