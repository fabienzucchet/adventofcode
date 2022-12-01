package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
	"gitlab.com/padok-team/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 5 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Init the program
	program, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}
	inputs := []int{1}
	intcode := opcode.Intcode{Program: program, Inputs: inputs}

	// Run the program
	outputs, err := opcode.RunIntcode(intcode)
	if err != nil {
		return fmt.Errorf("could not run intcode: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", outputs[len(outputs)-1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Init the program
	program, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}

	inputs := []int{5}
	intcode := opcode.Intcode{Program: program, Inputs: inputs}

	// Run the program
	outputs, err := opcode.RunIntcode(intcode)
	if err != nil {
		return fmt.Errorf("could not run intcode: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", outputs[len(outputs)-1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
