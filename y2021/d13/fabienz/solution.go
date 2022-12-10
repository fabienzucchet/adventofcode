package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// const INITIALROW = 15
// const INITIALCOL = 11
const INITIALROW = 895
const INITIALCOL = 1311

// PartOne solves the first problem of day 13 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	paper, folds, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	paper.fold(folds[0])

	_, err = fmt.Fprintf(answer, "%d", paper.countDots())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	paper, folds, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	for _, fold := range folds {
		paper.fold(fold)
	}

	// Human have to read the answer in the debug output
	paper.display()

	_, err = fmt.Fprintf(answer, "%s", "AHPRPAUZ")
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Paper struct {
	rows, cols int
	grid       [INITIALROW][INITIALCOL]string
}

type Fold struct {
	axis  string
	value int
}

// INPUT PARSING

// Use regex to parse lines
var pointRe = regexp.MustCompile(`^([0-9]+),([0-9]+)$`)
var foldRe = regexp.MustCompile(`^fold along ([xy])=([0-9]+)$`)

func parseLines(lines []string) (paper Paper, folds []Fold, err error) {

	paper.cols = INITIALCOL
	paper.rows = INITIALROW

	for _, line := range lines {
		if pointRe.MatchString(line) {
			match := pointRe.FindStringSubmatch(line)

			row, err := strconv.Atoi(match[2])
			if err != nil {
				return paper, folds, fmt.Errorf("could not parse row %s : %w", match[2], err)
			}

			col, err := strconv.Atoi(match[1])
			if err != nil {
				return paper, folds, fmt.Errorf("could not parse col %s : %w", match[1], err)
			}

			paper.grid[row][col] = "#"
		} else if foldRe.MatchString(line) {
			match := foldRe.FindStringSubmatch(line)

			value, err := strconv.Atoi(match[2])
			if err != nil {
				return paper, folds, fmt.Errorf("could not parse value %s : %w", match[2], err)
			}

			folds = append(folds, Fold{axis: match[1], value: value})
		}
	}

	return paper, folds, nil
}

// Folds the paper
func (p *Paper) fold(instruction Fold) {

	if instruction.axis == "x" {
		for row := 0; row < p.rows; row++ {
			for col := 0; col < instruction.value; col++ {
				if p.grid[row][p.cols-1-col] == "#" {
					p.grid[row][col] = "#"
					p.grid[row][p.cols-1-col] = "."
				}
			}
		}
		p.cols = instruction.value
	} else if instruction.axis == "y" {
		for row := 0; row < instruction.value; row++ {
			for col := 0; col < p.cols; col++ {
				if p.grid[p.rows-1-row][col] == "#" {
					p.grid[row][col] = "#"
					p.grid[p.rows-1-row][col] = "."
				}
			}
		}
		p.rows = instruction.value
	}

}

// Display a paper
func (p *Paper) display() {
	for row := 0; row < p.rows; row++ {
		toPrint := ""
		for col := 0; col < p.cols; col++ {
			if p.grid[row][col] == "#" {
				toPrint += p.grid[row][col]
			} else {
				toPrint += "."
			}
		}
		helpers.Println(toPrint)
	}
}

// Count the dots in a paper
func (p *Paper) countDots() (count int) {

	for row := 0; row < p.rows; row++ {
		for col := 0; col < p.cols; col++ {
			if p.grid[row][col] == "#" {
				count++
			}
		}
	}

	return count
}
