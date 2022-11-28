package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

const DIAGRAMSIZE = 1000 // Assume that no line has coordinates greater than 999

// PartOne solves the first problem of day 5 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	vents, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing vents from input : %w", err)
	}

	var d diagram

	for _, vent := range vents {
		d.drawVentWithoutDiagonals(vent)
	}

	_, err = fmt.Fprintf(answer, "%d", d.countOverlaps())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	vents, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing vents from input : %w", err)
	}

	var d diagram

	for _, vent := range vents {
		d.drawVentWithDiagonals(vent)
	}

	_, err = fmt.Fprintf(answer, "%d", d.countOverlaps())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

// The types used in this puzzle to structure the data are defined below
type Coordinates struct {
	x int
	y int
}

type Vent struct {
	start Coordinates
	end   Coordinates
}

// Type to represent the diagram
type diagram [DIAGRAMSIZE][DIAGRAMSIZE]int

// INPUT PARSING

// We use a regex to parse the lines
var re = regexp.MustCompile(`^([0-9]+),([0-9]+) -> ([0-9]+),([0-9]+)$`)

// Parse the input for lines
func parseLines(lines []string) (vents []Vent, err error) {

	for _, line := range lines {
		match := re.FindStringSubmatch(line)

		if len(match) != 5 {
			return nil, fmt.Errorf("error parsing the line %s", line)
		}

		x1, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing the coordinate %s : %w", match[1], err)
		}
		y1, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing the coordinate %s : %w", match[2], err)
		}
		x2, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, fmt.Errorf("error parsing the coordinate %s : %w", match[3], err)
		}
		y2, err := strconv.Atoi(match[4])
		if err != nil {
			return nil, fmt.Errorf("error parsing the coordinate %s : %w", match[4], err)
		}

		vents = append(vents, Vent{
			start: Coordinates{x: x1, y: y1},
			end:   Coordinates{x: x2, y: y2},
		})
	}

	return vents, nil
}

// METHODS ON THE DIAGRAM

// Draw the vent if it's a vertical or horizontal line (Part 1)
func (d *diagram) drawVentWithoutDiagonals(v Vent) {
	switch {
	case v.start.x == v.end.x:
		d.drawVerticalLine(v.start.x, v.start.y, v.end.y)

	case v.start.y == v.end.y:
		d.drawHorizontalLine(v.start.x, v.end.x, v.start.y)
	}
}

// Draw the vent if it's a vertical, horizontal or diagonal line (Part 2)
func (d *diagram) drawVentWithDiagonals(v Vent) {
	switch {
	case v.start.x == v.end.x:
		d.drawVerticalLine(v.start.x, v.start.y, v.end.y)

	case v.start.y == v.end.y:
		d.drawHorizontalLine(v.start.x, v.end.x, v.start.y)

	case abs(v.start.x-v.end.x) == abs(v.start.y-v.end.y):
		d.drawDiagonalLine(v.start.x, v.end.x, v.start.y, v.end.y)
	}
}

// Draw a vertical line
func (d *diagram) drawVerticalLine(x, y1, y2 int) {

	min, max := minmax(y1, y2)

	for i := min; i <= max; i++ {
		d[i][x]++
	}

}

// Draw an horizontal line
func (d *diagram) drawHorizontalLine(x1, x2, y int) {

	min, max := minmax(x1, x2)

	for i := min; i <= max; i++ {
		d[y][i]++
	}

}

func (d *diagram) drawDiagonalLine(x1, x2, y1, y2 int) {
	deltai := abs(x2 - x1)

	for i := 0; i <= deltai; i++ {
		d[y1+sgn(y2-y1)*i][x1+sgn(x2-x1)*i]++
	}
}

// Count the overlapping points
func (d *diagram) countOverlaps() (count int) {

	for _, row := range d {
		for _, cell := range row {
			// If two vents intercepts in the cell
			if cell > 1 {
				count++
			}
		}
	}

	return count
}

// HELPERS FUNCTIONS

// Absolute value
func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Returns the two arguments with the smallest first
func minmax(a, b int) (min, max int) {
	if a > b {
		return b, a
	}
	return a, b
}

// The sign function (1 if argument is positive, -1 else)
func sgn(x int) (sign int) {
	if x >= 0 {
		return 1
	}

	return -1
}
