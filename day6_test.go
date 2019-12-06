package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap_OrbitCount(t *testing.T) {
	in := `COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
`
	m := NewMap(bytes.NewBufferString(in))
	assert.Equal(t, 42, m.OrbitCount())
}
