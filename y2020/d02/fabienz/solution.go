package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numberValid := 0
	for _, line := range lines {

		min, max, letter, password, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		if isValidWrongRules(min, max, letter, password) {
			numberValid++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", numberValid)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numberValid := 0
	for _, line := range lines {

		min, max, letter, password, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("error parsing line %s : %w", line, err)
		}

		if isValidCorrectRules(min, max, letter, password) {
			numberValid++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", numberValid)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func isValidWrongRules(min int, max int, letter rune, password string) bool {
	count := 0

	for _, char := range password {
		if char == letter {
			count++
		}
	}

	return count >= min && count <= max
}

func isValidCorrectRules(pos1 int, pos2 int, letter rune, password string) bool {
	return (rune(password[pos1-1]) == letter || rune(password[pos2-1]) == letter) && !(rune(password[pos1-1]) == letter && rune(password[pos2-1]) == letter)
}

func parseLine(line string) (int, int, rune, string, error) {
	re := regexp.MustCompile(`^([0-9]*)-([0-9]*) ([a-z]): ([a-z]*)`)
	match := re.FindStringSubmatch(line)

	if len(match) == 1 {
		return 0, 0, ' ', "", fmt.Errorf("couldn't parse line %s", line)
	}

	min, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, ' ', "", fmt.Errorf("error parsing min value %s : %w", match[1], err)
	}
	max, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, ' ', "", fmt.Errorf("error parsing max value %s : %w", match[2], err)
	}
	letter := rune(match[3][0])
	password := match[4]

	return min, max, letter, password, nil
}
