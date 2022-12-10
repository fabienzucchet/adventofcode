package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := make([]int, 0, len(lines))

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("could not parse number %s : %w", line, err)
		}
		numbers = append(numbers, number)
	}

	product := 0

	for idx, num1 := range numbers[:len(numbers)-1] {
		for _, num2 := range numbers[idx+1:] {
			if num1+num2 == 2020 {
				product = num1 * num2
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := make([]int, 0, len(lines))

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("could not parse number %s : %w", line, err)
		}
		numbers = append(numbers, number)
	}

	product := 0

	for idx1, num1 := range numbers[:len(numbers)-2] {
		for idx2, num2 := range numbers[idx1+1 : len(numbers)-1] {
			for _, num3 := range numbers[idx2+1:] {
				if num1+num2+num3 == 2020 {
					product = num1 * num2 * num3
				}
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
