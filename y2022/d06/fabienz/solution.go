package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	buffer := lines[0]

	startOfPacketMarkerIdx := -1

	for idx := 4; idx < len(buffer); idx++ {
		if isDistinct(buffer[idx-4 : idx]) {
			startOfPacketMarkerIdx = idx
			break
		}
	}

	_, err = fmt.Fprintf(answer, "%d", startOfPacketMarkerIdx)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	buffer := lines[0]

	startOfMessageMarkerIdx := -1

	for idx := 14; idx < len(buffer); idx++ {
		if isDistinct(buffer[idx-14 : idx]) {
			startOfMessageMarkerIdx = idx
			break
		}
	}

	_, err = fmt.Fprintf(answer, "%d", startOfMessageMarkerIdx)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Check if all characters in a string are distinct.
func isDistinct(s string) bool {
	seen := make(map[rune]bool)
	for _, char := range s {
		seen[char] = true
	}
	return len(seen) == len(s)
}
