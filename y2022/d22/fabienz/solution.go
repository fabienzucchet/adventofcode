package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// const FACE_SIZE = 4

const FACE_SIZE = 50

// PartOne solves the first problem of day 22 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	b := boardFromLines(lines)

	// Execute the instructions.
	err = b.executeAll("PACMANWRAP")
	if err != nil {
		return fmt.Errorf("could not execute instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", b.password())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	b := boardFromLines(lines)

	// Execute the instructions.
	err = b.executeAll("CUBEWRAP")
	if err != nil {
		return fmt.Errorf("could not execute instructions: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", b.password())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type board struct {
	grid         map[helpers.Coord2D]string
	wrapMap      map[wrapCombo]wrapCombo
	instructions []string
	pos          helpers.Coord2D
	dir          helpers.Coord2D
}

func boardFromLines(lines []string) board {
	grid := make(map[helpers.Coord2D]string)
	for y, line := range lines {
		// The last line is the instructions.
		if y == len(lines)-1 {
			continue
		}

		for x, r := range line {
			if r == ' ' {
				continue
			}
			grid[helpers.Coord2D{X: x + 1, Y: y + 1}] = string(r)
		}
	}

	// Parse the instructions.
	var acc string
	var instructions []string
	for _, c := range lines[len(lines)-1] {
		if c == 'L' || c == 'R' {
			instructions = append(instructions, acc)
			acc = ""

			instructions = append(instructions, string(c))
			continue
		}

		acc += string(c)
	}
	instructions = append(instructions, acc)

	// Find the starting position.
	pos := helpers.Coord2D{X: 1, Y: 1}
	for {
		if grid[pos] == "." {
			break
		}
		pos.X++
	}

	// Create the wrapMap.
	wrapMap := buildWrapMap(FACE_SIZE)

	return board{
		grid:         grid,
		instructions: instructions,
		pos:          pos,
		dir:          helpers.Coord2D{X: 1, Y: 0},
		wrapMap:      wrapMap,
	}
}

// nextPos returns the next position in the given direction, handling wrapping.
func (b *board) nextPos() helpers.Coord2D {
	nextPos := b.pos.Add(b.dir)

	// Handle wrapping.
	if _, ok := b.grid[nextPos]; !ok {
		switch {
		case b.dir.X == 1:
			nextPos.X = 1
		case b.dir.X == -1:
			nextPos.X = len(b.grid)
		case b.dir.Y == 1:
			nextPos.Y = 1
		case b.dir.Y == -1:
			nextPos.Y = len(b.grid)
		}
	}

	for _, ok := b.grid[nextPos]; !ok; _, ok = b.grid[nextPos] {
		switch {
		case b.dir.X == 1:
			nextPos.X++
		case b.dir.X == -1:
			nextPos.X--
		case b.dir.Y == 1:
			nextPos.Y++
		case b.dir.Y == -1:
			nextPos.Y--
		}
	}

	return nextPos
}

type wrapCombo struct {
	pos helpers.Coord2D
	dir helpers.Coord2D
}

// Build the wrapMap.
func buildWrapMap(size int) map[wrapCombo]wrapCombo {
	m := make(map[wrapCombo]wrapCombo)

	switch size {
	case 4:
		// Top of the face 1 facing up ends top of the face 2 facing down.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: size - i, Y: size + 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Right of the face 1 facing right ends right of the face 6 facing left.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 4 * size, Y: 3*size - i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Left of the face 1 facing left ends top of the face 3 facing down.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1, Y: 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: size + 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Top of the face 2 facing up ends top of the face 1 facing down.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: size + 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: 3*size - i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Left of the face 2 facing left ends bottom of the face 6 facing up.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1, Y: size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 4*size - i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Bottom of the face 2 facing down ends bottom of the face 5 facing up.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: 2 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: 3*size - i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Top of the face 3 facing up ends left of the face 1 facing right.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: size + 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: 2*size + 1, Y: 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Bottom of the face 3 facing down ends left of the face 5 facing right.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: 2 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: 2*size + 1, Y: 3*size - i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Right of the face 4 facing right ends top of the face 6 facing down.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 4*size - i, Y: 2*size + 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Left of the face 5 facing left ends bottom of the face 3 facing up.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1, Y: 2*size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 2*size - i, Y: 2 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Bottom of the face 5 facing down ends bottom of face 2 facing up.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: size - i, Y: 2 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Top of the face 6 facing up ends right of the face 4 facing left.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 3*size + 1 + i, Y: 2*size + 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: 2*size - i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Right of the face 6 facing right ends right of the face 1 facing left.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 4 * size, Y: 2*size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: size - i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Bottom of the face 6 facing down ends left of the face 2 facing right.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 3*size + 1 + i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: 1, Y: 2*size - i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}
	case 50:
		// Top of the face 1 facing up ends left of the face 6 facing right with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: 1, Y: 3*size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Left of the face 1 facing left ends left of the face 4 facing right with idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1, Y: 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 1, Y: 3*size - i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Top of the face 2 facing up ends bottom of the face 6 facing up with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: 4 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Right of the face 2 facing right ends right of the face 5 facing left with idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 2 * size, Y: 3*size - i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Bottom of the face 2 facing down ends right of the face 3 facing left with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: 2 * size, Y: size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Left of the face 3 facing left ends top of the face 4 facing down with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1, Y: size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: 2*size + 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Right of the face 3 facing right ends bottom of the face 2 facing up with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2 * size, Y: size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Top of the face 4 facing up ends left of the face 3 facing right with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: 2*size + 1}, dir: helpers.Coord2D{X: 0, Y: -1}}] = wrapCombo{pos: helpers.Coord2D{X: size + 1, Y: size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Left of the face 4 facing left ends left of face 1 facing right with idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1, Y: 2*size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: size + 1, Y: size - i}, dir: helpers.Coord2D{X: 1, Y: 0}}
		}

		// Right of the face 5 facing right ends right of the face 2 facing left with idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 2 * size, Y: 2*size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: 3 * size, Y: size - i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Bottom of the face 5 facing down ends right of the face 6 facing left with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: size, Y: 3*size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}
		}

		// Left of the face 6 facing left ends top of the face 1 facing down with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1, Y: 3*size + 1 + i}, dir: helpers.Coord2D{X: -1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}

		// Right of the face 6 facing right ends bottom of the face 5 facing up with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: size, Y: 3*size + 1 + i}, dir: helpers.Coord2D{X: 1, Y: 0}}] = wrapCombo{pos: helpers.Coord2D{X: size + 1 + i, Y: 3 * size}, dir: helpers.Coord2D{X: 0, Y: -1}}
		}

		// Bottom of the face 6 facing down ends top of the face 2 facing down with no idx change.
		for i := 0; i < size; i++ {
			m[wrapCombo{pos: helpers.Coord2D{X: 1 + i, Y: 4 * size}, dir: helpers.Coord2D{X: 0, Y: 1}}] = wrapCombo{pos: helpers.Coord2D{X: 2*size + 1 + i, Y: 1}, dir: helpers.Coord2D{X: 0, Y: 1}}
		}
	default:
		panic("invalid size")
	}

	return m
}

