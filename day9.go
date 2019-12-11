package main

import (
	"bufio"
	"fmt"
	"os"
)

func day9() error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()

	{
		intcode, err := NewIntcode(line, []int64{1})
		if err != nil {
			return err
		}

		out := intcode.Exec()
		fmt.Println("Solution 1:", out)
	}

	return nil
}
