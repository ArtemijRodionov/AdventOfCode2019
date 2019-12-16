package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Traverse(orbits map[string]string, orbit string, do func(string)) {
	ok := true
	for ok {
		do(orbit)
		orbit, ok = orbits[orbit]
	}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	lineSeparated := strings.Split(strings.TrimSpace(string(input)), "\n")

	orbits := make(map[string]string, len(lineSeparated))
	for _, line := range lineSeparated {
		ab := strings.Split(line, ")")
		if len(ab) != 2 {
			log.Fatalf("Invalid orbit input: %v", ab)
		}

		orbit, around := ab[0], ab[1]
		orbits[around] = orbit
	}

	sum := 0
	for orbit := range orbits {
		if around, ok := orbits[orbit]; ok {
			Traverse(orbits, around, func(orbit string) { sum++ })
		}
	}
	fmt.Println("Part 1:", sum)

	seen := make(map[string]bool, 0)
	Traverse(orbits, orbits["YOU"], func(orbit string) { seen[orbit] = true })
	Traverse(orbits, orbits["SAN"], func(orbit string) {
		if _, beHere := seen[orbit]; beHere {
			seen[orbit] = false
		} else {
			seen[orbit] = true
		}
	})

	steps := 0
	for _, wantHere := range seen {
		if wantHere {
			steps++
		}
	}
	fmt.Println("Part 2:", steps)
}
