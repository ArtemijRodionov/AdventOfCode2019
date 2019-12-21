package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func Angle(l, r Asteroid) float64 {
	return math.Atan2(l.X-r.X, l.Y-r.Y)
}

func Hypot(l, r Asteroid) float64 {
	return math.Hypot(l.X-r.X, l.Y-r.Y)
}

func PrintSpace(rows []string) {
	w := bufio.NewWriter(os.Stdout)
	for y, row := range rows {
		for x, cell := range row {
			if cell == '#' {
				w.WriteString(fmt.Sprintf("%02d %02d|", x, y))
			} else {
				w.WriteString("## ##|")
			}
		}
		w.WriteByte('\n')
	}
	w.Flush()
}

type Asteroid struct{ X, Y float64 }
type Asteroids []Asteroid

func (xs Asteroids) Copy() Asteroids {
	cxs := make(Asteroids, len(xs))
	copy(cxs, xs)
	return cxs
}

func (xs *Asteroids) Push(x Asteroid) {
	*xs = append(*xs, x)
}

func (xs *Asteroids) Pop() Asteroid {
	i := len(*xs) - 1
	x := (*xs)[i]
	*xs = (*xs)[0:i]
	return x
}

type By func(a, b Asteroid) bool

func (b By) Sort(as Asteroids) {
	s := &AsteroidSorter{as, b}
	sort.Sort(s)
}

func ByAngle(base Asteroid) By {
	return func(l, r Asteroid) bool {
		return Angle(l, base) > Angle(r, base)
	}
}

func ByHypot(base Asteroid) By {
	return func(l, r Asteroid) bool {
		return Hypot(base, l) < Hypot(base, r)
	}
}

type AsteroidSorter struct {
	as Asteroids
	by By
}

func (s *AsteroidSorter) Len() int           { return len(s.as) }
func (s *AsteroidSorter) Swap(i, j int)      { s.as[i], s.as[j] = s.as[j], s.as[i] }
func (s *AsteroidSorter) Less(i, j int) bool { return s.by(s.as[i], s.as[j]) }

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an input file as the first argument.")
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Can't read input file: ", err)
	}

	rows := strings.Split(strings.TrimSpace(string(input)), "\n")
	asteroids := make(Asteroids, 0)
	for y, row := range rows {
		for x, cell := range row {
			if cell == '#' {
				asteroids.Push(Asteroid{float64(x), float64(y)})
			}
		}
	}

	max := 0
	best := Asteroid{}
	for i, a := range asteroids {
		visible := make(map[float64]bool)
		for j, b := range asteroids {
			if i == j {
				continue
			}
			visible[Angle(b, a)] = true
		}
		visibleCount := len(visible)
		if max < visibleCount {
			max = visibleCount
			best = a
		}
	}
	fmt.Println("Part 1:", best, max)

	sortedAsteroids := asteroids.Copy()
	ByAngle(best).Sort(sortedAsteroids)

	// unique sorted angles
	angles := make([]float64, 0)
	// asteroids at the same angle sorted by distance from station
	asteroidsLine := make(map[float64]*Asteroids)
	for i, a := range sortedAsteroids {
		if a == best {
			continue
		}
		angle := Angle(a, best)
		if _, ok := asteroidsLine[angle]; !ok {
			init := sortedAsteroids[i : i+1].Copy()
			asteroidsLine[angle] = &init
			angles = append(angles, angle)
		} else {
			asteroidsLine[angle].Push(a)
		}
	}

	sorter := ByHypot(best)
	for k := range asteroidsLine {
		sorter.Sort(*asteroidsLine[k])
	}

	vaporized := make(Asteroids, 0)
	for len(asteroidsLine) != 0 {
		for _, angle := range angles {
			if line, ok := asteroidsLine[angle]; !ok {
				continue
			} else {
				if len(*line) == 0 {
					delete(asteroidsLine, angle)
					continue
				}

				vaporized = append(vaporized, line.Pop())
			}
		}
	}

	result := vaporized[199]
	// Visualize for debug
	//PrintSpace(rows)
	//fmt.Println(vaporized[:200], len(vaporized), "\n")

	fmt.Println("Part 2: ", result, result.X*100+result.Y)
}
