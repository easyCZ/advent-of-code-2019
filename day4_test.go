package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCandidate(t *testing.T) {
	assert.Equal(t, false, isCandidate(111111))
	assert.Equal(t, false, isCandidate(223450))
	assert.Equal(t, false, isCandidate(123789))
}
