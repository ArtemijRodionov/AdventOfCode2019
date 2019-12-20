package main

import (
    "os"
    "io/ioutil"
    "log"
    "fmt"
    "strings"
    "math"
)

type Asteroid struct { X, Y float64 }

func (l Asteroid) Angle(r Asteroid) float64 {
    return math.Mod(math.Atan2(r.X - l.X, l.Y - r.Y), 2 * math.Pi)
}

func main() {
    if len(os.Args) != 2 { log.Fatal("Provide an input file as the first argument.") }
    input, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal("Can't read input file: ", err)
    }

    rows := strings.Split(strings.TrimSpace(string(input)), "\n")
    asteroids := make([]Asteroid, 0)
    for y, row := range rows {
        for x, cell := range row {
            if cell == '#' {
                asteroids = append(asteroids, Asteroid{float64(x), float64(y)})
            }
        }
    }


    max := 0
    asteroid := Asteroid{}
    for i, a := range asteroids {
        visible := make(map[float64]bool)
        for j, b := range asteroids {
            if i == j { continue }
            visible[a.Angle(b)] = true
        }
        visibleCount := len(visible)
        if max < visibleCount {
            max = visibleCount
            asteroid = a
        }
    }
    fmt.Println("Part 1:", asteroid, max)
}

