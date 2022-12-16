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

type Vulcano map[string]Valve

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
	vulcano := make(Vulcano)
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
		vulcano[name] = Valve{
			flowrate: flowrate,
			tunnel:   tunnels,
		}
	}
	return vulcano
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
}
