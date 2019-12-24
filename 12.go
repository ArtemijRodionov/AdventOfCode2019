package main

import (
    "os"
    "io/ioutil"
    "fmt"
    "log"
    "strings"
    "strconv"
    "math"
)

type Axis int
const (
    X Axis = iota
    Y
    Z
)

var All = [3]Axis{X, Y, Z}

type Vector map[Axis]int
func NewVector(x, y, z int) *Vector {
    return &Vector{ X: x, Y: y, Z: z }
}

func (v Vector) Norm() int {
    res := 0
    for _, a := range All {
        res += int(math.Abs(float64(v[a])))
    }
    return res
}

func (v *Vector) Add(r Vector) {
    for _, a := range All {
        (*v)[a] += r[a]
    }
}


type Moon struct {
    position *Vector
    velocity *Vector
}

func NewMoon(pos *Vector) Moon {
    return Moon{ pos, NewVector(0, 0, 0) }
}

type Space struct {
    moons []Moon
}

func (s Space) Energy() int {
    total := 0
    for _, m := range s.moons {
        total += m.position.Norm() * m.velocity.Norm()
    }
    return total
}

func (s Space) Step() {
    gravities := make([]*Vector, len(s.moons))
    for i := range s.moons {
        gravities[i] = NewVector(0, 0, 0)
        l := *s.moons[i].position

        for j := range s.moons {
            if i == j { continue }

            r := *s.moons[j].position
            gravity := NewVector(0, 0, 0)
            for _, a := range All {
                if l[a] == r[a] {
                    continue
                }

                if l[a] < r[a] {
                    (*gravity)[a]++
                } else {
                    (*gravity)[a]--
                }
            }

            gravities[i].Add(*gravity)
        }
    }

    for i := range s.moons {
        s.moons[i].velocity.Add(*gravities[i])
        s.moons[i].position.Add(*s.moons[i].velocity)
    }
}


func main() {
    if len(os.Args) != 2 { log.Fatal("Provide an input file as the first argument") }
    raw, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal("Can't read file: ", err)
    }

    input := string(raw)
    for _, s := range []string{"<", ">", " "} {
        input = strings.Replace(input, s, "", -1)
    }

    rows := strings.Split(strings.TrimSpace(string(input)), "\n")
    moons := make([]Moon, len(rows))
    for i, r := range rows {
        pos := NewVector(0, 0, 0)
        for j, t := range strings.Split(r, ",") {
            if (*pos)[All[j]], err = strconv.Atoi(t[2:]); err != nil {
                log.Fatal("Can't parse input: ", err)
            }
        }
        moons[i] = NewMoon(pos)
    }

    s := Space{moons}
    for i := 0; i < 1000; i++ {
        s.Step()
    }

    fmt.Println("Part 1:", s.Energy())
}

