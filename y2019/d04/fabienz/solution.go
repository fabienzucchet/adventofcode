package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// parse the input
	minValue, maxValue := parseInput(lines[0])

	// Count the number of valid passwords
	var count int
	for i := minValue; i <= maxValue; i++ {
		if isValid(intToPassword(i)) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// parse the input
	minValue, maxValue := parseInput(lines[0])

	// Count the number of valid passwords
	var count int
	for i := minValue; i <= maxValue; i++ {
		if isValid2(intToPassword(i)) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Represent a potential password.
type Password [6]int

// Convert an int to a password.
func intToPassword(i int) Password {
	var p Password
	for j := 5; j >= 0; j-- {
		p[j] = i % 10
		i /= 10
	}
	return p
}

// parse the input
func parseInput(line string) (int, int) {
	var min, max int
	fmt.Sscanf(line, "%d-%d", &min, &max)
	return min, max
}

// Check if a password is valid
func isValid(p Password) bool {
	// Check that the password is increasing
	for i := 0; i < 5; i++ {
		if p[i] > p[i+1] {
			return false
		}
	}

	// Check that there is at least one double
	for i := 0; i < 5; i++ {
		if p[i] == p[i+1] {
			return true
		}
	}

	return false
}

// Check if a password is valid for part 2
func isValid2(p Password) bool {
	// Check that the password is increasing
	for i := 0; i < 5; i++ {
		if p[i] > p[i+1] {
			return false
		}
	}

	// Check that there is at least one double
	for i := 0; i < 5; i++ {
		if p[i] == p[i+1] {
			if i == 0 || p[i] != p[i-1] {
				if i == 4 || p[i] != p[i+2] {
					return true
				}
			}
		}
	}

	return false
}
