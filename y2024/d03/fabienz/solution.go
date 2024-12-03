package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	program := concatProgramSections(lines)

	instructions, err := mulInstructionsFromProgram(program)
	if err != nil {
		return fmt.Errorf("could not parse mul instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", sumMulResults(instructions))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	program := concatProgramSections(lines)

	enabledSections := enabledSectionsFromProgram(program)

	instructions, err := mulInstructionsFromProgram(concatProgramSections(enabledSections))
	if err != nil {
		return fmt.Errorf("could not parse mul instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", sumMulResults(instructions))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type MulInstruction struct {
	Left  int
	Right int
}

func (m *MulInstruction) Mul() int {
	return m.Left * m.Right
}

// Join all the lines together to ignore newline characters since they don't matter in this case.
func concatProgramSections(lines []string) string {
	return strings.Join(lines, "")
}

func mulInstructionsFromProgram(s string) ([]MulInstruction, error) {
	re := regexp.MustCompile(`mul\(([0-9]{1,3}),([0-9]{1,3})\)`)

	matches := re.FindAllStringSubmatch(s, -1)

	mulInstructions := make([]MulInstruction, 0, len(matches))

	for i, submatch := range matches {
		mulInstruction, err := stringSubMatchToMulInstruction(submatch)
		if err != nil {
			return nil, fmt.Errorf("could not parse mul instruction %d: %w", i, err)
		}

		mulInstructions = append(mulInstructions, mulInstruction)
	}

	return mulInstructions, nil
}

func stringSubMatchToMulInstruction(submatch []string) (MulInstruction, error) {
	if len(submatch) != 3 {
		return MulInstruction{}, fmt.Errorf("could not parse mul instruction: %v", submatch)
	}

	left, err := strconv.Atoi(submatch[1])
	if err != nil {
		return MulInstruction{}, fmt.Errorf("could not parse left operand: %w", err)
	}

	right, err := strconv.Atoi(submatch[2])
	if err != nil {
		return MulInstruction{}, fmt.Errorf("could not parse right operand: %w", err)
	}

	return MulInstruction{
		Left:  left,
		Right: right,
	}, nil
}

func sumMulResults(mulInstructions []MulInstruction) int {
	sum := 0

	for _, mulInstruction := range mulInstructions {
		sum += mulInstruction.Mul()
	}

	return sum
}

func enabledSectionsFromProgram(s string) []string {
	enabledRe := regexp.MustCompile(`(?:^|do\(\))(.*?)(?:don\'t\(\)|$)`)

	return enabledRe.FindAllString(s, -1)
}
