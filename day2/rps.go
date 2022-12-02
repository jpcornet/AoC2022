package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

func (h Hand) Score() int {
	return int(h) - int(Rock) + 1
}

func (h Hand) Beats(o Hand) bool {
	return int(h) == (int(o)+1)%3
}

var lookuphand = map[byte]Hand{
	'A': Rock,
	'B': Paper,
	'C': Scissors,
	'X': Rock,
	'Y': Paper,
	'Z': Scissors,
}

func part1(inputlines []string) {
	starttime := time.Now()
	totalscore := 0
	for idx, line := range inputlines {
		if len(line) != 3 {
			if idx != len(inputlines)-1 {
				log.Fatal(fmt.Sprintf("Invalid line [%s] at line %d", line, idx))
			}
			continue
		}
		if line[1] != ' ' {
			log.Fatal(fmt.Sprintf("No space separator in [%s] at line %d", line, idx))
		}
		theirhand, ok := lookuphand[line[0]]
		if !ok {
			log.Fatal(fmt.Sprintf("Invalid item %c at line %d", line[0], idx))
		}
		ourhand, ok := lookuphand[line[2]]
		if !ok {
			log.Fatal(fmt.Sprintf("Invalid item %c at line %d", line[2], idx))
		}
		score := ourhand.Score()
		if theirhand == ourhand {
			score += 3
		} else if ourhand.Beats(theirhand) {
			score += 6
		}
		//fmt.Printf("Line [%s] they have %#v we have %#v scores %d\n", line, theirhand, ourhand, score)
		totalscore += score
	}
	elapsed := time.Now().Sub(starttime)
	fmt.Printf("Total score: %d\n", totalscore)
	fmt.Printf("Elapsed: %s\n", elapsed)
}

func part2(inputlines []string) {
	fmt.Printf("todo\n")
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Provide input file")
	}
	// read input file
	inputFile := os.Args[1]
	inputstr, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	inputlines := strings.Split(string(inputstr[:]), "\n")
	part1(inputlines)
	part2(inputlines)
}
