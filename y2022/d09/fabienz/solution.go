package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse input
	instructions := parseLines(lines)

	// Init a rope
	rope := newRope(2)

	// Apply instructions
	for _, instruction := range instructions {
		rope.applyInstruction(instruction)
	}

	// Count the number of visited coordinates
	visitedCount := 0
	for _, v := range rope.visited {
		if v {
			visitedCount++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", visitedCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse input
	instructions := parseLines(lines)

	// Init a rope
	rope := newRope(10)

	// Apply instructions
	for _, instruction := range instructions {
		rope.applyInstruction(instruction)
	}

	// Count the number of visited coordinates
	visitedCount := 0
	for _, v := range rope.visited {
		if v {
			visitedCount++
		}
	}
	_, err = fmt.Fprintf(answer, "%d", visitedCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Cartesian coordinates of a point in a 2D plane.
type Point struct {
	X, Y int
}

// Store the visited coordinates in the 2D plane in a map.
type visited map[Point]bool

// State of the rope.
type Rope struct {
	knots   []Point
	visited visited
}

// Represent an instruction
type instruction struct {
	direction string
	distance  int
}

// Parse the input and return a slice of instructions.
func parseLines(lines []string) []instruction {
	var instructions []instruction
	for _, line := range lines {
		var direction string
		var distance int
		fmt.Sscanf(line, "%s %d", &direction, &distance)
		instructions = append(instructions, instruction{direction, distance})
	}
	return instructions
}

// Init the rope
func newRope(length int) Rope {
	rope := Rope{
		knots:   make([]Point, length),
		visited: make(visited),
	}
	rope.visited[rope.knots[len(rope.knots)-1]] = true
	return rope
}

// Move the rope one step in the given direction.
func (r *Rope) move(direction string) {
	switch direction {
	case "U":
		r.knots[0].Y++
	case "D":
		r.knots[0].Y--
	case "L":
		r.knots[0].X--
	case "R":
		r.knots[0].X++
	}

	// Update the position of every knot that is not the head.
	for i := 1; i < len(r.knots); i++ {
		r.moveKnot(i)
	}

	// Mark the visited coordinates.
	r.visited[r.knots[len(r.knots)-1]] = true
}

// Move a knot of the rope.
func (r *Rope) moveKnot(index int) {
	// Update the position of the knots[index].
	// If the knots[index-1] and the knots[index] are touching i.e. of their coordinates are different of at most 1, do not move the knots[index].
	if helpers.AbsInt(r.knots[index-1].X-r.knots[index].X) > 1 || helpers.AbsInt(r.knots[index-1].Y-r.knots[index].Y) > 1 {
		// Move one step in the direction of the knots[index-1].
		switch {
		case r.knots[index-1].X > r.knots[index].X:
			r.knots[index].X++
		case r.knots[index-1].X < r.knots[index].X:
			r.knots[index].X--
		default:
			break
		}

		// Move one step in the direction of the knots[index-1].
		switch {
		case r.knots[index-1].Y > r.knots[index].Y:
			r.knots[index].Y++
		case r.knots[index-1].Y < r.knots[index].Y:
			r.knots[index].Y--
		default:
			break
		}
	}
}

// Apply an instruction to the rope.
func (r *Rope) applyInstruction(instruction instruction) {
	for i := 0; i < instruction.distance; i++ {
		r.move(instruction.direction)
	}
}
