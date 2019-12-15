package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type OpCode int
type Mode int

type Instruction interface {
	OpCode() OpCode
	Modes() []Mode
}

type VM interface {
	Address(Mode) int
	Execute(Instruction)
	Run()
}

// Parameters
const (
	ParamPosition Mode = iota
	ParamImmediate
)

// NewParam(32102, 0) => Param { mode 1}
// NewParam(32102, 1) => Param { mode 2}
// NewParam(32102, 2) => Param { mode 3}
func NewMode(code, position int) Mode {
	return Mode(int(code/int(math.Pow(10.0, float64(position+2)))) % 10)
}

// Instruction
const (
	// Start from 1
	OpSum OpCode = iota + 1
	OpMultiply
	OpInput
	OpOutput
	OpJmpIfTrue
	OpJmpIfFalse
	OpLessThan
	OpEquals
	OpRet OpCode = 99
)

type Op struct {
	code int
}

func (o Op) OpCode() OpCode {
	return OpCode(o.code % 100)
}

func (o Op) Modes() []Mode {
	length := 0
	switch o.OpCode() {
	case OpSum, OpMultiply, OpLessThan, OpEquals:
		length = 3
	case OpJmpIfTrue, OpJmpIfFalse:
		length = 2
	case OpInput, OpOutput:
		length = 1
	case OpRet:
		length = 0
	}

	params := make([]Mode, length)
	for i := range params {
		params[i] = NewMode(o.code, i)
	}
	return params
}

// VM
type IntCode struct {
	IP     int
	memory []int
}

func NewIntCode(memory []int) IntCode {
	machine := IntCode{0, make([]int, len(memory))}
	copy(machine.memory, memory)
	return machine
}

func (m *IntCode) Address(mode Mode) int {
	address := 0
	switch mode {
	case ParamPosition:
		address = m.memory[m.IP]
	case ParamImmediate:
		address = m.IP
	}
	return address
}

func (m *IntCode) Execute(inst Instruction) bool {
	modes := inst.Modes()
	address := make([]int, len(modes))
	for i, mode := range modes {
		address[i] = m.Address(mode)
		m.IP++
	}

	switch inst.OpCode() {
	case OpSum:
		ls, rs, dst := address[0], address[1], address[2]
		m.memory[dst] = m.memory[ls] + m.memory[rs]
	case OpMultiply:
		ls, rs, dst := address[0], address[1], address[2]
		m.memory[dst] = m.memory[ls] * m.memory[rs]
	case OpInput:
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		number, err := strconv.Atoi(input)
		if err != nil {
			log.Fatal(err)
		}
		m.memory[address[0]] = number
	case OpOutput:
		fmt.Println(m.memory[address[0]])
	case OpJmpIfTrue:
		if m.memory[address[0]] != 0 {
			m.IP = m.memory[address[1]]
		}
	case OpJmpIfFalse:
		if m.memory[address[0]] == 0 {
			m.IP = m.memory[address[1]]
		}
	case OpLessThan:
		ls, rs, dst := address[0], address[1], address[2]
		toSet := 0
		if m.memory[ls] < m.memory[rs] {
			toSet = 1
		}
		m.memory[dst] = toSet
	case OpEquals:
		ls, rs, dst := address[0], address[1], address[2]
		toSet := 0
		if m.memory[ls] == m.memory[rs] {
			toSet = 1
		}
		m.memory[dst] = toSet
	case OpRet:
		return true
	}
	return false
}

func (m *IntCode) Run() {
	for len(m.memory) > m.IP {
		op := Op{m.memory[m.IP]}
		m.IP++
		if m.Execute(op) {
			return
		}
	}
	log.Fatal("Ret instruction is missed")
}

func main() {
	input, err := ioutil.ReadFile("5.txt")
	if err != nil {
		log.Fatal(err)
	}

	tokens := strings.Split(strings.Trim(string(input), "\n"), ",")
	memory := make([]int, len(tokens))
	for i, code := range tokens {
		memory[i], err = strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
	}

	m := NewIntCode(memory)
	m.Run()
}
