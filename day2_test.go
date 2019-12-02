package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExec(t *testing.T) {
	for _, scenario := range []struct {
		in  []int
		out []int
	}{
		{in: []int{1, 0, 0, 0, 99}, out: []int{2, 0, 0, 0, 99}},
		{in: []int{2, 3, 0, 3, 99}, out: []int{2, 3, 0, 6, 99}},
		{in: []int{2, 4, 4, 5, 99, 0}, out: []int{2, 4, 4, 5, 99, 9801}},
		{in: []int{1, 1, 1, 4, 99, 5, 6, 0, 99}, out: []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
	} {
		assert.Equal(t, scenario.out, exec(scenario.in))
	}
}
