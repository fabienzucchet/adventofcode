package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	// For each assignments, count if it overlaps with another assignment.
	for i, line := range lines {
		start1, end1, start2, end2, err := parseAssignment(line)
		if err != nil {
			return fmt.Errorf("could not parse assignment %d: %w", i, err)
		}

		if isAssignmentContained(start1, end1, start2, end2) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	// For each assignments, count if it overlaps with another assignment.
	for i, line := range lines {
		start1, end1, start2, end2, err := parseAssignment(line)
		if err != nil {
			return fmt.Errorf("could not parse assignment %d: %w", i, err)
		}

		if isAssignmentOverlapping(start1, end1, start2, end2) {
			count++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Parse a pair of assignments.
func parseAssignment(line string) (int, int, int, int, error) {
	var start1, end1, start2, end2 int
	_, err := fmt.Sscanf(line, "%d-%d,%d-%d", &start1, &end1, &start2, &end2)
	return start1, end1, start2, end2, err
}

// Check if one of the two assignments is contained in the other.
func isAssignmentContained(AStart, AEnd, BStart, BEnd int) bool {
	return AStart >= BStart && AEnd <= BEnd || BStart >= AStart && BEnd <= AEnd
}

// Check if two assignments overlap.
func isAssignmentOverlapping(AStart, AEnd, BStart, BEnd int) bool {
	return AStart <= BStart && BStart <= AEnd || BStart <= AStart && AStart <= BEnd
}
