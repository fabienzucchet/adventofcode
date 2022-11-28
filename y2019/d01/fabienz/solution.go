package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	masses, err := modulesMassFromLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input ; %w", err)
	}

	sum := 0

	for _, mass := range masses {
		sum += fuelCostFromMass(mass)
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	masses, err := modulesMassFromLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input ; %w", err)
	}

	sum := 0

	for _, mass := range masses {
		sum += fuelCostFromMassWithFuel(mass)
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// INPUT PARSING
func modulesMassFromLines(lines []string) (masses []int, err error) {
	for _, line := range lines {
		mass, err := strconv.Atoi(line)
		if err != nil {
			return masses, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		masses = append(masses, mass)
	}

	return masses, nil
}

// Compute the fuel cost required for a given module
func fuelCostFromMass(mass int) (cost int) {

	return mass/3 - 2
}

// Compute the fuel cost required for a given module and its fuel
func fuelCostFromMassWithFuel(mass int) (totalCost int) {
	cost := fuelCostFromMass(mass)

	totalCost = cost

	for cost > 0 {
		cost = fuelCostFromMass(cost)
		totalCost += cost
	}

	return totalCost
}
