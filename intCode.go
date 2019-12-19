package main

import (
	"fmt"
	"math"
)

type OpCode int
type Mode int

// Parameters
const (
	ParamPosition Mode = iota
	ParamImmediate
	ParamRelative
)

// Instructions
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
	OpAdjustsRB
	OpRet OpCode = 99
)

// NewParam(32102, 0) => Param { mode 1}
// NewParam(32102, 1) => Param { mode 2}
// NewParam(32102, 2) => Param { mode 3}
func NewMode(code, position int) Mode {
	return Mode(int(code/int(math.Pow(10.0, float64(position+2)))) % 10)
}

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
	case OpInput, OpOutput, OpAdjustsRB:
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
	RB     int
	memory []int
	Output chan int
	Input  chan int
	Halt   chan error
}

func NewIntCode(memory []int) IntCode {
	machine := IntCode{
		0,
		0,
		make([]int, len(memory)),
		make(chan int),
		make(chan int),
		make(chan error)}
	copy(machine.memory, memory)
	return machine
}

func (m *IntCode) ReadOp() Op {
	op := Op{m.memory[m.IP]}
	m.IP++
	return op
}

func (m *IntCode) ReadParamAddrs(op Op) []int {
	modes := op.Modes()
	addrs := make([]int, len(modes))
	for i, mode := range modes {
		address := 0
		switch mode {
		case ParamPosition:
			address = m.memory[m.IP]
		case ParamImmediate:
			address = m.IP
		case ParamRelative:
			address = m.RB + m.memory[m.IP]
		}

		// expand memory in out of range case
		if address >= len(m.memory) {
			expandAt := make([]int, address-len(m.memory)+1)
			m.memory = append(m.memory, expandAt...)
		}

		addrs[i] = address
		m.IP++
	}

	return addrs
}

func (m *IntCode) Execute(op OpCode, address []int) bool {
	switch op {
	case OpSum:
		ls, rs, dst := address[0], address[1], address[2]
		m.memory[dst] = m.memory[ls] + m.memory[rs]
	case OpMultiply:
		ls, rs, dst := address[0], address[1], address[2]
		m.memory[dst] = m.memory[ls] * m.memory[rs]
	case OpInput:
		m.memory[address[0]] = <-m.Input
	case OpOutput:
		m.Output <- m.memory[address[0]]
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
	case OpAdjustsRB:
		m.RB += m.memory[address[0]]
	case OpRet:
		return true
	}
	return false
}

func (m *IntCode) Run() {
	go func() {
		for len(m.memory) > m.IP {
			if op := m.ReadOp(); m.Execute(op.OpCode(), m.ReadParamAddrs(op)) {
				m.Halt <- nil
				return
			}
		}
		m.Halt <- fmt.Errorf("Ret instruction is missed: %v", m.memory)
	}()
}
