package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Perms(input []int, out chan<- []int) {
	defer close(out)

	var inner func([]int, int)
	inner = func(arr []int, n int) {
		if n == 1 {
			output := make([]int, len(arr))
			copy(output, arr)
			out <- output
		} else {
			for i := 0; i < n; i++ {
				inner(arr, n-1)
				if n%2 == 1 {
					arr[i], arr[n-1] = arr[n-1], arr[i]
				} else {
					arr[0], arr[n-1] = arr[n-1], arr[0]
				}
			}
		}
	}

	arr := make([]int, len(input))
	copy(arr, input)
	inner(arr, len(arr))
}

func GenPerms(arr []int) chan []int {
	perms := make(chan []int)
	go Perms(arr, perms)
	return perms
}

type Controller func([]int, chan []int, chan int)

func Controller04(program []int, in chan []int, out chan int) {
	for settings := range in {
		res := 0
		for _, setting := range settings {
			vm := NewIntCode(program)
			vm.Run()
			vm.Input <- setting
			vm.Input <- res
			res = <-vm.Output
		}
		go func(res int) { out <- res }(res)
	}
}

func Controller59(program []int, in chan []int, out chan int) {
Serve:
	for settings := range in {
		res := 0

		VMs := make([]IntCode, len(settings))
		for i, setting := range settings {
			VMs[i] = NewIntCode(program)
			vm := VMs[i]
			vm.Run()
			vm.Input <- setting
			vm.Input <- res
			res = <-vm.Output
		}

		var err error
		for {
			for _, vm := range VMs {
				select {
				case vm.Input <- res:
				case err = <-vm.Halt:
					if err != nil {
						log.Fatal(err)
					}
					go func(res int) { out <- res }(res)
					continue Serve
				}

				res = <-vm.Output
			}
		}
	}
}

func Workers(program []int, controller Controller, count int) (chan []int, chan int) {
	in := make(chan []int)
	out := make(chan int)
	for i := 0; i < count; i++ {
		go controller(program, in, out)
	}
	return in, out
}

func Run(program, settings []int, controller Controller) int {
	in, out := Workers(program, controller, 4)
	wg := sync.WaitGroup{}

	for perm := range GenPerms(settings) {
		wg.Add(1)
		in <- perm
	}

	max := 0
	go func() {
		for res := range out {
			wg.Done()
			if max < res {
				max = res
			}
		}
	}()
	wg.Wait()
	return max
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide a file name as argument")
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Can't read input: ", input)
	}

	tokens := strings.Split(strings.TrimSpace(string(input)), ",")
	program := make([]int, len(tokens))
	for i, token := range tokens {
		program[i], err = strconv.Atoi(token)
		if err != nil {
			log.Fatal("Can't parse input: ", err)
		}
	}

	fmt.Println("Part 1:", Run(program, []int{0, 1, 2, 3, 4}, Controller04))
	fmt.Println("Part 2:", Run(program, []int{5, 6, 7, 8, 9}, Controller59))
}
