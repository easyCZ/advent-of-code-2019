package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
)

var solutions = map[int]func() error {
}

func main() {
	var day int

	cmd := cobra.Command{
		Use:                        "advent-of-code-2019",
		RunE: func(cmd *cobra.Command, args []string) error {

			if day < 1 || day > 25 {
				return errors.New("day must be between 1 and 24")
			}

			solution, ok := solutions[day]
			if !ok {
				return errors.New(fmt.Sprintf("solution for day %d not definied", day))
			}

			if err := solution(); err != nil {
				return errors.Wrap(err, "solution failed")
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&day, "day", "d", 0, "which day to run")

	if err := cmd.Execute(); err != nil {
		fmt.Println(fmt.Sprintf("Failed: %v", err))
	}
}
