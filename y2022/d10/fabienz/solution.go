package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	initialRegisters := make(map[string]int)
	initialRegisters["x"] = 1

	// Init the program.
	program, err := initProgram(lines, initialRegisters, false)
	if err != nil {
		return fmt.Errorf("could not init program: %w", err)
	}

	// Run the program.
	program.run()

	// Sum the signal strength for cycles 19, 59, 99, 139, 179 and 219.
	signalStrength := 0
	for i := 20; i <= 220; i += 40 {
		signalStrength += program.signalStrength(i)
	}

	_, err = fmt.Fprintf(answer, "%d", signalStrength)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	initialRegisters := make(map[string]int)
	initialRegisters["x"] = 1

	// Init the program.
	program, err := initProgram(lines, initialRegisters, true)
	if err != nil {
		return fmt.Errorf("could not init program: %w", err)
	}

	// Run the program.
	program.run()

	// The result (PAPKFKEJ) is written in the console.
	_, err = fmt.Fprintf(answer, "%s", "PAPKFKEJ")
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Type representing an instruction.
type instruction struct {
	operation string
	argument  int
}

type program struct {
	instructions []instruction
	registers    map[string]int
	cycle        int
	savedXStates map[int]int // Used to keep track of the x state at each cycle.
	shouldDraw   bool        // Used to draw the signal.
}

// Regex to match an instruction.
var noopRegex = regexp.MustCompile(`^noop$`)

// Regex to add an addx instruction.
var addxRegex = regexp.MustCompile(`^addx (-?\d+)$`)

// ParseInstruction parses an instruction.
func parseInstruction(line string) (instruction, error) {
	switch {
	case noopRegex.MatchString(line):
		return instruction{operation: "noop"}, nil
	case addxRegex.MatchString(line):
		matches := addxRegex.FindStringSubmatch(line)
		argument, err := strconv.Atoi(matches[1])
		if err != nil {
			return instruction{}, fmt.Errorf("could not parse argument: %w", err)
		}
		return instruction{operation: "addx", argument: argument}, nil
	default:
		return instruction{}, fmt.Errorf("could not parse instruction: %q", line)
	}
}

// InitProgram initializes a program based on a list of instructions, an initial value for registers.
func initProgram(lines []string, registers map[string]int, shouldDraw bool) (program, error) {
	instructions := make([]instruction, len(lines))
	for i, line := range lines {
		instruction, err := parseInstruction(line)
		if err != nil {
			return program{}, fmt.Errorf("could not parse instruction: %w", err)
		}
		instructions[i] = instruction
	}

	savedXStates := make(map[int]int)
	savedXStates[0] = registers["x"]

	return program{
		instructions: instructions,
		registers:    registers,
		cycle:        0,
		savedXStates: savedXStates,
		shouldDraw:   shouldDraw,
	}, nil
}

// Apply an instruction to the program.
func (p *program) apply(instruction instruction) {
	switch instruction.operation {
	case "noop":
		if p.shouldDraw {
			p.setPixel(p.cycle%40, p.registers["x"])
		}
		p.cycle++
		p.savedXStates[p.cycle] = p.registers["x"]
	case "addx":
		if p.shouldDraw {
			p.setPixel(p.cycle%40, p.registers["x"])
		}
		if p.shouldDraw {
			p.setPixel((p.cycle+1)%40, p.registers["x"])
		}
		// addx takes two cycles to execute.
		p.savedXStates[p.cycle+1] = p.registers["x"]
		p.savedXStates[p.cycle+2] = p.registers["x"]
		p.cycle += 2
		p.registers["x"] = p.registers["x"] + instruction.argument
	}
}

// Run all the instructions of the program.
func (p *program) run() {
	for _, inst := range p.instructions {
		p.apply(inst)
	}
}

// Compute the signal strength of the program at a given cycle.
func (p *program) signalStrength(cycle int) int {
	return p.savedXStates[cycle] * cycle
}

// Draw a pixel at a given position.
func (p *program) setPixel(x int, spritePos int) {
	pixel := "."
	if helpers.AbsInt(spritePos-x) <= 1 {
		pixel = "#"
	}

	fmt.Printf("%s", pixel)
	if x%40 == 39 {
		fmt.Printf("\n")
	}
}
