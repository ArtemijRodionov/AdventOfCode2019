package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Color int

const (
	Black Color = iota
	White
)

type Direction int

const (
	Left Direction = iota
	Right
	Top
	Down
)

func (d Direction) Turn(t Direction) Direction {
	transitions := map[Direction]map[Direction]Direction{
		Left:  {Left: Down, Right: Top, Top: Left, Down: Right},
		Right: {Left: Top, Right: Down, Top: Right, Down: Left},
	}

	return transitions[t][d]
}

type Position struct{ X, Y int }

func (p Position) Step(dir Direction) Position {
	transitions := map[Direction]Position{
		Left:  {-1, 0},
		Right: {1, 0},
		Top:   {0, 1},
		Down:  {0, -1},
	}
	to := transitions[dir]
	return Position{p.X + to.X, p.Y + to.Y}
}

type Surface map[Position]Color

func (s Surface) Color(p Position) Color    { return s[p] }
func (s Surface) Paint(p Position, c Color) { s[p] = c }
func NewSurface() Surface {
	return Surface(make(map[Position]Color, 0))
}
func (s Surface) Print() {
	var minx, miny, maxx, maxy int

	for p := range s {
		if minx > p.X {
			minx = p.X
		}
		if maxx < p.X {
			maxx = p.X
		}
		if miny > p.Y {
			miny = p.Y
		}
		if maxy < p.Y {
			maxy = p.Y
		}
	}

	w := bufio.NewWriter(os.Stdout)
	// top-down: max -> min
	for i := maxy; i >= miny; i-- {
		// left to right: min -> max
		for j := minx; j <= maxx; j++ {
			if v, ok := s[Position{j, i}]; !ok {
				w.WriteByte(' ')
			} else {
				w.WriteString(fmt.Sprintf("%d", v))
			}
		}
		w.WriteByte('\n')
	}
	w.Flush()
}

type Robot struct {
	surface Surface
	pos     Position
	dir     Direction
}

func NewRobot() Robot {
	return Robot{surface: NewSurface(), dir: Top}
}

func (r Robot) Color() Color {
	return r.surface.Color(r.pos)
}

func (r *Robot) Paint(c Color) {
	r.surface.Paint(r.pos, c)
}

func (r *Robot) Step(dir Direction) {
	r.dir = r.dir.Turn(dir)
	r.pos = r.pos.Step(r.dir)
}

func Run(program []int, initColor Color) Surface {
	vm := NewIntCode(program)
	robot := NewRobot()
	vm.Run()

	run := true
	vm.Input <- int(initColor)
	for run {
		robot.Paint(Color(<-vm.Output))
		robot.Step(Direction(<-vm.Output))

		select {
		case vm.Input <- int(robot.Color()):
		case <-vm.Halt:
			run = false
		}
	}

	return robot.surface
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an input file as the first argument")
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Can't read file: ", err)
	}

	tokens := strings.Split(strings.TrimSpace(string(input)), ",")
	program := make([]int, len(tokens))
	for i, t := range tokens {
		if program[i], err = strconv.Atoi(t); err != nil {
			log.Fatalf("Can't parse number %d. Err: %s", t, err)
		}
	}

	fmt.Println("Part 1: ", len(Run(program, Black)))
	Run(program, White).Print()
}
