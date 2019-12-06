package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuelForFuel(t *testing.T) {
	scenarios := []struct {
		in  int
		out int
	}{
		{in: 14, out: 2},
		{in: 1969, out: 966},
		{in: 100756, out: 50346},
	}

	for _, scenario := range scenarios {
		assert.Equal(t, scenario.out, fuelForFuel(scenario.in))
	}
}
