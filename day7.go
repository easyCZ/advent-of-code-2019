package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func day7() error {
	scanner := bufio.NewScanner(os.Stdin)

	var program string
	for scanner.Scan() {
		program = scanner.Text()
	}

	thrust, err := maxThrust(program)
	if err != nil {
		return err
	}
	fmt.Println("Solution 1:", thrust)

	thrustWithFeedback, err := maxThrustersWithFeedback(program)
	if err != nil {
		return err
	}

	fmt.Println("solution 2:", thrustWithFeedback)
	return nil
}

func maxThrust(program string) (int64, error) {
	var outputs []int64
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				for l := 0; l < 5; l++ {
					for m := 0; m < 5; m++ {

						if len(set([]int{i, j, k, l, m})) != 5 {
							continue
						}

						var amplifiers []*Intcode
						for i := 0; i < 5; i++ {
							intcode, err := NewIntcode(program, nil)
							if err != nil {
								return 0, err
							}
							amplifiers = append(amplifiers, intcode)
						}

						amplifiers[0].input = []int64{int64(i), 0}
						out0 := amplifiers[0].Exec()

						amplifiers[1].input = []int64{int64(j), int64(out0[0])}
						out1 := amplifiers[1].Exec()

						amplifiers[2].input = []int64{int64(k), int64(out1[0])}
						out2 := amplifiers[2].Exec()

						amplifiers[3].input = []int64{int64(l), int64(out2[0])}
						out3 := amplifiers[3].Exec()

						amplifiers[4].input = []int64{int64(m), int64(out3[0])}
						out4 := amplifiers[4].Exec()

						outputs = append(outputs, out4...)
					}
				}
			}
		}
	}

	return max(outputs), nil
}

func maxThrustersWithFeedback(program string) (int64, error) {
	var outputs []int64
	for i := 4; i < 10; i++ {
		for j := 4; j < 10; j++ {
			for k := 4; k < 10; k++ {
				for l := 4; l < 10; l++ {
					for m := 4; m < 10; m++ {

						if len(set([]int{i, j, k, l, m})) != 5 {
							continue
						}

						var amplifiers []*Intcode
						for i := 0; i < 5; i++ {
							intcode, err := NewIntcode(program, nil)
							if err != nil {
								return 0, err
							}
							amplifiers = append(amplifiers, intcode)
						}

						amplifiers[0].input = []int64{int64(i)}
						amplifiers[1].input = []int64{int64(j)}
						amplifiers[2].input = []int64{int64(k)}
						amplifiers[3].input = []int64{int64(l)}
						amplifiers[4].input = []int64{int64(m)}

						for _, amplifier := range amplifiers {
							amplifier.ExecUntil(StoreOpcode)
						}

						var lastOutput int64
						index := 0
						for !amplifiers[4].halted {
							amplifiers[index].AddInput(lastOutput)
							amplifiers[index].ExecUntil(OutputOpcode)
							out := amplifiers[index].Step(amplifiers[index].CurrentInstruction())
							if out != nil {
								lastOutput = int64(*out)
							}

							index = (index + 1) % 5
						}

						outputs = append(outputs, lastOutput)
					}
				}
			}
		}
	}

	return max(outputs), nil
}

func max(vals []int64) int64 {
	m := int64(math.MinInt64)
	for _, val := range vals {
		if val > m {
			m = val
		}
	}

	return m
}

func set(vals []int) map[int]struct{} {
	s := make(map[int]struct{})
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}
