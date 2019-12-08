package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseImage(t *testing.T) {
	img, err := parseImage("123456789012", 3, 2)
	assert.NoError(t, err)

	assert.Equal(t, []Layer{
		[][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
		[][]int{
			{7, 8, 9},
			{0, 1, 2},
		},
	}, img)

}
