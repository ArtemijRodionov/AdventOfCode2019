package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	IP     int
	memory []int
}

func newMachine(memory []int) *Machine {
	machineMemory := make([]int, len(memory))
	copy(machineMemory, memory)
	return &Machine{0, machineMemory}
}

func (m Machine) readParams() (int, int, int) {
	return m.memory[m.IP+1], m.memory[m.IP+2], m.memory[m.IP+3]
}

func (m *Machine) compute() []int {
	count := len(m.memory)
	for m.IP < count {
		switch m.memory[m.IP] {
		case 1:
			ls, rs, dst := m.readParams()
			m.memory[dst] = m.memory[ls] + m.memory[rs]
			m.IP += 4
		case 2:
			ls, rs, dst := m.readParams()
			m.memory[dst] = m.memory[ls] * m.memory[rs]
			m.IP += 4
		case 99:
			return m.memory
		default:
			log.Fatal("Something goes wrong", m)
		}
	}

	return []int{}
}

func main() {
	c, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	tokens := strings.Split(strings.Trim(string(c), "\n"), ",")
	memory := make([]int, len(tokens))
	for i, code := range tokens {
		memory[i], err = strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
	}

	m := newMachine(memory)
	m.memory[1] = 12
	m.memory[2] = 2
	fmt.Println("1:", m.compute()[0])

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			m := newMachine(memory)
			m.memory[1] = i
			m.memory[2] = j
			if res := m.compute(); res[0] == 19690720 {
				fmt.Println("2:", i, j)
				return
			}
		}
	}
}
