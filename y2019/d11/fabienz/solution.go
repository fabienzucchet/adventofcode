package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
	"github.com/fabienzucchet/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 11 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Create the hull painting robot intcode
	program, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode program: %w", err)
	}

	intcode := opcode.Intcode{Program: program}

	// Create the surface
	surface := Surface{
		Map:             make(map[Coordinate]bool),
		RobotCoordinate: Coordinate{X: 0, Y: 0},
		RobotDirection:  Coordinate{X: 0, Y: 1},
		Intcode:         &intcode,
		SurfacePainted:  make(map[Coordinate]bool),
	}

	// While the intcode is not finished, paint the surface
	for surface.Intcode.Pos >= 0 {
		surface.PaintOneStep()
	}

	_, err = fmt.Fprintf(answer, "%d", len(surface.SurfacePainted))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Create the hull painting robot intcode
	program, err := helpers.IntsFromString(lines[0], ",")
	if err != nil {
		return fmt.Errorf("could not parse intcode program: %w", err)
	}

	intcode := opcode.Intcode{Program: program}

	// Create the surface
	Map := make(map[Coordinate]bool)
	// The starting panel is white
	Map[Coordinate{X: 0, Y: 0}] = true

	surface := Surface{
		Map:             Map,
		RobotCoordinate: Coordinate{X: 0, Y: 0},
		RobotDirection:  Coordinate{X: 0, Y: 1},
		Intcode:         &intcode,
		SurfacePainted:  make(map[Coordinate]bool),
	}

	// While the intcode is not finished, paint the surface
	for surface.Intcode.Pos >= 0 {
		surface.PaintOneStep()
	}

	// Print the surface
	minX := 0
	maxX := 0
	minY := 0
	maxY := 0
	for coord := range surface.SurfacePainted {
		if coord.X < minX {
			minX = coord.X
		}
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.Y < minY {
			minY = coord.Y
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}

	}

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if surface.Map[Coordinate{X: x, Y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	// The answer can be read in the console : PFKHECZU
	_, err = fmt.Fprintf(answer, "%s", "PFKHECZU")
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Coordinate represents a coordinate in a 2D space.
type Coordinate struct {
	X int
	Y int
}

// Represent the ship' surface
type Surface struct {
	// Map of coordinates to color
	Map             map[Coordinate]bool
	RobotCoordinate Coordinate
	RobotDirection  Coordinate
	Intcode         *opcode.Intcode
	SurfacePainted  map[Coordinate]bool
}

// Iterate one step of the painting robot
func (s *Surface) PaintOneStep() {
	// Get the current color of the surface
	var currentColor int
	// if the surface is painted in white, the color is 1
	if s.Map[s.RobotCoordinate] {
		currentColor = 1
	}

	// Send the current color to the intcode
	s.Intcode.Inputs = append(s.Intcode.Inputs, currentColor)

	// Run the intcode
	s.Intcode.RunIntcode()

	// Get the new color
	newColor := s.Intcode.Outputs[0]
	s.Intcode.Outputs = s.Intcode.Outputs[1:]

	// Paint the surface
	if newColor != currentColor {
		s.SurfacePainted[s.RobotCoordinate] = true
		s.Map[s.RobotCoordinate] = newColor == 1
	}

	// Get the new direction
	newDirection := s.Intcode.Outputs[0]
	s.Intcode.Outputs = s.Intcode.Outputs[1:]

	// Update the direction
	if newDirection == 0 {
		// Turn left 90°
		s.RobotDirection = Coordinate{X: -s.RobotDirection.Y, Y: s.RobotDirection.X}
	} else {
		// Turn right 90°
		s.RobotDirection = Coordinate{X: s.RobotDirection.Y, Y: -s.RobotDirection.X}
	}

	// Move the robot
	s.RobotCoordinate = Coordinate{X: s.RobotCoordinate.X + s.RobotDirection.X, Y: s.RobotCoordinate.Y + s.RobotDirection.Y}
}
