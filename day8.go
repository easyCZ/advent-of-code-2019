package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func day8() error {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	text := scanner.Text()

	layers, err := parseImage(text, 25, 6)
	if err != nil {
		return err
	}

	fmt.Println("Solution 1:", checksum(layers))

	rendered := render(layers, 25, 6)
	fmt.Println("Solution 2:")
	for _, row := range rendered {
		fmt.Println(row)
	}

	return nil
}

type Layer [][]int

type Image []Layer

func render(img []Layer, width, height int) [][]int {
	out := make([][]int, height)
	for i := 0; i < height; i++ {
		out[i] = make([]int, width)
	}

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {

			var vals []int
			for _, layer := range img {
				vals = append(vals, layer[row][col])
			}

			out[row][col] = renderCell(vals)
		}
	}

	return out
}

func renderCell(layers []int) int {
	for _, l := range layers {
		if l == 0 {
			return l
		}
		if l == 1 {
			return l
		}
	}

	return 2
}

func parseImage(in string, width, height int) ([]Layer, error) {
	var layers []Layer

	for l := 0; l < len(in)/(width*height); l++ {
		layer := make([][]int, 0)

		for h := 0; h < height; h++ {
			layer = append(layer, make([]int, width))
		}

		for h := 0; h < height; h++ {
			for w := 0; w < width; w++ {
				index := l*(width*height) + h*width + w
				val, err := strconv.ParseInt(string(in[index]), 10, 32)
				if err != nil {
					return nil, err
				}
				layer[h][w] = int(val)
			}
		}

		layers = append(layers, layer)
	}

	return layers, nil
}

func (l Layer) countDigits(d int) int {
	c := 0
	for i := 0; i < len(l); i++ {
		for j := 0; j < len(l[0]); j++ {
			if l[i][j] == d {
				c += 1
			}
		}
	}
	return c
}

func checksum(layers []Layer) int {
	zeros := math.MaxInt32
	minZeros := layers[0]

	for _, layer := range layers {
		z := layer.countDigits(0)
		if z < zeros {
			zeros = z
			minZeros = layer
		}
	}

	return minZeros.countDigits(1) * minZeros.countDigits(2)
}
