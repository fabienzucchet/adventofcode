package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	PocketDimension = make(map[Coordinates]Cube)
	neighboursIndexes = []int{-1, 0, 1}

	parseLines(lines)

	for i := 0; i < 6; i++ {
		iterateCycle()
	}

	_, err = fmt.Fprintf(answer, "%d", countActiveCubes())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	PocketDimension4D = make(map[Coordinates4D]Cube)
	neighboursIndexes = []int{-1, 0, 1}

	parseLines4D(lines)

	for i := 0; i < 6; i++ {
		iterateCycle4D()
	}

	_, err = fmt.Fprintf(answer, "%d", countActiveCubes4D())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Coordinates struct {
	x int
	y int
	z int
}

type Coordinates4D struct {
	x int
	y int
	z int
	w int
}

type Cube struct {
	isActive  bool
	wasActive bool
}

var PocketDimension map[Coordinates]Cube
var PocketDimension4D map[Coordinates4D]Cube

var neighboursIndexes []int

func generateNeighborCoordinates(coor Coordinates) []Coordinates {
	var res []Coordinates

	for _, dx := range neighboursIndexes {
		for _, dy := range neighboursIndexes {
			for _, dz := range neighboursIndexes {
				if dx != 0 || dy != 0 || dz != 0 {
					res = append(res, Coordinates{coor.x + dx, coor.y + dy, coor.z + dz})
				}
			}
		}
	}

	return res
}

func generateNeighborCoordinates4D(coor Coordinates4D) []Coordinates4D {
	var res []Coordinates4D

	for _, dx := range neighboursIndexes {
		for _, dy := range neighboursIndexes {
			for _, dz := range neighboursIndexes {
				for _, dw := range neighboursIndexes {
					if dx != 0 || dy != 0 || dz != 0 || dw != 0 {
						res = append(res, Coordinates4D{coor.x + dx, coor.y + dy, coor.z + dz, coor.w + dw})
					}
				}
			}
		}
	}

	return res
}

// Read the lines from input and create the associated Cubes
func parseLines(lines []string) {

	// Create the cube
	for y, line := range lines {
		for x, rune := range line {

			// Create the cube
			c := Cube{}

			// Init the state
			if rune == '#' {
				c.isActive = true
				c.wasActive = true
			}

			coor := Coordinates{x, y, 0}

			PocketDimension[coor] = c
		}
	}

}

// Read the lines from input and create the associated Cubes
func parseLines4D(lines []string) {

	// Create the cube
	for y, line := range lines {
		for x, rune := range line {

			// Create the cube
			c := Cube{}

			// Init the state
			if rune == '#' {
				c.isActive = true
				c.wasActive = true
			}

			coor := Coordinates4D{x, y, 0, 0}

			PocketDimension4D[coor] = c
		}
	}

}

func countWasActiveNeighbors(coor Coordinates) int {
	count := 0

	for _, neighCoor := range generateNeighborCoordinates(coor) {
		if PocketDimension[neighCoor].wasActive {
			count++
		}
	}

	return count
}

func countWasActiveNeighbors4D(coor Coordinates4D) int {
	count := 0

	for _, neighCoor := range generateNeighborCoordinates4D(coor) {
		if PocketDimension4D[neighCoor].wasActive {
			count++
		}
	}

	return count
}

func iterateCycle() {

	// Make sure that wasActive corresponds to the state prior the iteration and isActive to the state after the iteration
	for coor, c := range PocketDimension {
		c.wasActive = c.isActive
		PocketDimension[coor] = c
	}

	// Make sure that every cube that has at least 1 active neighbor is in the map
	for coor, c := range PocketDimension {
		if c.wasActive {
			for _, neighborCoor := range generateNeighborCoordinates(coor) {
				_, present := PocketDimension[neighborCoor]

				if !present {
					PocketDimension[neighborCoor] = Cube{isActive: false, wasActive: false}
				}
			}
		}
	}

	// Update all cubes
	for coor, c := range PocketDimension {
		if nbActive := countWasActiveNeighbors(coor); c.wasActive && nbActive != 2 && nbActive != 3 {
			c.isActive = false
			PocketDimension[coor] = c
		} else if nbActive := countWasActiveNeighbors(coor); !c.wasActive && nbActive == 3 {
			c.isActive = true
			PocketDimension[coor] = c
		}
	}
}

func iterateCycle4D() {

	// Make sure that wasActive corresponds to the state prior the iteration and isActive to the state after the iteration
	for coor, c := range PocketDimension4D {
		c.wasActive = c.isActive
		PocketDimension4D[coor] = c
	}

	// Make sure that every cube that has at least 1 active neighbor is in the map
	for coor, c := range PocketDimension4D {
		if c.wasActive {
			for _, neighborCoor := range generateNeighborCoordinates4D(coor) {
				_, present := PocketDimension4D[neighborCoor]

				if !present {
					PocketDimension4D[neighborCoor] = Cube{isActive: false, wasActive: false}
				}
			}
		}
	}

	// Update all cubes
	for coor, c := range PocketDimension4D {
		if nbActive := countWasActiveNeighbors4D(coor); c.wasActive && nbActive != 2 && nbActive != 3 {
			c.isActive = false
			PocketDimension4D[coor] = c
		} else if nbActive := countWasActiveNeighbors4D(coor); !c.wasActive && nbActive == 3 {
			c.isActive = true
			PocketDimension4D[coor] = c
		}
	}
}

func countActiveCubes() int {
	count := 0

	for _, c := range PocketDimension {
		if c.isActive {
			count++
		}
	}

	return count
}

func countActiveCubes4D() int {
	count := 0

	for _, c := range PocketDimension4D {
		if c.isActive {
			count++
		}
	}

	return count
}