// The map is now a cube, the new nextPos function handles the wrapping around the cube for part 2.
func (b *board) nextPos2() (helpers.Coord2D, helpers.Coord2D) {
	nextPos := b.pos.Add(b.dir)

	// Check if we reached the end of the face.
	if _, ok := b.grid[nextPos]; !ok {
		wc := b.wrapMap[wrapCombo{pos: b.pos, dir: b.dir}]
		if b.grid[wc.pos] == "." {
			return wc.pos, wc.dir
		}
	}

	if b.grid[nextPos] == "." {
		return nextPos, b.dir
	}

	return b.pos, b.dir
}

// Execute an instruction.
func (b *board) execute(instruction string, mode string) error {
	switch instruction {
	case "R":
		b.dir = helpers.Coord2D{X: -b.dir.Y, Y: b.dir.X}
	case "L":
		b.dir = helpers.Coord2D{X: b.dir.Y, Y: -b.dir.X}
	default:
		dist, err := strconv.Atoi(instruction)
		if err != nil {
			return fmt.Errorf("could not parse instruction: %w", err)
		}
		var nextPos helpers.Coord2D
		for i := 0; i < dist; i++ {
			if mode == "CUBEWRAP" {
				nextPos, b.dir = b.nextPos2()
			} else {
				nextPos = b.nextPos()
			}
			if b.grid[nextPos] == "#" {
				break
			}
			b.pos = nextPos
		}
	}
	return nil
}

// Execute all instructions.
func (b *board) executeAll(mode string) error {
	for _, instruction := range b.instructions {
		if err := b.execute(instruction, mode); err != nil {
			return fmt.Errorf("could not execute instruction: %w", err)
		}
	}
	return nil
}

// Compute the password.
func (b *board) password() int {
	var facing int
	switch b.dir {
	case helpers.Coord2D{X: 1, Y: 0}:
		facing = 0
	case helpers.Coord2D{X: 0, Y: 1}:
		facing = 1
	case helpers.Coord2D{X: -1, Y: 0}:
		facing = 2
	case helpers.Coord2D{X: 0, Y: -1}:
		facing = 3
	}
	return 1000*b.pos.Y + 4*b.pos.X + facing
}
