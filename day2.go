package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func day2() error {
	reader := bufio.NewReader(os.Stdin)
	cmdString, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	{
		// part 1
		intcode, err := NewIntcode(strings.TrimSpace(cmdString), nil)
		if err != nil {
			return err
		}
		intcode.memory.Set(1, 12)
		intcode.memory.Set(2, 2)
		mem := intcode.Exec()
		fmt.Println(fmt.Sprintf("Solution 1: %d", mem[0]))
	}
	//{
	//	intcode, err := NewIntcode(strings.TrimSpace(cmdString), nil)
	//	if err != nil {
	//		return err
	//	}
	//
	//
	//	noun, verb := findVal(vals)
	//	fmt.Println("Noun:", noun)
	//	fmt.Println("Opcode:", verb)
	//	fmt.Println(fmt.Sprintf("Solution 2: %d", 100*noun+verb))
	//}

	return nil
}

//func findVal(program []int64) (int, int) {
//	for i := 0; i <= 99; i++ {
//		for j := 0; j <= 99; j++ {
//			p2 := append(program[:0:0], program...)
//			p2[1] = int64(i)
//			p2[2] = int64(j)
//
//			p := Intcode{memory: p2}
//			mem := p.Exec()
//
//			if mem[0] == 19690720 {
//				return i, j
//			}
//		}
//	}
//
//	return 0, 0
//}

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

func exec(program []int) []int {
	opIndex := 0
	for program[opIndex] != 99 {
		program = step(program, opIndex)
		opIndex += 4
	}

	return program
}
