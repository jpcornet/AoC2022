package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Field []string

type PathElem struct {
	dist   int
	rotate int
}

type Path []PathElem

type Pos [2]int

var directions [4]Pos

type PosDir struct {
	pos Pos
	dir int
}

func init_vars() {
	directions = [4]Pos{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}
}

func parse_input(filename string) (Field, Path) {
	inbuf, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fieldstr, pathstr, ok := strings.Cut(string(inbuf[:]), "\n\n")
	if !ok {
		panic("Invalid input")
	}
	field := strings.Split(fieldstr, "\n")
	if len(field[len(field)-1]) == 0 {
		field = field[:len(field)-1]
	}
	path := make(Path, 0, 10)
	pathstr = strings.TrimRight(pathstr, "\n")
	for len(pathstr) > 0 {
		var pe PathElem
		lroffset := strings.IndexAny(pathstr, "LR")
		if lroffset == -1 {
			lroffset = len(pathstr)
		} else if pathstr[lroffset] == 'L' {
			pe.rotate = -1
		} else if pathstr[lroffset] == 'R' {
			pe.rotate = 1
		}
		if lroffset != 0 {
			if i, err := strconv.Atoi(pathstr[:lroffset]); err == nil {
				pe.dist = i
			} else {
				panic(fmt.Sprintf("Invalid input in path, not a number: %s (%v)", pathstr[:lroffset], err))
			}
		}
		path = append(path, pe)
		if lroffset == len(pathstr) {
			break
		} else {
			pathstr = pathstr[lroffset+1:]
		}
	}
	return field, path
}

func get_startpos(field Field) PosDir {
	x := strings.Index(field[0], ".")
	if x == -1 {
		panic("Invalid input")
	}
	return PosDir{
		pos: Pos{x, 0},
		dir: 0,
	}
}

type WrapFunc func(PosDir) PosDir

func walk_path(field Field, path Path, posd PosDir, wrapper WrapFunc) PosDir {
	for _, p := range path {
		for i := 0; i < p.dist; i++ {
			newpos := Pos{posd.pos[0] + directions[posd.dir][0], posd.pos[1] + directions[posd.dir][1]}
			newposdir := PosDir{pos: newpos, dir: posd.dir}
			if newpos[0] < 0 || newpos[1] < 0 || newpos[1] >= len(field) || newpos[0] >= len(field[newpos[1]]) || field[newpos[1]][newpos[0]] == ' ' {
				//fmt.Printf("At %d,%d direction %c, wrapping\n", posd.pos[0], posd.pos[1], dirtochar(posd.dir))
				newposdir = wrapper(posd)
				newpos = newposdir.pos
				//fmt.Printf("... wrapped to %d,%d direction %c\n", newpos[0], newpos[1], dirtochar(newposdir.dir))
			}
			chr := field[newpos[1]][newpos[0]]
			if chr == '.' {
				posd = newposdir
			} else if chr == '#' {
				break
			} else {
				panic(fmt.Sprintf("Invalid character [%c] in field at %d,%d", chr, newpos[0], newpos[1]))
			}
		}
		posd.dir = (posd.dir + p.rotate + 4) % 4
		//fmt.Printf("At %d,%d rotated to %c\n", posd.pos[0], posd.pos[1], dirtochar(posd.dir))
	}
	return posd
}

func make_basic_wrapper(field Field) WrapFunc {
	return func(pd PosDir) PosDir {
		newpos := pd.pos
		for newpos[0] >= 0 && newpos[1] >= 0 && newpos[1] < len(field) && newpos[0] < len(field[newpos[1]]) && (field[newpos[1]][newpos[0]] == '.' || field[newpos[1]][newpos[0]] == '#') {
			pd.pos = newpos
			newpos[0] -= directions[pd.dir][0]
			newpos[1] -= directions[pd.dir][1]
		}
		return pd
	}
}

type OtherFace struct {
	facenr int
	dir    int
}

type FaceInfo struct {
	pos      Pos
	adjacent [4]OtherFace
}

type CubeLayout struct {
	dim     int
	face    [6]FaceInfo
	facepos map[Pos]int
}

type FoldConnect struct {
	normal, perp, rotate int
	free                 []Pos
}

func dirtochar(dir int) rune {
	dirs := "→↓←↑"
	for _, arrow := range dirs {
		if dir == 0 {
			return arrow
		}
		dir--
	}
	return '?'
}

