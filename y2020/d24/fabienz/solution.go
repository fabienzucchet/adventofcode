package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 24 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions := parseLines(lines)

	t := generateTilingAfterInstructions(instructions)

	_, err = fmt.Fprintf(answer, "%d", t.countBlack())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions := parseLines(lines)

	t := generateTilingAfterInstructions(instructions)

	// Iterate 100 days on the tiling
	for i := 0; i < 100; i++ {
		t = t.iterate()
	}

	_, err = fmt.Fprintf(answer, "%d", t.countBlack())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

/*
* We will use hexagonal coordinates to identify uniquely the tiles.
* See the second example of https://homepages.inf.ed.ac.uk/rbf/CVonline/LOCAL_COPIES/AV0405/MARTIN/Hex.pdf for an example
*
* In this coordinates system, every tile has a position depending on 3 coordinates (x, y, z)
*
* The translation of the problem movements into coordinates is :
* ne -> ( 0, 1, 1)
* nw -> (-1, 0, 1)
*  w -> (-1,-1, 0)
* sw -> ( 0,-1,-1)
* se -> ( 1, 0,-1)
*  e -> ( 1, 1, 0)
 */

// TYPES
type HexaCoordinates struct {
	x, y, z int
}

// Represents the tiles : a black tile is a true value in the map
type tiling map[HexaCoordinates]bool

// PARSE INPUT

// Transform an instruction string into a slice of HexaCoordinates corresponding to the instructions
func parseInstruction(instruction string) (coor []HexaCoordinates) {
	i := 0

	for i < len(instruction) {
		switch {
		case instruction[i:i+1] == "e":
			coor = append(coor, HexaCoordinates{1, 1, 0})
			i++
		case instruction[i:i+1] == "w":
			coor = append(coor, HexaCoordinates{-1, -1, 0})
			i++
		case instruction[i:i+2] == "ne":
			coor = append(coor, HexaCoordinates{0, 1, 1})
			i += 2
		case instruction[i:i+2] == "nw":
			coor = append(coor, HexaCoordinates{-1, 0, 1})
			i += 2
		case instruction[i:i+2] == "se":
			coor = append(coor, HexaCoordinates{1, 0, -1})
			i += 2
		case instruction[i:i+2] == "sw":
			coor = append(coor, HexaCoordinates{0, -1, -1})
			i += 2
		}
	}

	return coor
}

// Parse all input lines into a list of list of Coordinates
func parseLines(lines []string) (coordinates [][]HexaCoordinates) {

	for _, line := range lines {
		coordinates = append(coordinates, parseInstruction(line))
	}

	return coordinates
}

// PROCESSING FUNCTIONS

// Execute an instruction
func executeInstruction(pos HexaCoordinates, instruction HexaCoordinates) (newPos HexaCoordinates) {
	newPos.x = pos.x + instruction.x
	newPos.y = pos.y + instruction.y
	newPos.z = pos.z + instruction.z

	return newPos
}

// Execute recursively a serie of HexaCoordinates movements from a given position
func executeInstructions(pos HexaCoordinates, instructions []HexaCoordinates) (newPos HexaCoordinates) {

	if len(instructions) == 0 {
		return pos
	}

	newPos = executeInstruction(pos, instructions[0])

	return executeInstructions(newPos, instructions[1:])
}

// Generate tiling as described by the instructions
func generateTilingAfterInstructions(instructions [][]HexaCoordinates) (t tiling) {
	t = tiling{}

	for _, instruction := range instructions {
		tileCoor := executeInstructions(HexaCoordinates{}, instruction)
		t[tileCoor] = !t[tileCoor]
	}
	return t
}

// Count the black tiles in a tiling
func (t *tiling) countBlack() (count int) {

	for _, isBlack := range *t {
		if isBlack {
			count++
		}
	}

	return count
}

// Compute the neighbors of a give hexagonal coordinate
func getNeighbors(pos HexaCoordinates) (neighbors []HexaCoordinates) {

	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{1, 1, 0}))   // e
	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{-1, -1, 0})) // w
	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{0, 1, 1}))   // ne
	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{-1, 0, 1}))  // nw
	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{1, 0, -1}))  // se
	neighbors = append(neighbors, executeInstruction(pos, HexaCoordinates{0, -1, -1})) // sw

	return neighbors
}

// Count the black neighbors around a position
func (t *tiling) countBlackNeighbors(pos HexaCoordinates) (count int) {

	for _, neighbor := range getNeighbors(pos) {
		if (*t)[neighbor] {
			count++
		}
	}

	return count
}

// Iterate one step of swiching tiles
func (t *tiling) iterate() (newt tiling) {
	// First add all tiles that have a black neighbor (i.e. that could change color) to the tiling
	for coor, _ := range *t {
		for _, neighbor := range getNeighbors(coor) {
			if _, exists := (*t)[neighbor]; !exists {
				(*t)[neighbor] = false
			}
		}
	}

	// We can then build a new tiling considering the number of black neighbors for each tiles
	newt = tiling{}

	for coor, isBlack := range *t {
		blackNeighCount := t.countBlackNeighbors(coor)
		switch {
		case isBlack && blackNeighCount == 0:
			newt[coor] = false
		case isBlack && blackNeighCount > 2:
			newt[coor] = false
		case !isBlack && blackNeighCount == 2:
			newt[coor] = true
		case isBlack:
			newt[coor] = true
		}
	}

	return newt
}
