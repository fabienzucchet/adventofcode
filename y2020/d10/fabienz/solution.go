package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := []int{0}

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		numbers = append(numbers, number)
	}

	_, err = fmt.Fprintf(answer, "%d", computeJolt(sortSlice(numbers)))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := []int{0}

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		numbers = append(numbers, number)
	}

	computeJolt(sortSlice(numbers))

	computed := make(map[int]int)

	_, err = fmt.Fprintf(answer, "%d", countCombination(0, numbers[1:], computed))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func sortSlice(s []int) []int {
	n := len(s)

	for i := 1; i < n; i++ {
		for j := 1; j < n; j++ {
			if s[j-1] > s[j] {
				temp := s[j-1]
				s[j-1] = s[j]
				s[j] = temp
			}
		}
	}

	return s
}

func computeJolt(numbers []int) int {
	nb1 := 0
	nb3 := 1 // The joltage difference with the computer is always 3

	for i := 1; i < len(numbers); i++ {
		if numbers[i]-numbers[i-1] == 1 {
			nb1++
		} else if numbers[i]-numbers[i-1] == 3 {
			nb3++
		}
	}

	return nb1 * nb3
}

func countCombination(joltage int, numbers []int, computed map[int]int) int {
	n := len(numbers)

	// If already computed
	if _, ok := computed[joltage]; ok {
		return computed[joltage]
	}

	if n == 0 {
		return 1
	}

	comb := countCombination(numbers[0], numbers[1:], computed)
	computed[numbers[0]] = comb

	if n > 1 && numbers[1]-joltage <= 3 {
		c := countCombination(numbers[1], numbers[2:], computed)
		computed[numbers[1]] = c
		comb = comb + c

		if n > 2 && numbers[2]-joltage <= 3 {
			c := countCombination(numbers[2], numbers[3:], computed)
			computed[numbers[2]] = c
			comb = comb + c
		}
	}

	return comb
}
