package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
	"gitlab.com/padok-team/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 2 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	if len(lines) != 1 {
		return fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	// Parse the input.
	instructions, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}

	// Create the program
	intcode := opcode.Intcode{Program: instructions}

	// Init the program
	opcode.InitIntcode(intcode, 12, 2)

	// Run the program
	opcode.RunIntcode(&intcode)

	_, err = fmt.Fprintf(answer, "%d", intcode.Program[0])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	output := 0
	result := -1

	// Try all the possible combinations of noun and verb
	for noun := 0; noun < 100 && output != 19690720; noun++ {
		for verb := 0; verb < 100 && output != 19690720; verb++ {
			// Parse the input.
			instructions, err := helpers.IntsFromString(lines[0], ",")
			if err != nil {
				return fmt.Errorf("could not parse intcode: %w", err)
			}

			// Create the program
			intcode := opcode.Intcode{Program: instructions}

			// Init the program
			opcode.InitIntcode(intcode, noun, verb)

			// Run the program
			opcode.RunIntcode(&intcode)

			output = intcode.Program[0]
			if output == 19690720 {
				result = 100*noun + verb
			}
		}
	}

	// TODO: Write your solution to Part 2 below.
	_, err = fmt.Fprintf(answer, "%d", result)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}
