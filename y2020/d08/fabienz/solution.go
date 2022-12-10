package fabienz

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	operations := make(map[int]string)

	for idx, line := range lines {
		operations[idx] = line
	}

	executed := make(map[int]bool)

	acc, _ := runInstruction(0, 0, operations, executed)

	_, err = fmt.Fprintf(answer, "%d", acc)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	operations := make(map[int]string)

	for idx, line := range lines {
		operations[idx] = line
	}

	var endValue int

	for i := 0; i < len(lines); i++ {
		// If relevant, try to change the instruction

		operation := operations[i]

		executed := make(map[int]bool)

		if operation[:3] == "nop" {
			operations[i] = strings.Replace(operations[i], "nop", "jmp", 1)
			acc, err := runInstruction(0, 0, operations, executed)

			if err == nil {
				endValue = acc
				break
			} else {
				operations[i] = strings.Replace(operations[i], "jmp", "nop", 1)
			}

		} else if operation[:3] == "jmp" {
			operations[i] = strings.Replace(operations[i], "jmp", "nop", 1)
			acc, err := runInstruction(0, 0, operations, executed)

			if err == nil {
				endValue = acc
				break
			} else {
				operations[i] = strings.Replace(operations[i], "nop", "jmp", 1)
			}
		}

	}

	_, err = fmt.Fprintf(answer, "%d", endValue)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func runInstruction(idx int, acc int, operations map[int]string, executed map[int]bool) (int, error) {

	_, ok := executed[idx]

	// If the operation is not already executed
	if !ok {
		executed[idx] = true

		// If the operation doesn't exist : we reached either the end or an error
		if _, exists := operations[idx]; !exists {
			if _, exists := operations[idx-1]; !exists {
				return acc, errors.New("Error in the program")
			}
			return acc, nil
		}

		operation := operations[idx]

		delta, err := strconv.Atoi(operation[5:])
		if err != nil {
			return 0, fmt.Errorf("error parsing %v : %w", operation[5:], err)
		}

		if operation[:3] == "nop" {
			return runInstruction(idx+1, acc, operations, executed)
		} else if operation[:3] == "acc" {
			if operation[4] == '+' {
				return runInstruction(idx+1, acc+delta, operations, executed)
			}
			return runInstruction(idx+1, acc-delta, operations, executed)
		}

		if operation[4] == '+' {
			return runInstruction(idx+delta, acc, operations, executed)
		}
		return runInstruction(idx-delta, acc, operations, executed)
	}

	return acc, fmt.Errorf("infinite loop")
}
