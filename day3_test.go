package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinDist(t *testing.T) {
	for _, scenario := range []struct {
		w1   Wire
		w2   Wire
		dist int
	}{
		{
			w1: Wire{
				paths: []Path{
					NewPath("R8"), NewPath("U5"), NewPath("L5"), NewPath("D3"),
				},
			},
			w2: Wire{
				paths: []Path{
					NewPath("U7"), NewPath("R6"), NewPath("D4"), NewPath("L4"),
				},
			},
			dist: 6,
		},
	} {
		assert.Equal(t, scenario.dist, nearestIntersection(scenario.w1, scenario.w2))
	}
}

func TestMinSteps(t *testing.T) {
	w1 := Wire{
		paths: []Path{
			NewPath("R8"), NewPath("U5"), NewPath("L5"), NewPath("D3"),
		},
	}
	w2 := Wire{
		paths: []Path{
			NewPath("U7"), NewPath("R6"), NewPath("D4"), NewPath("L4"),
		},
	}

	assert.Equal(t, 30, minSteps(w1, w2))
}
