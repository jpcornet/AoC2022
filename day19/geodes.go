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
	max_ore_robot  int
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

type Path struct {
	bp     *Blueprint
	states []State
	score  int
}

type Solutions []Path

// Implement the interface needed to sort Solutions based on score. Note that it sorts highest score first
//func (s Solutions) Len() int { return len(s) }
//func (s Solutions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
//func (s Solutions) Less(i, j int) bool { return s[j].Score() < s[i].Score() }

func fill_solution(bp Blueprint, state State) Path {
	var path Path
	path.bp = &bp
	path.states = make([]State, 0, state.timeleft)
	elem := state
	for elem.timeleft > 0 {
		elem.ore += elem.ore_robot
		elem.clay += elem.clay_robot
		elem.obsidian += elem.obs_robot
		elem.geode += elem.geode_robot
		path.states = append(path.states, elem)
		elem.timeleft--
	}
	return path
}

func (p Path) Copy() Path {
	var new Path
	new.bp = p.bp
	new.states = make([]State, len(p.states))
	copy(new.states, p.states)
	return new
}

// return the minimum of 2 integers
func intmin(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func (p Path) Score() int {
	if p.score != 0 {
		return p.score
	}
	final := p.states[len(p.states)-1]
	// geodes are obviously best, score them good.
	p.score = final.geode * 1_000_000_000
	// geode cracking robots are second best, score according to how many we could have built
	p.score += intmin(final.ore*1_000_000/p.bp.geode_ore_cost, final.obsidian*1_000_000/p.bp.geode_obs_cost)
	// obsidian collecting robots are third, again score to how many could have been built
	p.score += intmin(final.ore*1000/p.bp.obs_ore_cost, final.clay*1000/p.bp.obs_clay_cost)
	// clay collecting robots are next, scoring even less
	p.score += final.ore / p.bp.clay_cost
	return p.score
}

// Find the first spot where there is enough ore, clay and obsidian as specified in the minstate
// note that other state elements aren't used
func (p Path) FindMinState(startpos int, minstate State) int {
	for i, state := range p.states[startpos : len(p.states)-1] {
		if state.ore >= minstate.ore && state.clay >= minstate.clay && state.obsidian >= minstate.obsidian {
			// we might've already added a robot at this point, so make sure the materials are still available at the next step
			nextstate := p.states[i+1]
			if nextstate.ore >= minstate.ore && nextstate.clay >= minstate.clay && nextstate.obsidian >= minstate.obsidian {
				// return next state, because elements only become available at the end of the minute.
				return startpos + i + 1
			}
		}
	}
	return -1
}

// Add a robot to the solution. Returns true on success
// return false on failure (when there is not enough material), and note that state of p will be inconsistent
func (p Path) AddRobot(pos int, rstate State) bool {
	for i := pos; i < len(p.states); i++ {
		if rstate.ore_robot > 0 {
			p.states[i].ore -= p.bp.ore_cost * rstate.ore_robot
			// number of available materials cannot drop below number of robots delivering this
			if p.states[i].ore < p.states[i].ore_robot {
				return false
			}
			p.states[i].ore_robot += rstate.ore_robot
			// each robot produces 0 the minute it is built, but 1 each minute afterwards
			p.states[i].ore += rstate.ore_robot * (i - pos)
		}
		if rstate.clay_robot > 0 {
			p.states[i].ore -= p.bp.clay_cost * rstate.clay_robot
			if p.states[i].ore < p.states[i].ore_robot {
				return false
			}
			p.states[i].clay_robot += rstate.clay_robot
			p.states[i].clay += rstate.clay_robot * (i - pos)
		}
		if rstate.obs_robot > 0 {
			p.states[i].ore -= p.bp.obs_ore_cost * rstate.obs_robot
			p.states[i].clay -= p.bp.obs_clay_cost * rstate.obs_robot
			if p.states[i].ore < p.states[i].ore_robot || p.states[i].clay < p.states[i].clay_robot {
				return false
			}
			p.states[i].obs_robot += rstate.obs_robot
			p.states[i].obsidian += rstate.obs_robot * (i - pos)
		}
		if rstate.geode_robot > 0 {
			p.states[i].ore -= p.bp.geode_ore_cost * rstate.geode_robot
			p.states[i].obsidian -= p.bp.geode_obs_cost * rstate.geode_robot
			if p.states[i].ore < p.states[i].ore_robot || p.states[i].obsidian < p.states[i].obs_robot {
				return false
			}
			p.states[i].geode_robot += rstate.geode_robot
			p.states[i].geode += rstate.geode_robot * (i - pos)
		}
	}
	return true
}

func improve(c Path) Solutions {
	// collect new solutions here
	new_solutions := make(Solutions, 0, 4)
	// try adding geode-cracking robots, if possible. Start searching at the beginning
	startpos := 0
	// loop until we inserted one geode-cracking robot
	for {
		insert_geode_robot := c.FindMinState(startpos, State{ore: c.bp.geode_ore_cost, obsidian: c.bp.geode_obs_cost})
		if insert_geode_robot < 0 {
			break
		}
		solution := c.Copy()
		if solution.AddRobot(insert_geode_robot, State{geode_robot: 1}) {
			new_solutions = append(new_solutions, solution)
			break
		}
		startpos = insert_geode_robot
	}
	// get the final state to check against amount of materials we have at the end
	final := c.states[len(c.states)-1]
	// if we need more obsidian, try inserting an obsidan-collecting robot
	// we only need obsidian if the ratio between ore and obsidian is below what is needed to build the robot
	need_obsidian := false
	if final.obsidian <= final.ore*c.bp.geode_obs_cost/c.bp.geode_ore_cost {
		need_obsidian = true
		startpos = 0
		for {
			insert_obs_robot := c.FindMinState(startpos, State{ore: c.bp.obs_ore_cost, clay: c.bp.obs_clay_cost})
			if insert_obs_robot < 0 {
				break
			}
			solution := c.Copy()
			if solution.AddRobot(insert_obs_robot, State{obs_robot: 1}) {
				new_solutions = append(new_solutions, solution)
				break
			}
			startpos = insert_obs_robot
		}
	}
	// if we need more clay, try inserting a clay-collecting robot
	// we only need more clay if we need more obsidian and the ratio of clay and ore, is less than what is needed for an obsidian robot
	if need_obsidian && final.clay <= final.ore*c.bp.obs_clay_cost/c.bp.obs_ore_cost {
		startpos = 0
		for {
			insert_clay_robot := c.FindMinState(startpos, State{ore: c.bp.clay_cost})
			if insert_clay_robot < 0 {
				break
			}
			solution := c.Copy()
			if solution.AddRobot(insert_clay_robot, State{clay_robot: 1}) {
				new_solutions = append(new_solutions, solution)
				break
			}
			startpos = insert_clay_robot
		}
	}
	// only try extra ore if allowed
	if final.ore_robot < c.bp.max_ore_robot {
		startpos = 0
		for {
			insert_ore_robot := c.FindMinState(startpos, State{ore: c.bp.ore_cost})
			if insert_ore_robot < 0 {
				break
			}
			solution := c.Copy()
			if solution.AddRobot(insert_ore_robot, State{ore_robot: 1}) {
				new_solutions = append(new_solutions, solution)
				break
			}
			startpos = insert_ore_robot
		}
	}
	return new_solutions
}

func get_max_geodes(bp Blueprint, state State) Path {
	// keep the best solution so far here
	best := Path{score: -1}
	// and keep a list of possible solutions
	possible := make(Solutions, 0, 10)
	max_ore_robot := 1
	for {
		fmt.Printf("Trying with max ore robots=%d\n", max_ore_robot)
		bp.max_ore_robot = max_ore_robot
		possible = append(possible, fill_solution(bp, state))
		seen := make(map[State]bool)
		showtime := time.Now()
		for len(possible) > 0 {
			if time.Since(showtime).Seconds() > 1.0 {
				fmt.Printf("Considering %d possible solutions\n", len(possible))
				showtime = time.Now()
			}
			// take first solution from the list
			candidate := possible[0]
			possible = possible[1:]
			new_solutions := improve(candidate)
			for _, s := range new_solutions {
				if _, already_done := seen[s.states[len(s.states)-1]]; already_done {
					continue
				} else {
					seen[s.states[len(s.states)-1]] = true
				}
				if s.Score() > best.Score() {
					fmt.Printf("New best solution, score=%d: %v\n", s.Score(), s)
					best = s.Copy()
				}
				position, _ := sort.Find(len(possible), func(i int) int {
					if s.Score() < possible[i].Score() {
						return 1
					} else if s.Score() > possible[i].Score() {
						return -1
					} else {
						return 0
					}
				})
				possible = append(possible, Path{})
				copy(possible[position+1:], possible[position:])
				possible[position] = s
			}
		}
		if best.states[len(best.states)-1].ore_robot < max_ore_robot {
			break
		} else {
			// try to see if adding 1 ore robot gives a better solution
			max_ore_robot += 1
		}
	}
	return best
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
