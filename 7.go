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
            for i:= 0; i < n; i++ {
                inner(arr, n - 1)
                if n % 2 == 1 {
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
            output := vm.Run([]int{setting, res})
            res = output[0]
        }
        out <- res
    }
}

func Controller59(program []int, in chan []int, out chan int) {
    for settings := range in {
        res := 0
        vm := NewIntCode(program)
        for _, setting := range settings {
            output := vm.Run([]int{setting, res})
            res = output[0]
        }
        out <- res
    }
}

func Runs(program []int, controller Controller, count int) (chan []int, chan int) {
    in := make(chan []int)
    out := make(chan int)
    for i := 0; i < count; i++ {
        go controller(program, in, out)
    }
    return in, out
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
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

    // run 4 workers
    in, out := Runs(program, Controller04, 4)
    wg := sync.WaitGroup{}

    for perm := range GenPerms([]int{0, 1, 2, 3, 4}) {
        wg.Add(1)
        go (func(perm []int) { in <- perm })(perm)
    }

    max := 0
    go (func() {
        for res := range out {
            wg.Done()
            if max < res {
                max = res
            }
        }
    })()
    wg.Wait()
    fmt.Println("Part 1:", max)


    // run 4 workers
    in, out = Runs(program, Controller59, 4)
    wg = sync.WaitGroup{}

    for perm := range GenPerms([]int{5,6,7,8,9}) {
        wg.Add(1)
        go (func(perm []int) { in <- perm })(perm)
    }

    max = 0
    go (func() {
        for res := range out {
            wg.Done()
            if max < res {
                max = res
            }
        }
    })()
    wg.Wait()
    fmt.Println("Part 2:", max)
}
