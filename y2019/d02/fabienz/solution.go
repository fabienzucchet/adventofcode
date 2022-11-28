package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
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
	intcode, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode: %w", err)
	}

	// Init the program
	initIntcode(intcode, 12, 2)

	// Run the program
	pos := 0
	for pos >= 0 {
		pos, err = computeStep(intcode, pos)
		if err != nil {
			return fmt.Errorf("could not compute step: %w", err)
		}
	}

	_, err = fmt.Fprintf(answer, "%d", intcode[0])
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
			intcode, err := helpers.IntsFromString(lines[0], ",")
			if err != nil {
				return fmt.Errorf("could not parse intcode: %w", err)
			}

			// Init the program
			initIntcode(intcode, noun, verb)

			// Run the program
			pos := 0
			for pos >= 0 {
				pos, err = computeStep(intcode, pos)
				if err != nil {
					return fmt.Errorf("could not compute step: %w", err)
				}
			}

			output = intcode[0]
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

// Compute a step of the intcode program
func computeStep(intcode []int, pos int) (newPos int, err error) {
	opcode := intcode[pos]
	switch opcode {
	case 1:
		intcode[intcode[pos+3]] = intcode[intcode[pos+1]] + intcode[intcode[pos+2]]
		return pos + 4, nil
	case 2:
		intcode[intcode[pos+3]] = intcode[intcode[pos+1]] * intcode[intcode[pos+2]]
		return pos + 4, nil
	case 99:
		return -1, nil
	default:
		return -1, fmt.Errorf("unknown opcode %d", opcode)
	}
}

// Init the intcode program with the 1202 alarm state
func initIntcode(intcode []int, noun int, verb int) {
	intcode[1] = noun
	intcode[2] = verb
}
