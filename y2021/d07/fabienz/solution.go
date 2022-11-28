package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	positions, err := parseInput(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input %s : %w", lines[0], err)
	}

	_, err = fmt.Fprintf(answer, "%d", findMinimumCostWithLinearCost(positions))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	positions, err := parseInput(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input %s : %w", lines[0], err)
	}

	_, err = fmt.Fprintf(answer, "%d", findMinimumCostWithArithmeticCost(positions))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// INPUT PARSING

// Parse input
func parseInput(input string) (positions []int, err error) {
	for _, pos := range strings.Split(input, ",") {
		position, err := strconv.Atoi(pos)
		if err != nil {
			return nil, fmt.Errorf("error parsing position %s : %w", pos, err)
		}

		positions = append(positions, position)
	}

	return positions, nil
}

// PROCESSING FUNCTIONS

// Move all crabs to a given target position
func moveCrabsToWithLinearCost(target int, positions []int) (cost int) {

	for _, pos := range positions {
		cost += abs(pos - target)
	}

	return cost
}

func moveCrabsToWithArithmeticCost(target int, positions []int) (cost int) {

	for _, pos := range positions {
		n := abs(pos - target)
		cost += n * (n + 1) / 2 // Because 1+2+..+n = n(n+1)/2
	}

	return cost
}

// Find the best solution
func findMinimumCostWithLinearCost(positions []int) (minimalCost int) {

	minimalCost = 1<<63 - 1

	min, max := minmax(positions)

	for target := min; target <= max; target++ {
		cost := moveCrabsToWithLinearCost(target, positions)
		if cost < minimalCost {
			minimalCost = cost
		}
	}

	return minimalCost
}

func findMinimumCostWithArithmeticCost(positions []int) (minimalCost int) {

	minimalCost = 1<<63 - 1

	min, max := minmax(positions)

	for target := min; target <= max; target++ {
		cost := moveCrabsToWithArithmeticCost(target, positions)
		if cost < minimalCost {
			minimalCost = cost
		}
	}

	return minimalCost
}

// HELPER FUNCTIONS

// Returns the absolute value of x
func abs(x int) int {
	if x >= 0 {
		return x
	}

	return -x
}

// Find the minimum and maximum of a list
func minmax(positions []int) (min, max int) {

	min, max = 1<<63-1, -1

	for _, pos := range positions {
		if pos < min {
			min = pos
		}
		if pos > max {
			max = pos
		}
	}

	return min, max
}
