package opcode

import (
	"fmt"
)

// Structure of the intcode program: instructions, inputs and outputs.
type Intcode struct {
	Program      []int
	Inputs       []int
	Outputs      []int
	Pos          int
	Halted       bool
	RelativeBase int
}

// Compute a step of the intcode program
func (intcode *Intcode) ComputeStep() (err error) {
	// Get the instructions
	pos := intcode.Pos

	// Parse the instruction and the parameters depending on the parameters mode
	opcode := intcode.GetValue(pos) % 100
	param1Mode := (intcode.GetValue(pos) / 100) % 10
	param2Mode := (intcode.GetValue(pos) / 1000) % 10
	param3Mode := (intcode.GetValue(pos) / 10000) % 10

	switch opcode {
	case 1: // Add
		intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode))+intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode)))
		intcode.Pos = pos + 4
		return nil
	case 2: // Multiply
		intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode))*intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode)))
		intcode.Pos = pos + 4
		return nil
	case 3: // Input
		if len(intcode.Inputs) == 0 {
			// Save the position of the instruction in the program
			intcode.Pos = pos
			intcode.Halted = true
			return nil
		}
		intcode.SetValue(intcode.getParamAddress(pos+1, param1Mode), intcode.Inputs[0])
		intcode.Inputs = intcode.Inputs[1:]
		intcode.Pos = pos + 2
		return nil
	case 4: // Output
		intcode.Outputs = append(intcode.Outputs, intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode)))
		intcode.Pos = pos + 2
		return nil
	case 5: // Jump if true
		if intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode)) != 0 {
			intcode.Pos = intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode))
			return nil
		}
		intcode.Pos = pos + 3
		return nil
	case 6: // Jump if false
		if intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode)) == 0 {
			intcode.Pos = intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode))
			return nil
		}
		intcode.Pos = pos + 3
		return nil
	case 7: // Less than
		if intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode)) < intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode)) {
			intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), 1)
		} else {
			intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), 0)
		}
		intcode.Pos = pos + 4
		return nil
	case 8: // Equals
		if intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode)) == intcode.GetValue(intcode.getParamAddress(pos+2, param2Mode)) {
			intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), 1)
		} else {
			intcode.SetValue(intcode.getParamAddress(pos+3, param3Mode), 0)
		}
		intcode.Pos = pos + 4
		return nil
	case 9: // Adjust relative base
		intcode.RelativeBase += intcode.GetValue(intcode.getParamAddress(pos+1, param1Mode))
		intcode.Pos = pos + 2
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
func (intcode *Intcode) InitIntcode(noun int, verb int) {
	intcode.Program[1] = noun
	intcode.Program[2] = verb
}

// Get a parameter value depending on the parameter mode
func (intcode Intcode) getParamAddress(pos int, mode int) int {
	switch mode {
	case 0:
		return intcode.GetValue(pos)
	case 2:
		return intcode.RelativeBase + intcode.GetValue(pos)
	default: // default is mode 1
		return pos
	}
}

// Run the intcode program
func (intcode *Intcode) RunIntcode() ([]int, error) {
	intcode.Halted = false
	for !intcode.Halted {
		err := intcode.ComputeStep()
		if err != nil {
			return intcode.Outputs, err
		}
	}
	return intcode.Outputs, nil
}

// Getter of a value in the intcode program at a given position. Adds 0s if the position is out of the program.
func (intcode Intcode) GetValue(pos int) int {
	if pos < len(intcode.Program) {
		return intcode.Program[pos]
	}
	return 0
}

// Setter of a value in the intcode program at a given position. Adds 0s if the position is out of the program.
func (intcode *Intcode) SetValue(pos int, value int) {
	if pos < 0 {
		panic(fmt.Sprintf("trying to set value {%d} at negative position {%d}", value, pos))
	}
	if pos < len(intcode.Program) {
		intcode.Program[pos] = value
		return
	}
	for i := len(intcode.Program); i < pos; i++ {
		intcode.Program = append(intcode.Program, 0)
	}
	intcode.Program = append(intcode.Program, value)
}
