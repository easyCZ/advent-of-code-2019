package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Opcode struct {
	id     int
	length int
}

var (
	UnknownOpcode     = Opcode{id: 0}
	AddOpcode         = Opcode{id: 1, length: 4}
	MultiplyOpcode    = Opcode{id: 2, length: 4}
	StoreOpcode       = Opcode{id: 3, length: 2}
	OutputOpcode      = Opcode{id: 4, length: 2}
	JumpIfTrueOpcode  = Opcode{id: 5, length: 3}
	JumpIfFalseOpcode = Opcode{id: 6, length: 3}
	LessThanOpcode    = Opcode{id: 7, length: 4}
	EqualsOpcode      = Opcode{id: 8, length: 4}
	SetRelBaseOpcode  = Opcode{id: 9, length: 2}
	HaltOpcdoe        = Opcode{id: 99, length: 1}
)

type Memory struct {
	m map[int64]int64
}

func (m *Memory) Set(loc int64, val int64) {
	if loc < 0 {
		panic(fmt.Sprintf("Attempted to set negative memroy: %d", loc))
	}
	m.m[loc] = val
}

func (m *Memory) Get(loc int64) int64 {
	if loc < 0 {
		panic(fmt.Sprintf("Requested negative memroy: %d", loc))
	}
	if val, ok := m.m[loc]; ok {
		return val
	}

	return 0
}

type Output int64

type Intcode struct {
	memory  Memory
	cursor  int64
	input   []int64
	halted  bool
	relBase int64
}

func (i *Intcode) readInput() (int64, error) {
	if len(i.input) > 0 {
		val := i.input[0]
		i.input = i.input[1:]
		return val, nil
	}

	return 0, errors.New("no input available")
}

func (i *Intcode) val(param Param) int64 {
	switch param.mode {
	case ImmediateMode:
		return param.val
	case PositionMode:
		return i.memory.Get(param.val)
	default:
		panic(fmt.Sprintf("Unknown mode encountered: %v", param.mode))
	}
}

func (i *Intcode) Step(instruction Instruction) *Output {
	switch instruction.opcode {
	case AddOpcode:

		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val

		sum := left + right
		i.memory.Set(target, sum)
		i.Advance(instruction)
		return nil
	case MultiplyOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val

		product := left * right
		i.memory.Set(target, product)
		i.Advance(instruction)
		return nil
	case StoreOpcode:
		target := instruction.params[0].val

		input, err := i.readInput()
		if err != nil {
			panic(err)
		}
		i.memory.Set(target, input)
		i.Advance(instruction)
		return nil
	case OutputOpcode:
		val := i.val(instruction.params[0])
		out := Output(val)
		i.Advance(instruction)
		return &out

	case JumpIfTrueOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		if left != 0 {
			i.cursor = right
		} else {
			i.Advance(instruction)
		}
		return nil

	case JumpIfFalseOpcode:
		left := i.val(instruction.params[0])
		target := i.val(instruction.params[1])
		if left == 0 {
			i.cursor = target
		} else {
			i.Advance(instruction)
		}
		return nil

	case LessThanOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val

		if left < right {
			i.memory.Set(target, 1)
		} else {
			i.memory.Set(target, 0)
		}
		i.Advance(instruction)
		return nil

	case EqualsOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val
		if left == right {
			i.memory.Set(target, 1)
		} else {
			i.memory.Set(target, 0)
		}
		i.Advance(instruction)
		return nil

	case SetRelBaseOpcode:
		left := instruction.params[0]
		i.relBase += left.val
		i.Advance(instruction)
		return nil

	case HaltOpcdoe:
		i.halted = true
		i.Advance(instruction)
		return nil

	default:
		panic(fmt.Sprintf("unknown instruction %v encountered", instruction))
	}

	return nil
}

func (i *Intcode) Advance(ins Instruction) {
	i.cursor += int64(ins.opcode.length)
}

func (i *Intcode) AddInput(val int64) {
	i.input = append(i.input, val)
}

func (i *Intcode) CurrentInstruction() Instruction {
	s := fmt.Sprintf("%d", i.memory.Get(i.cursor))
	for k := len(s); k < 5; k++ {
		s = "0" + s
	}

	opcode := parseOpcode(s[3:])

	var params []Param

	if opcode.length >= 2 {
		params = append(params, Param{
			val:  i.memory.Get(i.cursor + 1),
			mode: parseMode(string(s[2])),
		})
	}
	if opcode.length >= 3 {
		params = append(params, Param{
			val:  i.memory.Get(i.cursor + 2),
			mode: parseMode(string(s[1])),
		})
	}
	if opcode.length >= 4 {
		params = append(params, Param{
			val:  i.memory.Get(i.cursor + 3),
			mode: parseMode(string(s[0])),
		})
	}

	return Instruction{
		opcode: opcode,
		params: params,
	}
}

func (i *Intcode) ExecUntil(opcode Opcode) []int64 {
	var out []int64

	for {
		instruction := i.CurrentInstruction()
		if instruction.opcode == opcode || instruction.opcode == HaltOpcdoe {
			return out
		}

		step := i.Step(instruction)
		if step != nil {
			out = append(out, int64(*step))
		}
	}
}

func (i *Intcode) Exec() []int64 {
	return i.ExecUntil(HaltOpcdoe)
}

func NewIntcode(s string, input []int64) (*Intcode, error) {
	var vals []int64
	tokens := strings.Split(s, ",")
	for _, token := range tokens {
		i, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse entry %v", token))
		}

		vals = append(vals, int64(i))
	}

	return &Intcode{
		memory: NewMemory(vals),
		cursor: 0,
		input:  input,
	}, nil
}

func NewMemory(vals []int64) Memory {
	mem := Memory{m: map[int64]int64{}}
	for loc, val := range vals {
		mem.Set(int64(loc), val)
	}
	return mem
}

type Mode int

const (
	PositionMode  Mode = 0
	ImmediateMode Mode = 1
	RelativeMode  Mode = 2
)

type Param struct {
	val  int64
	mode Mode
}

type Instruction struct {
	opcode Opcode
	params []Param
}

func (i *Instruction) String() string {
	return fmt.Sprintf("{opcode: %v, params: %v}", i.opcode, i.params)
}

func parseOpcode(s string) Opcode {
	switch s {
	case "02":
		return MultiplyOpcode
	case "01":
		return AddOpcode
	case "03":
		return StoreOpcode
	case "04":
		return OutputOpcode
	case "05":
		return JumpIfTrueOpcode
	case "06":
		return JumpIfFalseOpcode
	case "07":
		return LessThanOpcode
	case "08":
		return EqualsOpcode
	case "09":
		return SetRelBaseOpcode
	case "99":
		return HaltOpcdoe
	default:
		panic(fmt.Sprintf("Unknown opcode: %v", s))
	}
}

func parseMode(s string) Mode {
	switch s {
	case "0":
		return PositionMode
	case "1":
		return ImmediateMode
	case "2":
		return RelativeMode
	default:
		panic(fmt.Sprintf("UnknownOpcode mode: %v", s))
	}
}
