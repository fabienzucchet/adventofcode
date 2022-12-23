package fabienz

import (
	"fmt"
	"io"
	"math"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	elves := elvesFromLines(lines)

	// Simulate the rounds.
	for i := 0; i < 10; i++ {
		elves, _ = simulate(elves, i)
	}

	_, err = fmt.Fprintf(answer, "%d", countEmptyTiles(elves))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	elves := elvesFromLines(lines)

	// Simulate the rounds while elves are still moving.
	i := 0
	for moved := true; moved; {
		elves, moved = simulate(elves, i)
		i++
	}
	_, err = fmt.Fprintf(answer, "%d", i)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// The directions to consider when moving an elf.
var directions = []string{"N", "S", "W", "E"}

// Key is the current position of an elf, value is the proposed position.
type elves map[helpers.Coord2D]helpers.Coord2D

// Print the elves.
func (e elves) String() string {
	minX, maxX, minY, maxY := math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32

	for pos := range e {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	var s string

	s += fmt.Sprintf("minX: %d, maxX: %d, minY: %d, maxY: %d\n", minX, maxX, minY, maxY)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if _, ok := e[helpers.Coord2D{X: x, Y: y}]; ok {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}

	return s
}

// Parse the input.
func elvesFromLines(lines []string) elves {
	e := make(elves)

	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				e[helpers.Coord2D{X: x, Y: y}] = helpers.Coord2D{}
			}
		}
	}

	return e
}

// Simulate one round.
func simulate(e elves, roundIdx int) (elves, bool) {
	targetMap := make(map[helpers.Coord2D]int)
	moved := false

	// Phase 1: all elves propose a new position.
	for pos := range e {
		// If nobody is around, they stay here.
		shouldMove := false
		for _, n := range pos.Neighbors() {
			if _, ok := e[n]; ok {
				shouldMove = true
				break
			}
		}

		if !shouldMove {
			e[pos] = pos
			targetMap[pos]++
			continue
		}

		moved = true

		// Else, they examine the proposition in the correct order to propose a new position.
		for i := roundIdx; i < roundIdx+len(directions); i++ {
			dir := directions[i%len(directions)]

			canPropose := true

			switch dir {
			case "N":
				for offset := -1; offset <= 1; offset++ {
					if _, ok := e[helpers.Coord2D{X: pos.X + offset, Y: pos.Y - 1}]; ok {
						canPropose = false
						break
					}
				}
			case "S":
				for offset := -1; offset <= 1; offset++ {
					if _, ok := e[helpers.Coord2D{X: pos.X + offset, Y: pos.Y + 1}]; ok {
						canPropose = false
						break
					}
				}
			case "W":
				for offset := -1; offset <= 1; offset++ {
					if _, ok := e[helpers.Coord2D{X: pos.X - 1, Y: pos.Y + offset}]; ok {
						canPropose = false
						break
					}
				}
			case "E":
				for offset := -1; offset <= 1; offset++ {
					if _, ok := e[helpers.Coord2D{X: pos.X + 1, Y: pos.Y + offset}]; ok {
						canPropose = false
						break
					}
				}
			}

			if canPropose {
				switch dir {
				case "N":
					e[pos] = helpers.Coord2D{X: pos.X, Y: pos.Y - 1}
				case "S":
					e[pos] = helpers.Coord2D{X: pos.X, Y: pos.Y + 1}
				case "W":
					e[pos] = helpers.Coord2D{X: pos.X - 1, Y: pos.Y}
				case "E":
					e[pos] = helpers.Coord2D{X: pos.X + 1, Y: pos.Y}
				}
				targetMap[e[pos]]++
				break
			}
		}
	}

	// Phase 2: all elves move to their new position if they are the only one to go there.
	newElves := make(elves)
	for pos, newPos := range e {
		// If more than an elf wants to go there, they stay here.
		if targetMap[newPos] == 1 {
			newElves[newPos] = newPos
			continue
		}

		newElves[pos] = pos
	}

	return newElves, moved
}

// Count the number of empty tiles in the smallest rectangle that contains all elves.
func countEmptyTiles(e elves) int {
	minX, minY, maxX, maxY := math.MaxInt, math.MaxInt, math.MinInt, math.MinInt

	for pos := range e {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	return (maxX-minX+1)*(maxY-minY+1) - len(e)
}
