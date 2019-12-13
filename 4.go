package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Match func(int) bool

func isMatching1(number int) bool {
	digits := strconv.Itoa(number)
	ok := false

	for i := range digits[1:] {
		prev, next := digits[i], digits[i+1]
		if prev > next {
			return false
		}

		if prev == next {
			ok = true
		}
	}

	return ok
}

func isMatching2(number int) bool {
	digits := strconv.Itoa(number)

	ok := false
	adjacent := 0

	for i := range digits[1:] {
		prev, next := digits[i], digits[i+1]
		if prev > next {
			return false
		}

		if prev == next {
			adjacent++
			continue
		}

		if adjacent == 1 {
			ok = true
		}
		adjacent = 0
	}

	lastOk := adjacent == 1
	return ok || lastOk
}

func run(low, high int, isMatching Match) int {
	matched := 0

	for number := low; number <= high; number++ {
		if isMatching(number) {
			matched++
		}
	}

	return matched
}

func main() {
	digitsRange, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	digits := strings.Split(string(digitsRange), "-")
	if len(digits) != 2 {
		log.Fatal("Invalid range")
	}

	numbers := [2]int{}
	for i := 0; i < 2; i++ {
		numbers[i], err = strconv.Atoi(digits[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("1:", run(numbers[0], numbers[1], isMatching1))
	fmt.Println("2:", run(numbers[0], numbers[1], isMatching2))
}
