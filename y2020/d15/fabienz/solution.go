package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 15 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var numbers []int

	for _, number := range strings.Split(lines[0], ",") {
		n, err := strconv.Atoi(number)
		if err != nil {
			return fmt.Errorf("error parsing number %s : %w", number, err)
		}

		numbers = append(numbers, n)
	}

	spoken := map[int]int{}

	last := -1

	for i := 1; i <= len(numbers); i++ {
		spoken[last] = i
		last = numbers[i-1]
	}

	for i := len(numbers) + 1; i <= 2020; i++ {
		if t, alreadySpoken := spoken[last]; alreadySpoken {
			spoken[last] = i
			last = i - t
		} else {
			spoken[last] = i
			last = 0
		}
	}

	_, err = fmt.Fprintf(answer, "%d", last)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var numbers []int

	for _, number := range strings.Split(lines[0], ",") {
		n, err := strconv.Atoi(number)
		if err != nil {
			return fmt.Errorf("error parsing number %s : %w", number, err)
		}

		numbers = append(numbers, n)
	}

	spoken := map[int]int{}

	last := -1

	for i := 1; i <= len(numbers); i++ {
		spoken[last] = i
		last = numbers[i-1]
	}

	for i := len(numbers) + 1; i <= 30000000; i++ {
		if t, alreadySpoken := spoken[last]; alreadySpoken {
			spoken[last] = i
			last = i - t
		} else {
			spoken[last] = i
			last = 0
		}
	}

	_, err = fmt.Fprintf(answer, "%d", last)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
