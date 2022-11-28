package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 9 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	floor, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", computeRiskLevels(floor))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 9 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	floor, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	basins := computeBasins(floor)

	_, err = fmt.Fprintf(answer, "%d", multiplyThreeBiggestBasinsSize(basins))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Floor [][]int

type Coordinates struct {
	row, col int
}

type Basin struct {
	border    []Coordinates
	locations []Coordinates
	size      int
}

// INPUT PARSING
func parseLines(lines []string) (floor Floor, err error) {
	for _, line := range lines {
		var row []int
		for _, char := range line {
			number, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("error parsing char %s : %w", string(char), err)
			}

			row = append(row, number)
		}
		floor = append(floor, row)
	}
	return floor, nil
}

// PROCESING FUNCTIONS

// Check if n is lower than every number in numbers
func isLowerThan(n int, numbers []int) bool {
	for _, number := range numbers {
		if n >= number {
			return false
		}
	}

	return true
}

// Returns a list of all neighbors of a given position
func getNeighbors(row, col int, floor Floor) (numbers []int) {

	if row > 0 {
		numbers = append(numbers, floor[row-1][col])
	}

	if row < len(floor)-1 {
		numbers = append(numbers, floor[row+1][col])
	}

	if col > 0 {
		numbers = append(numbers, floor[row][col-1])
	}

	if col < len(floor[row])-1 {
		numbers = append(numbers, floor[row][col+1])
	}

	return numbers
}

// Compute risk levels
func computeRiskLevels(floor Floor) (risk int) {
	for row := range floor {
		for col := range floor[row] {
			if isLowerThan(floor[row][col], getNeighbors(row, col, floor)) {
				risk += 1 + floor[row][col]
			}
		}
	}

	return risk
}

// Compute the basins recursively
func computeBasins(floor Floor) (basins []Basin) {
	// Array to store if the locations is already in a basin
	var isAffected [][]bool

	for row := range floor {
		var isAffectedRow []bool
		for col := range floor[row] {
			if floor[row][col] == 9 {
				isAffectedRow = append(isAffectedRow, true)
			} else {
				isAffectedRow = append(isAffectedRow, false)
			}
		}
		isAffected = append(isAffected, isAffectedRow)
	}

	// Initialise the basins with the low points

	for row := range floor {
		for col := range floor[row] {
			if isLowerThan(floor[row][col], getNeighbors(row, col, floor)) {
				coords := []Coordinates{{row: row, col: col}}
				basins = append(basins, Basin{
					border:    coords,
					locations: coords,
					size:      1,
				})
				isAffected[row][col] = true
			}
		}
	}

	// Flag the know when the situation is not evolving
	stop := false

	// helper
	helper := func() {

		stop = true

		for idx, basin := range basins {
			var newBorder []Coordinates
			for _, borderLocation := range basin.border {
				for _, neighbordCoordinates := range getNeighborsCoordinates(borderLocation, floor) {
					if !isAffected[neighbordCoordinates.row][neighbordCoordinates.col] && floor[neighbordCoordinates.row][neighbordCoordinates.col] > floor[borderLocation.row][borderLocation.col] {
						isAffected[neighbordCoordinates.row][neighbordCoordinates.col] = true
						newBorder = append(newBorder, neighbordCoordinates)
						basin.locations = append(basin.locations, neighbordCoordinates)
						basin.size++
						stop = false
					}
				}
			}
			basin.border = newBorder
			basins[idx] = basin
		}

	}

	// Call the recursive function
	for !stop {
		helper()
	}

	return basins
}

// Get Neighbors coordinates
func getNeighborsCoordinates(coor Coordinates, floor Floor) (neighborsCoordinates []Coordinates) {

	if coor.row > 0 {
		neighborsCoordinates = append(neighborsCoordinates, Coordinates{coor.row - 1, coor.col})
	}

	if coor.row < len(floor)-1 {
		neighborsCoordinates = append(neighborsCoordinates, Coordinates{coor.row + 1, coor.col})
	}

	if coor.col > 0 {
		neighborsCoordinates = append(neighborsCoordinates, Coordinates{coor.row, coor.col - 1})
	}

	if coor.col < len(floor[coor.row])-1 {
		neighborsCoordinates = append(neighborsCoordinates, Coordinates{coor.row, coor.col + 1})
	}

	return neighborsCoordinates
}

// Find the 3 biggest basins
func multiplyThreeBiggestBasinsSize(basins []Basin) (product int) {

	if len(basins) < 3 {
		return product
	}

	// Sort the basins per size decreasing
	for i := 0; i < len(basins); i++ {
		for j := 1; j < len(basins); j++ {
			if basins[j-1].size < basins[j].size {
				basins[j-1], basins[j] = basins[j], basins[j-1]
			}
		}
	}

	product = 1

	for _, basin := range basins[:3] {
		product *= basin.size
	}

	return product
}
