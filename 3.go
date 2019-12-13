package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Dot struct {
	X, Y int
}

var labelDirection = map[byte]Dot{
	'R': {1, 0},
	'U': {0, 1},
	'L': {-1, 0},
	'D': {0, -1},
}

func parse(input string) map[Dot]int {
	moves := strings.Split(input, ",")
	// +1 to store unspecified starting dot (0, 0)
	dots := make(map[Dot]int, len(moves)+1)
	pivot := Dot{}
	cost := 0

	for _, move := range moves {
		dir, ok := labelDirection[move[0]]
		if !ok {
			log.Fatal("Can't parse move direction: ", move)
		}

		value, err := strconv.Atoi(move[1:])
		if err != nil {
			log.Fatal("Can't parse move value: ", move)
		}

		for i := 0; i < value; i++ {
			dots[pivot] = cost
			cost++
			pivot.X += dir.X
			pivot.Y += dir.Y
		}
	}

	return dots
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	rows := strings.Split(strings.Trim(string(input), "\n"), "\n")
	if len(rows) != 2 {
		log.Fatal("Invalid number of lines: ", len(rows))
	}

	a, b := parse(rows[0]), parse(rows[1])
	min := 0.0
	for dotA := range a {
		if _, ok := b[dotA]; ok {
			sum := math.Abs(float64(dotA.X + dotA.Y))
			if min != 0 && min < sum {
				continue
			}
			min = sum
		}
	}
	fmt.Println("1:", min)

	minCost := 0
	minDot := Dot{}
	for dotA, costA := range a {
		if costB, ok := b[dotA]; ok {
			cost := costA + costB
			if minCost != 0 && cost > minCost {
				continue
			}
			minCost = cost
			minDot = dotA
		}
	}
	fmt.Println("2:", minDot, minCost)
}
