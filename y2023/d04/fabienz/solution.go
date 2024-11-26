package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2023.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	totalScore := 0

	for _, line := range lines {
		splitLine := strings.Split(line, ": ")
		if len(splitLine) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		winning, present, err := splitNumbers(splitLine[1])
		if err != nil {
			return fmt.Errorf("could not split numbers: %w", err)
		}

		totalScore += getScore(winning, present)
	}

	_, err = fmt.Fprintf(answer, "%d", totalScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2023.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	scratchcardsCount := make(map[int]int)
	for i := 0; i < len(lines); i++ {
		scratchcardsCount[i] = 1
	}

	for i, line := range lines {

		splitLine := strings.Split(line, ": ")
		if len(splitLine) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		winning, present, err := splitNumbers(splitLine[1])
		if err != nil {
			return fmt.Errorf("could not split numbers: %w", err)
		}

		matchCount := countMatches(winning, present)
		for j := 0; j < matchCount; j++ {
			scratchcardsCount[i+j+1] += scratchcardsCount[i]
		}
	}

	totalCount := 0
	for _, count := range scratchcardsCount {
		totalCount += count
	}

	_, err = fmt.Fprintf(answer, "%d", totalCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func splitNumbers(s string) (winning []int, present []int, err error) {
	splitLine := strings.Split(s, " | ")
	if len(splitLine) != 2 {
		return nil, nil, fmt.Errorf("invalid line: %s", s)
	}

	winning, err = parseNumbers(strings.Split(splitLine[0], " "))
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse numbers: %w", err)
	}

	present, err = parseNumbers(strings.Split(splitLine[1], " "))
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse numbers: %w", err)
	}

	return winning, present, nil
}

func parseNumbers(numbersToParse []string) (numbers []int, err error) {
	for _, number := range numbersToParse {
		if len(number) == 0 {
			continue
		}

		n, err := strconv.Atoi(number)
		if err != nil {
			return nil, fmt.Errorf("%q is not an integer", number)
		}

		numbers = append(numbers, n)
	}

	return numbers, nil
}

func getScore(winning []int, present []int) (score int) {
	matchCount := 0

	for _, p := range present {
		if isInSlice(p, winning) {
			matchCount++
		}
	}

	return 1 << uint(matchCount-1)
}

func isInSlice(n int, slice []int) bool {
	for _, s := range slice {
		if s == n {
			return true
		}
	}

	return false
}

func countMatches(winning []int, present []int) (matchCount int) {
	for _, p := range present {
		if isInSlice(p, winning) {
			matchCount++
		}
	}

	return matchCount
}
