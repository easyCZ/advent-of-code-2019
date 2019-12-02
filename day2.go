package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

func day2() error {
	scanner := bufio.NewScanner(os.Stdin)

	var vals []int
	for scanner.Scan() {
		t := scanner.Text()
		tokens := strings.Split(t, ",")
		for _, token := range tokens {
			i, err := strconv.ParseInt(token, 10, 32)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to parse entry %v", token))
			}

			vals = append(vals, int(i))
		}
	}

	{
		// part 1
		programAlarm := append(vals[:0:0], vals...)
		programAlarm[1] = 12
		programAlarm[2] = 2
		p := Intcode{memory: programAlarm}
		mem := p.Exec()
		fmt.Println(fmt.Sprintf("Solution 1: %d", mem[0]))
	}
	{
		noun, verb := findVal(vals)
		fmt.Println("Noun:", noun)
		fmt.Println("Verb:", verb)
		fmt.Println(fmt.Sprintf("Solution 2: %d", 100*noun+verb))
	}

	return nil
}

func findVal(program []int) (int, int) {
	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			p2 := append(program[:0:0], program...)
			p2[1] = i
			p2[2] = j

			p := Intcode{memory: p2}
			mem := p.Exec()

			if mem[0] == 19690720 {
				return i, j
			}
		}
	}

	return 0, 0
}

func step(program []int, opIndex int) []int {
	op := program[opIndex]
	switch op {
	case 1:
		// ADD
		left := program[program[opIndex+1]]
		right := program[program[opIndex+2]]
		target := program[opIndex+3]
		sum := left + right
		//fmt.Println(left, right, target, sum)
		program[target] = sum
	case 2:
		// Multiple
		left := program[program[opIndex+1]]
		right := program[program[opIndex+2]]
		target := program[opIndex+3]
		product := left * right
		//fmt.Println(left, right, target, product)
		program[target] = product
	case 99:
	// halt
	default:
		panic(fmt.Sprintf("unknown code %v encountered", op))
	}

	return program
}

func exec(program [] int) []int {
	opIndex := 0
	for program[opIndex] != 99 {
		program = step(program, opIndex)
		opIndex += 4
	}

	return program
}

type Verb int

const (
	Unknown  Verb = 0
	Add      Verb = 1
	Multiply Verb = 2
	Halt     Verb = 99
)

type Intcode struct {
	memory []int
	cursor int
}

func (i *Intcode) Step() Verb {
	op := i.memory[i.cursor]
	switch op {
	case int(Add): // Add
		left := i.memory[i.memory[i.cursor+1]]
		right := i.memory[i.memory[i.cursor+2]]
		target := i.memory[i.cursor+3]
		sum := left + right
		i.memory[target] = sum
		i.cursor += 4
		return Add
	case int(Multiply): // Multiply
		left := i.memory[i.memory[i.cursor+1]]
		right := i.memory[i.memory[i.cursor+2]]
		target := i.memory[i.cursor+3]
		product := left * right
		i.memory[target] = product
		i.cursor += 4
		return Multiply
	case int(Halt): // halt
		i.cursor += 1
		return Halt
	default:
		panic(fmt.Sprintf("unknown code %v encountered", op))
	}

	return Unknown
}

func (i *Intcode) Exec() []int {
	for i.memory[i.cursor] != int(Halt) {
		i.Step()
	}

	return i.memory
}