func analyze_cube(field Field) CubeLayout {
	var layout CubeLayout
	// first, get the cube dimension. Scan the total surface and the max line length
	surface := 0
	maxline := 0
	for _, l := range field {
		left := strings.IndexAny(l, ".#")
		right := strings.LastIndexAny(l, ".#")
		if left == -1 || right == -1 {
			panic("Invalid input")
		}
		linelen := right - left + 1
		surface += linelen
		if linelen > maxline {
			maxline = linelen
		}
	}
	// usually, maxline will be the dimension, but maybe we need to divide by 2 or 3
	for div := 1; div <= 3; div++ {
		if maxline%div != 0 {
			continue
		}
		testsurf := 6 * maxline * maxline / div / div
		if surface == testsurf {
			layout.dim = maxline / div
			break
		} else if testsurf < surface {
			panic(fmt.Sprintf("Cannot determine cube dimensions. maxline=%d surface=%d\n", maxline, surface))
		}
	}
	if layout.dim == 0 {
		panic(fmt.Sprintf("Cannot determine cube dimensions, no divider found. maxline=%d surface=%d\n", maxline, surface))
	}
	// now get the positions of all the cube faces
	facenr := 0
	layout.facepos = make(map[Pos]int, 6)
	for y := 0; y < len(field); y += layout.dim {
		left := strings.IndexAny(field[y], ".#")
		right := strings.LastIndexAny(field[y], ".#")
		for x := left; x < right; x += layout.dim {
			// sanity check
			if facenr >= 6 {
				panic("Logic error, too many cube faces")
			}
			layout.face[facenr].pos = Pos{x, y}
			layout.facepos[Pos{x, y}] = facenr
			//fmt.Printf("Found face#%d at %d,%d\n", facenr, x, y)
			facenr++
		}
	}
	// prepare a data structure that allows us to find the "connecting" face, and the direction
	foldconnects := []FoldConnect{
		{normal: 1, perp: 0, rotate: 0},                           // Just directly connected
		{normal: 1, perp: 1, rotate: 1},                           // in a 90 degree angle to the right
		{normal: 1, perp: -1, rotate: -1},                         // in a 90 degree angle to the left
		{normal: 1, perp: 2, rotate: 2},                           // 2 removed, doing 180 degree rotate
		{normal: 1, perp: -2, rotate: 2},                          // 2 removed the other way
		{normal: -1, perp: 2, rotate: 2, free: []Pos{{0, 1}}},     // 2 removed going "backwards"
		{normal: -1, perp: -2, rotate: 2, free: []Pos{{0, -1}}},   // going backwards the other way
		{normal: -3, perp: 0, rotate: 2},                          // loop around 3 to the other side
		{normal: 1, perp: 3, rotate: 3},                           // 3 over
		{normal: 1, perp: -3, rotate: -3},                         // 3 over other side
		{normal: -1, perp: 3, rotate: -1, free: []Pos{{0, 2}}},    // backwards and 3 to the side
		{normal: -1, perp: -3, rotate: 1, free: []Pos{{0, -2}}},   // other side
		{normal: -3, perp: 1, rotate: 3, free: []Pos{{-1, 1}}},    // 3 backwards, and to the right
		{normal: -3, perp: -1, rotate: -3, free: []Pos{{-1, -1}}}, // 3 backwards and to the left
		{normal: -1, perp: 4, rotate: 0},                          // 4 over, no rotation
		{normal: -1, perp: -4, rotate: 0},
		{normal: -2, perp: 2, rotate: -1, free: []Pos{{0, 1}, {-1, 2}}},
		{normal: -2, perp: -2, rotate: 1, free: []Pos{{0, -1}, {-1, -2}}},
		{normal: -2, perp: 3, rotate: 0},
		{normal: -2, perp: -3, rotate: 0},
		{normal: -3, perp: 2, rotate: 0},
		{normal: -3, perp: -2, rotate: 0},
	}
	// Now, for each face where the adjacent one isn't next to it, find the real adjacent face and the direction
	for fi, f := range layout.face {
		for di, d := range directions {
			found := false
			for _, fold := range foldconnects {
				// note: to get consistent turn directions, we calculate the perpendicular movement als x += -dy, y += dx. Or multiply movement by matrix (0 -1) // (1 0)
				dpos := Pos{f.pos[0] + layout.dim*d[0]*fold.normal - layout.dim*d[1]*fold.perp, f.pos[1] + layout.dim*d[1]*fold.normal + layout.dim*d[0]*fold.perp}
				if otherfi, ok := layout.facepos[dpos]; ok {
					// there is a face where we expected it
					found = true
					// make sure every block in "absent" is actually absent, or we didn't really find it
					for _, absent := range fold.free {
						apos := Pos{f.pos[0] + layout.dim*d[0]*absent[0] - layout.dim*d[1]*absent[1], f.pos[1] + layout.dim*d[1]*absent[0] + layout.dim*d[0]*absent[1]}
						if _, notok := layout.facepos[apos]; notok {
							//fmt.Printf("Would go from face#%d at %d,%d going %c to face #%d, but there is face#%d at %d,%d that blocks it\n", fi, f.pos[0], f.pos[1], dirtochar(di), otherfi, blk, apos[0], apos[1])
							found = false
							break
						}
					}
					if found {
						layout.face[fi].adjacent[di].facenr = otherfi
						layout.face[fi].adjacent[di].dir = (di + fold.rotate + 4) % 4
						//fmt.Printf("From face#%d at %d,%d going %c we find face#%d at %d,%d going %c\n", fi, f.pos[0], f.pos[1], dirtochar(di), otherfi, layout.face[otherfi].pos[0], layout.face[otherfi].pos[1], dirtochar(layout.face[fi].adjacent[di].dir))
						break
					}
				}
			}
			if !found {
				panic(fmt.Sprintf("Cannot find connecting face from face#%d going in direction %c", fi, dirtochar(di)))
			}
		}
	}
	return verify_layout(layout)
}

