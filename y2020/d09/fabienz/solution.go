package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var numbers []int

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		numbers = append(numbers, number)
	}

	invalidNumber := findInvalidNumber(25, numbers)

	_, err = fmt.Fprintf(answer, "%d", invalidNumber)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var numbers []int

	for _, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		numbers = append(numbers, number)
	}

	invalidNumber := findInvalidNumber(25, numbers)

	_, err = fmt.Fprintf(answer, "%d", findWeakness(invalidNumber, numbers))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func checkIsValid(number int, previous []int) bool {
	n := len(previous)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if previous[i]+previous[j] == number {
				return true
			}
		}
	}

	return false
}

func findInvalidNumber(preamble_size int, numbers []int) int {
	for i := preamble_size; i < len(numbers); i++ {
		if !checkIsValid(numbers[i], numbers[i-preamble_size:i]) {
			return numbers[i]
		}
	}

	return -1
}

func sumNumbers(i int, j int, numbers []int) int {
	sum := 0

	for k := 0; k <= j; k++ {
		sum = sum + numbers[i+k]
	}

	return sum
}

func computeWeakness(numbers []int) int {
	min := numbers[0]
	max := numbers[0]

	for i := 1; i < len(numbers); i++ {
		if numbers[i] > max {
			max = numbers[i]
		} else if numbers[i] < min {
			min = numbers[i]
		}
	}

	return min + max
}

func findWeakness(target int, numbers []int) int {
	for i := 0; i < len(numbers); i++ {
		for j := 0; sumNumbers(i, j, numbers) <= target; j++ {
			if sumNumbers(i, j, numbers) == target {
				return computeWeakness(numbers[i : i+j])
			}
		}
	}

	return 0
}
