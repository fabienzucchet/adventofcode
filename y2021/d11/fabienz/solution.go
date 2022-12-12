package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// Size of yhe square of octopuses
const OCTOPUSESSIZE = 10

// PartOne solves the first problem of day 11 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	octopuses, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	var flashesCount int

	for i := 0; i < 100; i++ {
		flashesCount += octopuses.iterate()
	}

	_, err = fmt.Fprintf(answer, "%d", flashesCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	octopuses, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	var step int

	for octopuses.iterate() < OCTOPUSESSIZE*OCTOPUSESSIZE {
		step++
	}

	_, err = fmt.Fprintf(answer, "%d", step+1)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Octopuses [OCTOPUSESSIZE][OCTOPUSESSIZE]int

type Flashes [OCTOPUSESSIZE][OCTOPUSESSIZE]bool

type Coordinates struct {
	row, col int
}

// Parsing Input
func parseLines(lines []string) (octopuses Octopuses, err error) {

	for row, line := range lines {
		for col, char := range line {
			energy, err := strconv.Atoi(string(char))
			if err != nil {
				return octopuses, fmt.Errorf("error parsing energy level %s : %w", string(char), err)
			}

			octopuses[row][col] = energy
		}
	}

	return octopuses, nil
}

// Iterate on step
func (o *Octopuses) iterate() (flashesCount int) {
	// Increase the energy level of all octopuses by 1
	for row := range o {
		for col := range o[row] {
			o[row][col]++
		}
	}

	// Use a bool array to prevent from flashing twice
	var flashes Flashes

	for i := 0; i < 100; i++ {
		for row := range o {
			for col := range o[row] {
				if o[row][col] > 9 && !flashes[row][col] {
					flashes[row][col] = true
					for _, coor := range o.getNeighbors(row, col) {
						o[coor.row][coor.col]++
					}
				}
			}
		}
	}

	// Reset the energy level of the flashing octopuses to 0
	for row := range o {
		for col := range o[row] {
			if flashes[row][col] {
				o[row][col] = 0
				flashesCount++
			}
		}
	}

	return flashesCount
}

// Get then neighbors coordinates
func (o *Octopuses) getNeighbors(row, col int) (coordinates []Coordinates) {

	if row > 0 {
		coordinates = append(coordinates, Coordinates{row - 1, col})
		if col > 0 {
			coordinates = append(coordinates, Coordinates{row - 1, col - 1})
		}
		if col < len(o[row])-1 {
			coordinates = append(coordinates, Coordinates{row - 1, col + 1})
		}
	}
	if row < len(o)-1 {
		coordinates = append(coordinates, Coordinates{row + 1, col})
		if col > 0 {
			coordinates = append(coordinates, Coordinates{row + 1, col - 1})
		}
		if col < len(o[row])-1 {
			coordinates = append(coordinates, Coordinates{row + 1, col + 1})
		}
	}
	if col > 0 {
		coordinates = append(coordinates, Coordinates{row, col - 1})
	}
	if col < len(o[row])-1 {
		coordinates = append(coordinates, Coordinates{row, col + 1})
	}

	return coordinates
}
