package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntcodeReadStore(t *testing.T) {
	i, err := NewIntcode("3,0,4,0,99", []int{5})
	assert.NoError(t, err)
	out := i.Exec()
	assert.Equal(t, []int{5}, out)
}

func TestIncodeModes(t *testing.T) {
	i, err := NewIntcode("1002,4,3,4,33", nil)
	assert.NoError(t, err)
	_ = i.Exec()
	assert.Equal(t, []int{1002, 4, 3, 4, 99}, i.memory)
}

func TestIncodeNegative(t *testing.T) {
	i, err := NewIntcode("1101,100,-1,4,0", nil)
	assert.NoError(t, err)
	_ = i.Exec()
	assert.Equal(t, []int{1101, 100, -1, 4, 99}, i.memory)
}

func TestComparisons(t *testing.T) {
	for _, scenario := range []struct {
		program string
		in      []int
		out     []int
	}{
		{
			program: "3,9,8,9,10,9,4,9,99,-1,8",
			in:      []int{8},
			out:     []int{1},
		},
		{
			program: "3,9,8,9,10,9,4,9,99,-1,8",
			in:      []int{7},
			out:     []int{0},
		},
		{
			program: "3,9,7,9,10,9,4,9,99,-1,8",
			in:      []int{7},
			out:     []int{1},
		},
		{
			program: "3,9,7,9,10,9,4,9,99,-1,8",
			in:      []int{8},
			out:     []int{0},
		},
		{
			program: "3,3,1108,-1,8,3,4,3,99",
			in:      []int{8},
			out:     []int{1},
		},
		{
			program: "3,3,1108,-1,8,3,4,3,99",
			in:      []int{7},
			out:     []int{0},
		},
		{
			program: "3,3,1107,-1,8,3,4,3,99",
			in:      []int{7},
			out:     []int{1},
		},
		{
			program: "3,3,1107,-1,8,3,4,3,99",
			in:      []int{8},
			out:     []int{0},
		},

		{
			program: "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			in:      []int{0},
			out:     []int{0},
		},
		{
			program: "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			in:      []int{1},
			out:     []int{1},
		},

		{
			program: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			in:      []int{0},
			out:     []int{0},
		},
		{
			program: "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			in:      []int{1},
			out:     []int{1},
		},
		{
			program: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			in:      []int{7},
			out:     []int{999},
		},
		{
			program: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			in:      []int{8},
			out:     []int{1000},
		},
		{
			program: "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99",
			in:      []int{9},
			out:     []int{1001},
		},
	} {
		i, err := NewIntcode(scenario.program, scenario.in)
		assert.NoError(t, err)
		assert.Equal(t, scenario.out, i.Exec())
	}
}
