package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func day6() error {
	m := NewMap(os.Stdin)
	fmt.Println("Solution 1:", m.OrbitCount())
	fmt.Println("Solution 2:", m.Route("YOU", "SAN"))

	return nil
}

func NewMap(r io.Reader) Map {
	graph := make(map[string][]string)
	locations := make(map[string]bool)
	parents := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		tokens := strings.Split(t, ")")
		parent, child := tokens[0], tokens[1]

		if _, ok := graph[parent]; !ok {
			graph[parent] = []string{child}
		} else {
			graph[parent] = append(graph[parent], child)
		}

		if _, ok := locations[parent]; !ok {
			locations[parent] = true
		}
		if _, ok := locations[child]; !ok {
			locations[child] = true
		}

		if _, ok := parents[child]; !ok {
			parents[child] = parent
		}
	}
	var locs []string

	for key := range locations {
		locs = append(locs, key)
	}

	m := Map{graph: graph, locations: locs, parents: parents}
	return m
}

type Map struct {
	graph     map[string][]string
	locations []string
	parents   map[string]string
}

func (m *Map) OrbitCountFrom(loc string) int {
	var count int
	for _, child := range m.graph[loc] {
		count += 1 + m.OrbitCountFrom(child)
	}
	return count
}

func (m *Map) OrbitCount() int {
	var oc func(map[string][]string, string) int
	oc = func(g map[string][]string, start string) int {
		var count int
		for _, child := range g[start] {
			count += 1 + oc(g, child)
		}
		return count
	}
	var sum int
	for _, location := range m.locations {
		sum += oc(m.graph, location)
	}

	return sum
}

func (m *Map) Path(from string) []string {
	curr := from
	var path []string
	for {
		parent, ok := m.parents[curr]
		if !ok {
			return path
		}

		path = append(path, parent)
		curr = parent
	}
}

func (m *Map) Route(from, to string) int {
	toRootFrom := reverse(m.Path(from))
	toRootTo := reverse(m.Path(to))

	var commonIdx int
	for i := range toRootFrom {
		if toRootFrom[i] != toRootTo[i] {
			// it is the index one before that is common
			commonIdx = i - 1
			break
		}
	}

	// -1 to count edges, not vertices of the graph
	// another -1 to remove duplicate of the common ancestor
	return len(toRootFrom[commonIdx:]) + len(toRootTo[commonIdx:]) - 1 - 1
}

func reverse(vals []string) []string {
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	return vals
}
