package fabienz

import (
	"fmt"
	"io"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2023.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0

	for _, line := range lines {
		first, last, err := findFirstAndLastDigit(line)
		if err != nil {
			return fmt.Errorf("could not find first and last digit: %w", err)
		}

		sum += 10*first + last
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2023.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0

	for _, line := range lines {
		first, last, err := findFirstAndLastDigitWithWords(line)
		if err != nil {
			return fmt.Errorf("could not find first and last digit: %w", err)
		}

		sum += 10*first + last
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Returns the first and last digit of a string.
func findFirstAndLastDigit(s string) (first int, last int, err error) {
	first, last = -1, -1

	for i := 0; i < len(s); i++ {
		if isDigit(rune(s[i])) {
			if first == -1 {
				first, err = toDigit(rune(s[i]))
				if err != nil {
					return -1, -1, fmt.Errorf("could not convert first digit: %w", err)
				}
			}

			last, err = toDigit(rune(s[i]))
			if err != nil {
				return -1, -1, fmt.Errorf("could not convert last digit: %w", err)
			}
		}
	}

	return first, last, nil
}

// Checks if a rune is a digit.
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Converts a digit rune to an integer.
func toDigit(r rune) (int, error) {
	d := int(r - '0')

	if d < 0 {
		return 0, fmt.Errorf("rune %q is not a digit", r)
	}

	return d, nil
}

// Find the first and last digit also checking words.
func findFirstAndLastDigitWithWords(s string) (first int, last int, err error) {
	first, last = -1, -1

	for i, r := range s {
		if isDigit(r) {
			if first == -1 {
				first, err = toDigit(r)
				if err != nil {
					return -1, -1, fmt.Errorf("could not convert first digit: %w", err)
				}
			}

			last, err = toDigit(r)
			if err != nil {
				return -1, -1, fmt.Errorf("could not convert last digit: %w", err)
			}

			continue
		}

		for word, digit := range wordsToDigits {
			if strings.HasPrefix(s[i:], word) {
				if first == -1 {
					first = digit
				}

				last = digit
			}
		}

	}

	return first, last, nil
}

var wordsToDigits = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}