// verify that every cube face is connected both ways with the other cube face, and vice versa.
func verify_layout(l CubeLayout) CubeLayout {
	for facenr, face := range l.face {
		for dir, finfo := range face.adjacent {
			reversedir := (dir + 2) % 4
			reverse_on_face := (finfo.dir + 2) % 4
			pointsback := l.face[finfo.facenr].adjacent[reverse_on_face]
			if pointsback.facenr != facenr || pointsback.dir != reversedir {
				panic(fmt.Sprintf("Hm, on face#%d, direction %c, we point to face#%d, direction %c, but that points back to face#%d, direction %c", facenr, dirtochar(dir), finfo.facenr, dirtochar(finfo.dir), pointsback.facenr, dirtochar(pointsback.dir)))
			}
		}
	}
	return l
}

func make_cube_wrapper(layout CubeLayout) WrapFunc {
	return func(pd PosDir) PosDir {
		// determine offset within the cube face that we're in
		x := pd.pos[0] % layout.dim
		y := pd.pos[1] % layout.dim
		// determine the cube face we're in, based on the top left corner
		facepos := Pos{pd.pos[0] - x, pd.pos[1] - y}
		facenr := layout.facepos[facepos]
		// determine the cube face that we're going to
		other := layout.face[facenr].adjacent[pd.dir]
		// we go from direction pd.dir to other.dir. Rotate x and y accordingly.
		var ox, oy int
		switch (other.dir - pd.dir + 4) % 4 {
		case 0:
			// no need to rotate
			ox, oy = x, y
		case 1:
			// rotate clockwise 90
			ox, oy = layout.dim-1-y, x
		case 2:
			// rotate 180
			ox, oy = layout.dim-1-x, layout.dim-1-y
		case 3:
			// rotate counterclockwise 90
			ox, oy = y, layout.dim-1-x
		}
		// ox, oy is now the corresponding point where we left off. Map this to the other side
		ox -= (layout.dim - 1) * directions[other.dir][0]
		oy -= (layout.dim - 1) * directions[other.dir][1]
		newfacepos := layout.face[other.facenr].pos
		mapped := PosDir{pos: Pos{newfacepos[0] + ox, newfacepos[1] + oy}, dir: other.dir}
		//fmt.Printf("Cubewrap from %d,%d facing %c on face#%d (relative pos %d,%d) to %d,%d facing %c on face#%d (relative pos %d,%d)\n",
		//	pd.pos[0], pd.pos[1], dirtochar(pd.dir), facenr, x, y, mapped.pos[0], mapped.pos[1], dirtochar(other.dir), other.facenr, ox, oy)
		return mapped
	}
}

func to_pass(posd PosDir) int {
	return 1000*(posd.pos[1]+1) + 4*(posd.pos[0]+1) + posd.dir
}

func main() {
	if len(os.Args) != 2 {
		panic("Provide input file")
	}
	starttime := time.Now()
	field, path := parse_input(os.Args[1])
	init_vars()
	parsetime := time.Now()
	startpos := get_startpos(field)
	endpos := walk_path(field, path, startpos, make_basic_wrapper(field))
	walktime := time.Now()
	fmt.Printf("endpos part 1: %v. Password: %d\n", endpos, to_pass(endpos))
	cube_layout := analyze_cube(field)
	endpos2 := walk_path(field, path, startpos, make_cube_wrapper(cube_layout))
	walk2time := time.Now()
	fmt.Printf("endpos part 2: %v, Password: %d\n", endpos2, to_pass(endpos2))
	fmt.Printf("Parse took: %s\n", parsetime.Sub(starttime))
	fmt.Printf("part 1 took: %s\n", walktime.Sub(parsetime))
	fmt.Printf("part 2 took: %s\n", walk2time.Sub(walktime))
}
