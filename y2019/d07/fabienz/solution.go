package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
	"gitlab.com/padok-team/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 7 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Compute the thrust for each possible phase setting. Store the result in a slice.
	var thrusts []int
	for _, sequence := range permutations([]int{0, 1, 2, 3, 4}) {
		// Temporary variable to store the output of the previous amplifier.
		previousOutput := 0

		// Process the 5 amplifiers.
		for i := 0; i < 5; i++ {
			// Parse the instructions.
			instructions, err := helpers.IntsFromString(lines[0], ",")
			if err != nil {
				return fmt.Errorf("could not parse instructions: %w", err)
			}

			inputs := []int{sequence[i], previousOutput}

			// Init the computer.
			intcode := opcode.Intcode{Program: instructions, Inputs: inputs}

			// Run the computer.
			outputs, err := opcode.RunIntcode(&intcode)
			if err != nil {
				return fmt.Errorf("could not run intcode: %w", err)
			}

			// Store the output of the current amplifier.
			previousOutput = outputs[0]

			// Store the output of the last amplifier.
			if i == 4 {
				thrusts = append(thrusts, outputs[0])
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", helpers.MaxInts(thrusts))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Compute the thrust for each possible phase setting. Store the result in a slice.
	var thrusts []int
	for _, sequence := range permutations([]int{5, 6, 7, 8, 9}) {
		// Init the 5 amplifiers.
		var amplifiers []opcode.Intcode
		for i := 0; i < 5; i++ {
			// Parse the instructions.
			instructions, err := helpers.IntsFromString(lines[0], ",")
			if err != nil {
				return fmt.Errorf("could not parse instructions: %w", err)
			}

			// Init the computer.
			intcode := opcode.Intcode{Program: instructions, Inputs: []int{sequence[i]}}
			amplifiers = append(amplifiers, intcode)
		}

		// Store the output of the previous amplifier.
		var previousOutput int

		// Process the feedback loop while the last amplifier is not halted with a 99 code.
		for amplifiers[4].Pos != -1 {

			// Process the 5 amplifiers.
			for i := 0; i < 5; i++ {
				// Add the output of the previous amplifier to the inputs of the current amplifier.
				amplifiers[i].Inputs = append(amplifiers[i].Inputs, previousOutput)

				// Run the computer.
				outputs, err := opcode.RunIntcode(&amplifiers[i])
				if err != nil {
					return fmt.Errorf("could not run intcode: %w", err)
				}

				// Store the output of the current amplifier.
				previousOutput = outputs[len(outputs)-1]
			}
		}

		// Store the output of the last amplifier.
		thrusts = append(thrusts, previousOutput)
	}

	_, err = fmt.Fprintf(answer, "%d", helpers.MaxInts(thrusts))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Generate all permutations of distinct integers between 0 and 4.
func permutations(ints []int) [][]int {
	var result [][]int

	// Base case.x
	if len(ints) == 1 {
		return [][]int{ints}
	}

	// Recursive case.
	for i, v := range ints {
		// Copy the slice.
		rest := make([]int, len(ints))
		copy(rest, ints)

		// Remove the current element.
		rest = append(rest[:i], rest[i+1:]...)

		// Generate all permutations of the remaining elements.
		permutations := permutations(rest)

		// Append the current element to each permutation.
		for _, p := range permutations {
			result = append(result, append(p, v))
		}
	}

	return result
}
