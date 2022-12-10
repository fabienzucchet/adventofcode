package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	depths, err := intsFromStrings(lines)
	if err != nil {
		return fmt.Errorf("error when parsing input: %w", err)
	}

	res := countIncreasingMeasurements(depths, 1)

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	depths, err := intsFromStrings(lines)
	if err != nil {
		return fmt.Errorf("error when parsing input: %w", err)
	}

	res := countIncreasingMeasurements(depths, 3)

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func intsFromStrings(lines []string) (numbers []int, err error) {

	for _, line := range lines {

		number, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, number)
	}

	return numbers, nil
}

func countIncreasingMeasurements(measures []int, window int) int {
	count := 0

	for i := window; i < len(measures); i++ {
		if measures[i-window] < measures[i] {
			count++
		}
	}

	return count
}
