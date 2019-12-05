package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func day3() error {
	scanner := bufio.NewScanner(os.Stdin)

	var wires []Wire
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")

		var paths []Path
		for _, token := range tokens {
			paths = append(paths, NewPath(token))
		}

		wires = append(wires, Wire{paths: paths})
	}

	fmt.Println("Solution 1: ", nearestIntersection(wires[0], wires[1]))
	fmt.Println("Solution 2: ", minSteps(wires[0], wires[1]))

	return nil
}

type Wire struct {
	paths []Path
}

func (w *Wire) Trace() []Point {
	var points []Point
	var x, y int

	for _, path := range w.paths {
		switch path.direction {
		case Left:
			for i := 0; i < path.units; i++ {
				x -= 1
				points = append(points, Point{x: x, y: y})
			}
		case Right:
			for i := 0; i < path.units; i++ {
				x += 1
				points = append(points, Point{x: x, y: y})
			}
		case Up:
			for i := 0; i < path.units; i++ {
				y += 1
				points = append(points, Point{x: x, y: y})
			}
		case Down:
			for i := 0; i < path.units; i++ {
				y -= 1
				points = append(points, Point{x: x, y: y})
			}
		}
	}

	return points
}

type Direction int

const (
	Up    Direction = 1
	Left  Direction = 2
	Down  Direction = 3
	Right Direction = 4
)

type Path struct {
	direction Direction
	units     int
}

func NewPath(s string) Path {
	units, err := strconv.ParseInt(s[1:], 10, 32)
	if err != nil {
		panic(fmt.Sprintf("failed to parse path '%v'", s))
	}

	switch s[0] {
	case 'U':
		return Path{direction: Up, units: int(units)}
	case 'D':
		return Path{direction: Down, units: int(units)}
	case 'L':
		return Path{direction: Left, units: int(units)}
	case 'R':
		return Path{direction: Right, units: int(units)}
	default:
		panic("UnknownOpcode path")
	}
}

type Point struct {
	x int
	y int
}

func (c *Point) ManhattanDistance(from Point) int {
	return int(math.Abs(float64(c.x-from.x)) + math.Abs(float64(c.y-from.y)))
}

func intersections(w1 Wire, w2 Wire) []Point {
	pts1 := w1.Trace()
	pts2 := w2.Trace()

	var crosses []Point
	for _, p1 := range pts1 {
		for _, p2 := range pts2 {
			if p1.x == 0 && p1.y == 0 {
				continue
			}
			if p1 == p2 {
				crosses = append(crosses, p1)
			}
		}
	}

	return crosses
}

func nearest(coords []Point) Point {
	min := math.MaxInt32
	var minCoord Point
	origin := Point{}

	for i := 0; i < len(coords); i++ {
		if dist := coords[i].ManhattanDistance(origin); dist < min {
			min = dist
			minCoord = coords[i]
		}
	}

	return minCoord
}

func nearestIntersection(w1 Wire, w2 Wire) int {
	ints := intersections(w1, w2)
	closest := nearest(ints)
	return closest.ManhattanDistance(Point{0, 0})
}

func minSteps(w1 Wire, w2 Wire) int {
	ints := intersections(w1, w2)

	t1 := w1.Trace()
	t2 := w2.Trace()

	min := math.MaxInt32

	for _, intersection := range ints {
		s1 := findFirst(t1, intersection)
		s2 := findFirst(t2, intersection)

		sum := s1 + s2
		if sum < min {
			min = sum
		}
	}

	return min
}

func findFirst(points []Point, p Point) int {
	for i, pp := range points {
		if pp == p {
			return i + 1
		}
	}

	return 0
}
