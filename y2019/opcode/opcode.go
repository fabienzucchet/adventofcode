package opcode

import (
	"fmt"
)

// Structure of the intcode program: instructions, inputs and outputs.
type Intcode struct {
	Program []int
	Inputs  []int
	Outputs []int
	Pos     int
	Halted  bool
}

// Compute a step of the intcode program
func ComputeStep(intcode *Intcode) (err error) {
	// Get the instructions
	instructions := intcode.Program
	pos := intcode.Pos

	// Parse the instruction and the parameters depending on the parameters mode
	opcode := instructions[pos] % 100
	param1Mode := (instructions[pos] / 100) % 10
	param2Mode := (instructions[pos] / 1000) % 10

	switch opcode {
	case 1: // Add
		instructions[instructions[pos+3]] = getParamValue(instructions, pos+1, param1Mode) + getParamValue(instructions, pos+2, param2Mode)
		intcode.Pos = pos + 4
		return nil
	case 2: // Multiply
		instructions[instructions[pos+3]] = getParamValue(instructions, pos+1, param1Mode) * getParamValue(instructions, pos+2, param2Mode)
		intcode.Pos = pos + 4
		return nil
	case 3: // Input
		if len(intcode.Inputs) == 0 {
			// Save the position of the instruction in the program
			intcode.Pos = pos
			intcode.Halted = true
			return nil
		}
		instructions[instructions[pos+1]] = intcode.Inputs[0]
		intcode.Inputs = intcode.Inputs[1:]
		intcode.Pos = pos + 2
		return nil
	case 4: // Output
		intcode.Outputs = append(intcode.Outputs, getParamValue(instructions, pos+1, param1Mode))
		intcode.Pos = pos + 2
		return nil
	case 5: // Jump if true
		if getParamValue(instructions, pos+1, param1Mode) != 0 {
			intcode.Pos = getParamValue(instructions, pos+2, param2Mode)
			return nil
		}
		intcode.Pos = pos + 3
		return nil
	case 6: // Jump if false
		if getParamValue(instructions, pos+1, param1Mode) == 0 {
			intcode.Pos = getParamValue(instructions, pos+2, param2Mode)
			return nil
		}
		intcode.Pos = pos + 3
		return nil
	case 7: // Less than
		if getParamValue(instructions, pos+1, param1Mode) < getParamValue(instructions, pos+2, param2Mode) {
			instructions[instructions[pos+3]] = 1
		} else {
			instructions[instructions[pos+3]] = 0
		}
		intcode.Pos = pos + 4
		return nil
	case 8: // Equals
		if getParamValue(instructions, pos+1, param1Mode) == getParamValue(instructions, pos+2, param2Mode) {
			instructions[instructions[pos+3]] = 1
		} else {
			instructions[instructions[pos+3]] = 0
		}
		intcode.Pos = pos + 4
		return nil
	case 99:
		intcode.Pos = -1
		intcode.Halted = true
		return nil
	default:
		intcode.Pos = -1
		intcode.Halted = true
		return fmt.Errorf("unknown opcode %d", opcode)
	}
}

// Init the intcode program with the given noun and verb
func InitIntcode(intcode Intcode, noun int, verb int) {
	intcode.Program[1] = noun
	intcode.Program[2] = verb
}

// Get a parameter value depending on the parameter mode
func getParamValue(instructions []int, pos int, mode int) int {
	if mode == 0 {
		return instructions[instructions[pos]]
	}
	return instructions[pos]
}

// Run the intcode program
func RunIntcode(intcode *Intcode) ([]int, error) {
	intcode.Halted = false
	for !intcode.Halted {
		err := ComputeStep(intcode)
		if err != nil {
			return intcode.Outputs, err
		}
	}
	return intcode.Outputs, nil
}
