package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sub Submarine

	for _, line := range lines {
		direction, value, err := parseInstruction(line)
		if err != nil {
			return fmt.Errorf("could not parse instruction: %w", err)
		}

		sub.moveSubmarineWithoutManual(direction, value)
	}

	_, err = fmt.Fprintf(answer, "%d", sub.hPosition*sub.depth)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var sub Submarine

	for _, line := range lines {
		direction, value, err := parseInstruction(line)
		if err != nil {
			return fmt.Errorf("could not parse instruction: %w", err)
		}

		sub.moveSubmarineWithManual(direction, value)
	}

	_, err = fmt.Fprintf(answer, "%d", sub.hPosition*sub.depth)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Submarine struct {
	depth     int
	hPosition int
	aim       int
}

func parseInstruction(instruction string) (string, int, error) {
	splittedInstruction := strings.SplitN(instruction, " ", 2)
	if len(splittedInstruction) < 2 {
		return "", 0, fmt.Errorf("could not parse instruction %v", instruction)
	}

	val, err := strconv.Atoi(splittedInstruction[1])
	if err != nil {
		return "", 0, err
	}

	return splittedInstruction[0], val, nil
}

func (sub *Submarine) moveSubmarineWithoutManual(direction string, val int) error {
	switch direction {
	case "forward":
		sub.hPosition += val
	case "up":
		sub.depth -= val
	case "down":
		sub.depth += val
	default:
		return fmt.Errorf("invalid direction %v", direction)
	}

	return nil
}

func (sub *Submarine) moveSubmarineWithManual(direction string, val int) error {
	switch direction {
	case "forward":
		sub.hPosition += val
		sub.depth += val * sub.aim
	case "up":
		sub.aim -= val
	case "down":
		sub.aim += val
	default:
		return fmt.Errorf("invalid direction %v", direction)
	}

	return nil
}
