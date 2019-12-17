package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func Run(i int, program []int) error {
	var (
		res int
		err error
	)

	m := NewIntCode(program)
	m.Run()
	m.Input <- i
	for {
		select {
		case res = <-m.Output:
			fmt.Println("Output:", res)
		case err = <-m.Halt:
			return err
		}
	}
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	tokens := strings.Split(strings.TrimSpace(string(input)), ",")
	memory := make([]int, len(tokens))
	for i, code := range tokens {
		memory[i], err = strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
	}

	Run(1, memory)
	Run(5, memory)
}
