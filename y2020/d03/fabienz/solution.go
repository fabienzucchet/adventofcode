package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", ride(3, 1, lines))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	res11 := ride(1, 1, lines)
	res31 := ride(3, 1, lines)
	res51 := ride(5, 1, lines)
	res71 := ride(7, 1, lines)
	res12 := ride(1, 2, lines)

	_, err = fmt.Fprintf(answer, "%d", res11*res31*res51*res71*res12)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Position struct {
	X int // Horizontal axis
	Y int // Vertical axis
}

func (p *Position) move(dx int, dy int, mapWidth int) {
	p.X = (p.X + dx) % mapWidth
	p.Y = p.Y + dy
}

func ride(dx int, dy int, lines []string) int {
	mapHeight := len(lines)
	mapWidth := len(lines[0])

	pos := Position{}

	treeCounter := 0

	for pos.Y < mapHeight {
		if lines[pos.Y][pos.X] == '#' {
			treeCounter++
		}
		pos.move(dx, dy, mapWidth)
	}

	return treeCounter
}
