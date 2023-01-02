package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type NumEntry struct {
	val, pos int
}

type EntryList []NumEntry

// the entries in the list, with their updated positions, and the possible positions, sorted.
type NumList struct {
	entries       EntryList
	positions     []int
	sortedentries EntryList
	zeropos       int
}

const spacing = 1 << 32

const minspacing = 1 << 8

func parse_input(filename string) NumList {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	entries := make([]NumEntry, 0)
	positions := make([]int, 0)
	pos := spacing / 2
	for _, line := range strings.Split(string(inbuf[:]), "\n") {
		if len(line) == 0 {
			continue
		}
		if val, err := strconv.Atoi(line); err == nil {
			entries = append(entries, NumEntry{val: val, pos: pos})
			positions = append(positions, pos)
			pos += spacing
		} else {
			panic(fmt.Sprintf("Cannot convert `%s' to numeric: %s", line, err))
		}
	}
	return NumList{entries: entries, positions: positions}
}

func (nl NumList) Len() int { return len(nl.entries) }

func (el EntryList) Len() int { return len(el) }

func (el EntryList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

func (el EntryList) Less(i, j int) bool { return el[i].pos < el[j].pos }

// Populate sortedentries, and set zeropos
func (nl *NumList) MakeSorted() {
	if nl.sortedentries != nil {
		return
	}
	nl.sortedentries = make(EntryList, len(nl.entries))
	copy(nl.sortedentries, nl.entries)
	sort.Sort(nl.sortedentries)
	for i, entry := range nl.sortedentries {
		if entry.val == 0 {
			nl.zeropos = i
			return
		}
	}
	panic("No zero entry in list")
}

func (nl NumList) Str() string {
	nl.MakeSorted()
	result := ""
	for _, entry := range nl.sortedentries {
		if len(result) > 0 {
			result += " "
		}
		result += fmt.Sprintf("%d", entry.val)
	}
	return result
}

func (nl NumList) Move(i int) {
	val := nl.entries[i].val
	curpos := nl.entries[i].pos
	realpos, ok := sort.Find(len(nl.positions), func(i int) int {
		if curpos < nl.positions[i] {
			return -1
		} else if curpos == nl.positions[i] {
			return 0
		} else {
			return 1
		}
	})
	if !ok {
		panic(fmt.Sprintf("Internal error, positions not updated, cannot find %d", curpos))
	}
	newrealpos := (realpos + val) % (len(nl.entries) - 1)
	if newrealpos < 0 {
		newrealpos += len(nl.entries) - 1
	}
	if newrealpos == realpos {
		// nothing to do
		return
	} else if newrealpos < realpos {
		// moving down, we really need to insert before newrealpos. So decrement newrealpos so we can still insert after it.
		newrealpos--
	}
	// invalidate the sorted entries, if any
	nl.sortedentries = nil
	// we need to insert current value between newrealpos and newrealpos+1
	// or, if newrealpos < 0, just take half the value at position 0
	var diff, prev_val int
	if newrealpos >= 0 {
		prev_val = nl.positions[newrealpos]
		diff = nl.positions[newrealpos+1] - prev_val
	} else {
		diff = nl.positions[0]
		prev_val = 0
	}
	if diff <= 1 {
		panic(fmt.Sprintf("Cannot insert a number between real positions %d and %d, rel positions %d and %d", newrealpos, newrealpos+1, prev_val, nl.positions[newrealpos+1]))
	}
	newpos := prev_val + diff/2
	nl.entries[i].pos = newpos
	// fix the positions array. "realpos" goes away, and add an entry after "newrealpos".
	if realpos < newrealpos {
		copy(nl.positions[realpos:newrealpos], nl.positions[realpos+1:newrealpos+1])
		nl.positions[newrealpos] = newpos
		//fmt.Printf("Moving up val=%d from realpos=%d to newrealpos=%d. Original rel position=%d, new rel position=%d\n", val, realpos, newrealpos, curpos, newpos)
	} else { // realpos > newrealpos
		copy(nl.positions[newrealpos+2:realpos+1], nl.positions[newrealpos+1:realpos])
		nl.positions[newrealpos+1] = newpos
		//fmt.Printf("Moving down val=%d from realpos=%d to newrealpos=%d. Original rel position=%d, new rel position=%d\n", val, realpos, newrealpos, curpos, newpos)
	}
	// fix the positions array. "realpos" goes away
	//fmt.Printf("numlist=%v\n", nl)
}

func (nl NumList) Offset0(i int) int {
	nl.MakeSorted()
	wantedpos := (i + nl.zeropos) % len(nl.sortedentries)
	return nl.sortedentries[wantedpos].val
}

func (nl *NumList) Rebalance() {
	// check if rebalance is necessary
	is_needed := false
	prev := nl.positions[0]
	for _, pos := range nl.positions[1:] {
		if pos-prev < minspacing {
			is_needed = true
			break
		}
		prev = pos
	}
	if !is_needed {
		return
	}
	//fmt.Printf("Rebalancing...\n")
	remap := make(map[int]int, len(nl.positions))
	newpos := spacing / 2
	new_positions := make([]int, 0, len(nl.positions))
	for _, pos := range nl.positions {
		remap[pos] = newpos
		new_positions = append(new_positions, newpos)
		newpos += spacing
	}
	nl.positions = new_positions
	for i := range nl.entries {
		nl.entries[i].pos = remap[nl.entries[i].pos]
	}
}

func (nl *NumList) Reset() {
	nl.sortedentries = nil
	pos := spacing / 2
	nl.positions = make([]int, 0, len(nl.entries))
	for i := range nl.entries {
		nl.entries[i].pos = pos
		nl.positions = append(nl.positions, pos)
		pos += spacing
	}
}

func (nl *NumList) Decrypt(key int) {
	nl.sortedentries = nil
	for i := range nl.entries {
		nl.entries[i].val *= key
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	numlist := parse_input(os.Args[1])
	parsetime := time.Now()
	for i := 0; i < numlist.Len(); i++ {
		numlist.Move(i)
	}
	part1 := 0
	for _, offset := range []int{1000, 2000, 3000} {
		num := numlist.Offset0(offset)
		fmt.Printf("Number at offset %d is %d\n", offset, num)
		part1 += num
	}
	part1time := time.Now()
	fmt.Println("part1 sum: ", part1)
	numlist.Reset()
	numlist.Decrypt(811589153)
	for round := 1; round <= 10; round++ {
		for i := 0; i < numlist.Len(); i++ {
			numlist.Move(i)
		}
		// Make sure there is enough room between the positions to insert some numbers
		numlist.Rebalance()
		//fmt.Printf("After round %d, list = %s\n", round, numlist.Str())
	}
	part2 := 0
	for _, offset := range []int{1000, 2000, 3000} {
		num := numlist.Offset0(offset)
		fmt.Printf("Number at offset %d is %d\n", offset, num)
		part2 += num
	}
	part2time := time.Now()
	fmt.Printf("part 2 sum: %d\n", part2)
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("Part 1 took: %s\n", part1time.Sub(parsetime))
	fmt.Printf("Part 2 took: %s\n", part2time.Sub(part1time))
}
