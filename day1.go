package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func day1() error {
	scanner := bufio.NewScanner(os.Stdin)

	extraFuel := 0
	fuelNeeded := 0
	for scanner.Scan() {
		line := scanner.Text()
		module, err := strconv.ParseInt(line, 10, 32)
		if err != nil {
			return err
		}

		needed := fuelForModule(int(module))
		fuelNeeded += needed
		extraFuel +=  fuelForFuel(needed)
	}

	fmt.Println(fmt.Sprintf("Solution: %d", fuelNeeded))
	fmt.Println(fmt.Sprintf("Extra fuel: %d", extraFuel))
	fmt.Println(fmt.Sprintf("Total: %d", fuelNeeded + extraFuel))
	return nil
}

func fuelForModule(mass int) int {
	return int(math.Floor(float64(mass/3))) - 2
}

func fuelForFuel(mass int) int {
	if mass <= 0 {
		return 0
	}

	extraFuel := int(math.Floor(float64(mass/3))) - 2
	if extraFuel < 0 {
		return 0
	}
	return extraFuel + fuelForFuel(extraFuel)
}
