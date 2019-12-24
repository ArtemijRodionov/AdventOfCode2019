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

type Vector [3]int

func (v Vector) Norm() int {
    res := 0
    for i := range v {
        res += int(math.Abs(float64(v[i])))
    }
    return res
}

func Sum(l, r Vector) Vector {
    v := Vector{}
    for i := range v {
        v[i] = l[i] + r[i]
    }
    return v
}

type Satelite struct {
    position Vector
    velocity Vector
}

func NewSatelite(pos Vector) Satelite {
    return Satelite{ pos, Vector{} }
}

type Space struct {
    satelites []Satelite
}

func (s Space) Energy() int {
    total := 0
    for _, m := range s.satelites {
        total += m.position.Norm() * m.velocity.Norm()
    }
    return total
}

func (s Space) Copy() Space {
    m := make([]Satelite, len(s.satelites))
    for i, x := range s.satelites {
        m[i] = x
    }
    return Space{m}
}

func (s Space) Step() {
    gravities := make([]Vector, len(s.satelites))
    for i := range s.satelites {
        gravities[i] = Vector{}
        l := s.satelites[i].position

        for j := range s.satelites {
            if i == j { continue }

            r := s.satelites[j].position
            gravity := Vector{}
            for ip := range r {
                if l[ip] == r[ip] {
                    continue
                }

                if l[ip] < r[ip] {
                    gravity[ip]++
                } else {
                    gravity[ip]--
                }
            }

            gravities[i] = Sum(gravities[i], gravity)
        }
    }

    for i := range s.satelites {
        s.satelites[i].velocity = Sum(s.satelites[i].velocity, gravities[i])
        s.satelites[i].position = Sum(s.satelites[i].position, s.satelites[i].velocity)
    }
}

func Gcd(a, b float64) float64 {
    mod := math.Mod(a, b)
    for mod != 0.0 {
        a, b = b, mod
        mod = math.Mod(a, b)
    }
    return b
}

func Lcm(a, b float64) float64 {
    return math.Abs(a * b) / Gcd(a, b)
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
    satelites := make([]Satelite, len(rows))
    for i, r := range rows {
        pos := Vector{}
        for j, t := range strings.Split(r, ",") {
            if pos[j], err = strconv.Atoi(t[2:]); err != nil {
                log.Fatal("Can't parse input: ", err)
            }
        }
        satelites[i] = NewSatelite(pos)
    }

    ref := Space{satelites}
    s := ref.Copy()
    for i := 0; i < 1000; i++ {
        s.Step()
    }
    fmt.Println("Part 1:", s.Energy())

    s = ref.Copy()
    // Each satelite position affects on velocity of other satelites on the same axis.
    // So we need to find revolving period for each axis of all satelites and find LCM for these periods.
    // https://en.wikipedia.org/wiki/Least_common_multiple#Planetary_alignment
    steps := [3]int{}
    orbits := [3]map[[4][2]int]bool{}
    for i := range orbits {
        orbits[i] = make(map[[4][2]int]bool, 0)
    }

    for n := 0; true; n++ {
        for i := range steps {
            if steps[i] != 0 { continue }

            posVel := [4][2]int{}
            for j := range posVel {
                posVel[j][0] = s.satelites[j].position[i]
                posVel[j][1] = s.satelites[j].velocity[i]
            }

            if _, ok := orbits[i][posVel]; ok {
                steps[i] = n
            } else {
                orbits[i][posVel] = true
            }
        }

        done := true
        for _, s := range steps {
            if s == 0 {
                done = false
            }
        }

        if done { break }

        s.Step()
    }

    res := float64(steps[0])
    for _, s := range steps[1:] {
        res = Lcm(res, float64(s))
    }

    fmt.Println("Part 2:", int(res))
}

