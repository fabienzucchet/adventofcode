package fabienz

import (
	"fmt"
	"io"
	"regexp"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	ship, instructions, err := ParseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Execute the instructions.
	for _, instruction := range instructions {
		err := ship.Execute(instruction)
		if err != nil {
			return fmt.Errorf("could not execute instruction: %w", err)
		}
	}

	result, err := ship.String()
	if err != nil {
		return fmt.Errorf("could not format result: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%s", result)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	ship, instructions, err := ParseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Execute the instructions.
	for _, instruction := range instructions {
		err := ship.ExecuteKeepOrder(instruction)
		if err != nil {
			return fmt.Errorf("could not execute instruction: %w", err)
		}
	}

	result, err := ship.String()
	if err != nil {
		return fmt.Errorf("could not format result: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%s", result)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Define a pile.

type Pile []string

// Push to the pile.
func (p *Pile) Push(s string) {
	*p = append(*p, s)
}

// Pop from the pile.
func (p *Pile) Pop() string {
	if len(*p) == 0 {
		return ""
	}
	s := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return s
}

// Peek at the top of the pile.
func (p *Pile) Peek() string {
	if len(*p) == 0 {
		return ""
	}
	return (*p)[len(*p)-1]
}

// IsEmpty returns true if the pile is empty.
func (p *Pile) IsEmpty() bool {
	return len(*p) == 0
}

// Flip a pile.
func (p *Pile) Flip() {
	for i := 0; i < len(*p)/2; i++ {
		(*p)[i], (*p)[len(*p)-1-i] = (*p)[len(*p)-1-i], (*p)[i]
	}
}

// Define a type for a ship.
type Ship map[int]*Pile

// Define a type to contain an instruction.
type Instruction struct {
	source int
	target int
	count  int
}

// Parse the input.
func ParseLines(lines []string) (ship Ship, instructions []Instruction, err error) {
	ship = make(Ship)

	for i := 1; i < 10; i++ {
		ship[i] = &Pile{}
	}

	// Regex used to parse the initial state of the ship.
	re := regexp.MustCompile(`^(\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3}) (\[[A-Z]\]|\s{3})$`)

	for i, line := range lines {
		// The first lines contains the initial state of the ship.
		if i < 8 {
			matches := re.FindStringSubmatch(line)
			if len(matches) != 10 {
				return nil, nil, fmt.Errorf("could not parse line %d", i)
			}
			for j := 1; j < 10; j++ {
				if matches[j] != "   " {
					ship[j].Push(matches[j])
				}
			}
			continue
		}

		// The remaining lines contain the instructions.
		if i >= 10 {
			var instruction Instruction
			_, err := fmt.Sscanf(line, "move %d from %d to %d", &instruction.count, &instruction.source, &instruction.target)
			if err != nil {
				return nil, nil, fmt.Errorf("could not parse line %d", i)
			}
			instructions = append(instructions, instruction)
		}
	}

	// We need to flip the piles because the input is in reverse order.
	for i := 1; i < 10; i++ {
		ship[i].Flip()
	}

	return
}

// Execute an instruction.
func (ship Ship) Execute(instruction Instruction) (err error) {
	for i := 0; i < instruction.count; i++ {
		if ship[instruction.source].IsEmpty() {
			return fmt.Errorf("could not run instructionexecution %d x %d -> %d : source pile %d is empty", instruction.count, instruction.source, instruction.target, instruction.source)
		}
		ship[instruction.target].Push(ship[instruction.source].Pop())
	}

	return nil
}

// Execute an instruction and keep the order.
func (ship Ship) ExecuteKeepOrder(instruction Instruction) (err error) {

	// Store the crates to move in a list to keep order.
	crates := make([]string, instruction.count)

	// Fetch the crates to move from the source pile.
	for i := 0; i < instruction.count; i++ {
		if ship[instruction.source].IsEmpty() {
			return fmt.Errorf("could not run instructionexecution %d x %d -> %d : source pile %d is empty", instruction.count, instruction.source, instruction.target, instruction.source)
		}
		crates[i] = ship[instruction.source].Pop()
	}

	// Move the crates to the target pile.
	for i := instruction.count - 1; i >= 0; i-- {
		ship[instruction.target].Push(crates[i])
	}

	return nil
}

// Format the result.
func (ship Ship) String() (result string, err error) {
	for i := 1; i < 10; i++ {
		if ship[i].IsEmpty() {
			return "", fmt.Errorf("pile %d is empty", i)
		} else {
			result += ship[i].Peek()[1:2]
		}
	}

	return
}

// Display the contents of the ship.
func (ship Ship) Display() {
	for i := 1; i < 10; i++ {
		helpers.Println("pile", i, ship[i])
	}
}
