package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	grid := gridFromLines(lines)

	count := grid.countWordInGrid("XMAS")

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	grid := gridFromLines(lines)

	count := grid.countXMas()

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Direction int

const (
	Up          Direction = iota
	Down        Direction = iota
	Left        Direction = iota
	Right       Direction = iota
	TopRight    Direction = iota
	TopLeft     Direction = iota
	BottomRight Direction = iota
	BottomLeft  Direction = iota
)

var allDirections = []Direction{Up, Down, Left, Right, TopRight, TopLeft, BottomRight, BottomLeft}

type Character rune

type Grid [][]Character

func gridFromLines(lines []string) Grid {
	grid := make(Grid, len(lines))
	for i, line := range lines {
		grid[i] = make([]Character, len(line))
		for j, char := range line {
			grid[i][j] = Character(char)
		}
	}

	return grid
}

func (g Grid) String() string {
	var str string
	for _, row := range g {
		for _, char := range row {
			str += string(char)
		}
		str += "\n"
	}
	return str
}

func (g Grid) checkIdx(i, j int) bool {
	return i >= 0 && i < len(g) && j >= 0 && j < len((g)[0])
}

func (g Grid) checkWordFromPositionInGridForDirection(word string, i, j int, d Direction) bool {
	for k, char := range word {
		switch d {
		case Up:
			if !g.checkIdx(i-k, j) || g[i-k][j] != Character(char) {
				return false
			}

		case Down:
			if !g.checkIdx(i+k, j) || g[i+k][j] != Character(char) {
				return false
			}

		case Left:
			if !g.checkIdx(i, j-k) || g[i][j-k] != Character(char) {
				return false
			}

		case Right:
			if !g.checkIdx(i, j+k) || g[i][j+k] != Character(char) {
				return false
			}

		case TopRight:
			if !g.checkIdx(i-k, j+k) || g[i-k][j+k] != Character(char) {
				return false
			}

		case TopLeft:
			if !g.checkIdx(i-k, j-k) || g[i-k][j-k] != Character(char) {
				return false
			}

		case BottomRight:
			if !g.checkIdx(i+k, j+k) || g[i+k][j+k] != Character(char) {
				return false
			}

		case BottomLeft:
			if !g.checkIdx(i+k, j-k) || g[i+k][j-k] != Character(char) {
				return false
			}
		}
	}

	return true
}

func (g Grid) countWordFromPositionInGridForDirections(word string, i, j int) int {
	count := 0
	for _, d := range allDirections {
		if g.checkWordFromPositionInGridForDirection(word, i, j, d) {
			count++
		}
	}

	return count
}

func (g Grid) countWordInGrid(word string) int {
	count := 0
	for i, row := range g {
		for j := range row {
			count += g.countWordFromPositionInGridForDirections(word, i, j)
		}
	}

	return count
}

// A position of valid if:
// - its value is an A
// each diagonal of length 3 centered on the position contains a MAS in any direction
//
// Hence the algorithm to check a position is the following:
//   if the value is not an A, return false
//   check if the word MAS is in the grid for the direction BottomRight from the position (i-1, j-1)
//   check if the word MAS is in the grid for the direction TopLeft from the position (i+1, j+1)
//      return false if none of the two checks is true
//   check if the word MAS is in the grid for the direction BottomLeft from the position (i-1, j+1)
//   check if the word MAS is in the grid for the direction TopRight from the position (i+1, j-1)
//   		return false if none of the two checks is true
//   return true

func (g Grid) isXMasCenteredOnPosition(i, j int) bool {
	if g[i][j] != 'A' {
		return false
	}

	if !g.checkWordFromPositionInGridForDirection("MAS", i-1, j-1, BottomRight) && !g.checkWordFromPositionInGridForDirection("MAS", i+1, j+1, TopLeft) {
		return false
	}

	if !g.checkWordFromPositionInGridForDirection("MAS", i-1, j+1, BottomLeft) && !g.checkWordFromPositionInGridForDirection("MAS", i+1, j-1, TopRight) {
		return false
	}

	return true
}

func (g Grid) countXMas() int {
	count := 0
	for i, row := range g {
		for j := range row {
			if g.isXMasCenteredOnPosition(i, j) {
				count++
			}
		}
	}

	return count
}
