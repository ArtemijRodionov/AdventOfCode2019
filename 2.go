package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
    "os"
)

type Machine struct {
	pos     int
	opCodes []int
}

func newMachine(opCodes []int) *Machine {
	return &Machine{0, opCodes}
}

func (m Machine) readOpArgs() (int, int, int) {
	return m.opCodes[m.pos+1], m.opCodes[m.pos+2], m.opCodes[m.pos+3]
}

func (m *Machine) compute() []int {
	count := len(m.opCodes)
	for m.pos < count {
		switch m.opCodes[m.pos] {
		case 1:
			ls, rs, dst := m.readOpArgs()
			m.opCodes[dst] = m.opCodes[ls] + m.opCodes[rs]
			m.pos += 4
		case 2:
			ls, rs, dst := m.readOpArgs()
			m.opCodes[dst] = m.opCodes[ls] * m.opCodes[rs]
			m.pos += 4
		case 99:
			return m.opCodes
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
	opCodes := make([]int, len(tokens))
	for i, code := range tokens {
		opCodes[i], err = strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
	}

	m := newMachine(opCodes)
	fmt.Println(m.compute())
}
