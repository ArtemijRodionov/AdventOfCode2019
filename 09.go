package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func Run(program []int, input int) int {
	vm := NewIntCode(program)
	vm.Run()
	var (
		res int
		err error
	)

	vm.Input <- input

	for {
		select {
		case res = <-vm.Output:
		case err = <-vm.Halt:
			if err != nil {
				log.Fatal("VM error: ", err)
			}
			return res
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an input file name as the first argument")
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Can't read file: ", err)
	}

	tokens := strings.Split(strings.TrimSpace(string(input)), ",")
	program := make([]int, len(tokens))
	for i, token := range tokens {
		program[i], err = strconv.Atoi(token)
		if err != nil {
			log.Fatal("Can't parse a number: ", err)
		}
	}

	fmt.Println("Part 1: ", Run(program, 1))
	fmt.Println("Part 2: ", Run(program, 2))
}
