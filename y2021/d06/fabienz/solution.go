package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ages, err := parseInput(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input data : %w", err)
	}

	for i := 0; i < 80; i++ {
		ages = iterate(ages)
	}

	_, err = fmt.Fprintf(answer, "%d", count(ages))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	ages, err := parseInput(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input data : %w", err)
	}

	for i := 0; i < 256; i++ {
		ages = iterate(ages)
	}

	_, err = fmt.Fprintf(answer, "%d", count(ages))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

// Count the number of fishes with a given age
type Ages map[int]int

// PARSE INPUT

// Parse the input to create a slice of ages
func parseInput(line string) (ages Ages, err error) {

	ages = make(map[int]int)

	for _, age := range strings.Split(line, ",") {
		ageInt, err := strconv.Atoi(age)
		if err != nil {
			return nil, fmt.Errorf("error when parsing age of fish %s : %w", age, err)
		}

		ages[ageInt]++
	}

	return ages, nil
}

// Iterate 1 day
func iterate(ages Ages) (newAges Ages) {

	newAges = make(map[int]int)

	for age, count := range ages {
		switch age {
		case 0:
			newAges[6] += count
			// Create a new fish for each fish
			newAges[8] += count
		default:
			newAges[age-1] += count
		}
	}

	return newAges
}

// Count fishes
func count(ages Ages) (population int) {
	for _, count := range ages {
		population += count
	}

	return population
}
