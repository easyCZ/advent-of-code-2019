package main

import (
	"fmt"
	"strconv"
)

func day4() error {

	from, to := 359282, 820401
	var candidates []int

	for i := from; i <= to; i++ {
		if isIncreasing(i) && hasSameAdjacent(i) {
			candidates = append(candidates, i)
		}
	}

	fmt.Println(fmt.Sprintf("Solution 1: %d", len(candidates)))

	return nil
}

func isCandidate(n int) bool {
	return isIncreasing(n) && hasSameAdjacent(n)
}

func isIncreasing(n int) bool {
	vals := digits(n)

	for i := 0; i < len(vals)-1; i++ {
		if vals[i+1] < vals[i] {
			return false
		}
	}

	return true
}

func hasSameAdjacent(n int) bool {
	s := fmt.Sprintf("%d", n)
	set := make(map[string]int)

	for _, i := range s {
		if _, ok := set[string(i)]; !ok {
			set[string(i)] = 0
		}

		set[string(i)] += 1
	}

	for _, val := range set {
		if val == 2 {
			return true
		}
	}

	return false
}

func digits(n int) []int {
	s := fmt.Sprintf("%d", n)

	var vals []int
	for _, i := range s {
		k, err := strconv.ParseInt(fmt.Sprintf("%d", i), 10, 32)
		if err != nil {
			panic("failed to parse")
		}

		vals = append(vals, int(k))
	}
	return vals
}
