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
	HaltOpcdoe        = Opcode{id: 99, length: 1}
)

type Output int

type Intcode struct {
	memory []int
	cursor int
	input  []int
}

func (i *Intcode) readInput() (int, error) {
	if len(i.input) > 0 {
		val := i.input[0]
		i.input = i.input[1:]
		return val, nil
	}

	return 0, errors.New("no input available")
}

func (i *Intcode) val(param Param) int {
	switch param.mode {
	case ImmediateMode:
		return param.val
	case PositionMode:
		return i.memory[param.val]
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
		i.memory[target] = sum
		i.Advance(instruction)
		return nil
	case MultiplyOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val

		product := left * right
		i.memory[target] = product
		i.Advance(instruction)
		return nil
	case StoreOpcode:
		target := instruction.params[0].val

		input, err := i.readInput()
		if err != nil {
			panic(err)
		}
		i.memory[target] = input
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
			i.memory[target] = 1
		} else {
			i.memory[target] = 0
		}
		i.Advance(instruction)
		return nil

	case EqualsOpcode:
		left := i.val(instruction.params[0])
		right := i.val(instruction.params[1])
		target := instruction.params[2].val
		if left == right {
			i.memory[target] = 1
		} else {
			i.memory[target] = 0
		}
		i.Advance(instruction)
		return nil
	case HaltOpcdoe:
		i.Advance(instruction)
		return nil
	default:
		panic(fmt.Sprintf("unknown instruction %v encountered", instruction))
	}

	return nil
}

func (i *Intcode) Advance(ins Instruction) {
	i.cursor += ins.opcode.length
}

func (i *Intcode) CurrentInstruction() Instruction {
	s := fmt.Sprintf("%d", i.memory[i.cursor])
	for k := len(s); k < 5; k++ {
		s = "0" + s
	}

	opcode := parseOpcode(s[3:])

	var params []Param

	if opcode.length >= 2 {
		params = append(params, Param{
			val:  i.memory[i.cursor+1],
			mode: parseMode(string(s[2])),
		})
	}
	if opcode.length >= 3 {
		params = append(params, Param{
			val:  i.memory[i.cursor+2],
			mode: parseMode(string(s[1])),
		})
	}
	if opcode.length >= 4 {
		params = append(params, Param{
			val:  i.memory[i.cursor+3],
			mode: parseMode(string(s[0])),
		})
	}

	return Instruction{
		opcode: opcode,
		params: params,
	}
}

func (i *Intcode) Exec() []int {
	var out []int

	for {
		instruction := i.CurrentInstruction()
		if instruction.opcode == HaltOpcdoe {
			return out
		}

		step := i.Step(instruction)
		if step != nil {
			out = append(out, int(*step))
		}
	}
}

func NewIntcode(s string, input []int) (*Intcode, error) {
	var vals []int
	tokens := strings.Split(s, ",")
	for _, token := range tokens {
		i, err := strconv.ParseInt(token, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse entry %v", token))
		}

		vals = append(vals, int(i))
	}

	return &Intcode{
		memory: vals,
		cursor: 0,
		input:  input,
	}, nil
}

type Mode int

const (
	PositionMode  Mode = 0
	ImmediateMode Mode = 1
)

type Param struct {
	val  int
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
	default:
		panic(fmt.Sprintf("UnknownOpcode mode: %v", s))
	}
}
