package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
	"github.com/fabienzucchet/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 9 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	instructions, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}

	// Create the program
	intcode := opcode.Intcode{Program: instructions, Inputs: []int{1}}

	// Run the program
	intcode.RunIntcode()

	// Fail if the program output is not a single value.
	if len(intcode.Outputs) != 1 {
		return fmt.Errorf("expected 1 output, got %d", len(intcode.Outputs))
	}

	_, err = fmt.Fprintf(answer, "%d", intcode.Outputs[0])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	instructions, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}

	// Create the program
	intcode := opcode.Intcode{Program: instructions, Inputs: []int{2}}

	// Run the program
	intcode.RunIntcode()

	// Fail if the program output is not a single value.
	if len(intcode.Outputs) != 1 {
		return fmt.Errorf("expected 1 output, got %d", len(intcode.Outputs))
	}

	_, err = fmt.Fprintf(answer, "%d", intcode.Outputs[0])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
