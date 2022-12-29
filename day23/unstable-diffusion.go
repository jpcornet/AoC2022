package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Field [][]byte

func parse_input(filename string) Field {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	field := strings.Split(string(inbuf[:]), "\n")
	if len(field[len(field)-1]) == 0 {
		field = field[:len(field)-1]
	}
	bfield := make(Field, 0, len(field))
	for _, s := range field {
		bfield = append(bfield, []byte(s))
	}
	return bfield
}

func (field *Field) Expand() {
	top := strings.Contains(string((*field)[0]), "#")
	bottom := strings.Contains(string((*field)[len(*field)-1]), "#")
	var left, right bool
	for _, s := range *field {
		if s[0] == '#' {
			left = true
		}
		if s[len(s)-1] == '#' {
			right = true
		}
		if left && right {
			break
		}
	}
	if left || right {
		for i, bs := range *field {
			string_s := string(bs)
			if left {
				string_s = "." + string_s
			}
			if right {
				string_s += "."
			}
			(*field)[i] = []byte(string_s)
		}
	}
	if top || bottom {
		empty := []byte(strings.Repeat(".", len((*field)[0])))
		*field = append(*field, empty)
		if top {
			copy((*field)[1:], *field)
			(*field)[0] = empty
			if bottom {
				empty2 := []byte(strings.Repeat(".", len((*field)[0])))
				*field = append(*field, empty2)
			}
		}
	}
}

type Dir struct{ dx, dy int }

type Pos struct{ x, y int }

func (field Field) ProposedMove(p Pos, firstdir int) Pos {
	directions := [][]Dir{
		{{0, -1}, {-1, -1}, {1, -1}},
		{{0, 1}, {-1, 1}, {1, 1}},
		{{-1, 0}, {-1, 1}, {-1, -1}},
		{{1, 0}, {1, -1}, {1, 1}},
	}
	// first, look around everywhere to see if we want to move
	need_move := false
LookAround:
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if (dx != 0 || dy != 0) && field[p.y+dy][p.x+dx] == '#' {
				need_move = true
				break LookAround
			}
		}
	}
	if !need_move {
		return p
	}
	for dir := 0; dir < len(directions); dir++ {
		look := directions[(dir+firstdir)%len(directions)]
		other_elf := false
		for _, l := range look {
			if field[p.y+l.dy][p.x+l.dx] == '#' {
				other_elf = true
				break
			}
		}
		if !other_elf {
			// there is no other elf in this direction, propose moving there
			return Pos{p.x + look[0].dx, p.y + look[0].dy}
		}
	}
	// there are elfes around us everywhere! propose not to move
	return p
}

func (field Field) Evolve(firstdir int) {
	// "proposed" contains the proposed moves.
	// if true, move is possible, if false, more than 1 elf proposed to move there
	proposed := make(map[Pos]bool)
	elves := make(map[Pos]Pos)
	for y, line := range field {
		for x, item := range line {
			if item == '#' {
				elf := Pos{x, y}
				moveto := field.ProposedMove(elf, firstdir)
				// only remember the elf if it wants to move
				if moveto != elf {
					if _, alreadyproposed := proposed[moveto]; !alreadyproposed {
						elves[elf] = moveto
						proposed[moveto] = true
						//fmt.Printf("Propose moving elf at %v to %v\n", elf, moveto)
					} else {
						//fmt.Printf("Elf at %v proposed moving to %v but another elf goes there too\n", elf, moveto)
						proposed[moveto] = false
					}
				} else {
					//fmt.Printf("Elf at %v does not move\n", elf)
				}
			}
		}
	}
	// part 2, actually move the elves, if they are the only one that proposed this move
	for elf, moveto := range elves {
		if proposed[moveto] {
			if field[elf.y][elf.x] != '#' {
				panic(fmt.Sprintf("Expected an elf at %v, found [%c]", elf, field[elf.y][elf.x]))
			}
			if field[moveto.y][moveto.x] != '.' {
				panic(fmt.Sprintf("Moving elf from %v to %v, but that is not free, got: [%c]", elf, moveto, field[moveto.y][moveto.x]))
			}
			field[elf.y][elf.x] = '.'
			field[moveto.y][moveto.x] = '#'
		}
	}
}

func (f Field) String() string {
	ret := ""
	for _, s := range f {
		ret += string(s) + "\n"
	}
	return ret
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	field := parse_input(os.Args[1])
	parsetime := time.Now()
	for round := 1; round <= 10; round++ {
		field.Expand()
		fmt.Printf("\nStarting round %d:\n%s", round, field)
		field.Evolve(round - 1)
	}
	part1time := time.Now()
	fmt.Printf("\nEnd result:\n%s", field)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", part1time.Sub(parsetime))
}
