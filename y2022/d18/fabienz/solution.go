package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the cubes from the input
	cubes, _, err := cubesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse cubes: %w", err)
	}

	// Update the number of hidden faces for all cubes
	for k, v := range cubes {
		v.facesHidden = v.countHiddenFaces(cubes)
		cubes[k] = v
	}

	// Sum all the visible faces
	var visibleFaces int
	for _, c := range cubes {
		visibleFaces += 6 - c.facesHidden
	}

	_, err = fmt.Fprintf(answer, "%d", visibleFaces)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the cubes from the input
	cubes, maxPos, err := cubesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse cubes: %w", err)
	}

	// Create a water map.
	water := percolate(cubes, maxPos)

	count := 0
	for _, c := range cubes {
		count += c.countWaterFaces(water)
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type cube struct {
	pos         helpers.Coord3D
	facesHidden int
}

// parse the input and return a map of cubes
func cubesFromLines(lines []string) (map[helpers.Coord3D]cube, helpers.Coord3D, error) {
	cubes := make(map[helpers.Coord3D]cube)
	var maxPos helpers.Coord3D

	for _, line := range lines {
		var x, y, z int
		_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			return nil, helpers.Coord3D{}, fmt.Errorf("could not parse line: %w", err)
		}
		if x > maxPos.X {
			maxPos.X = x
		}
		if y > maxPos.Y {
			maxPos.Y = y
		}
		if z > maxPos.Z {
			maxPos.Z = z
		}

		pos := helpers.Coord3D{X: x, Y: y, Z: z}
		cubes[pos] = cube{pos: pos}
	}

	// Add 1 to maxPos to account for the fact that the cubes are 1x1x1
	maxPos.X++
	maxPos.Y++
	maxPos.Z++

	return cubes, maxPos, nil
}

// Count the number of hidden faces of a cube
func (c cube) countHiddenFaces(cubes map[helpers.Coord3D]cube) int {
	count := 0
	for _, pos := range c.pos.NeighborsWithoutDiagonals() {
		if _, ok := cubes[pos]; ok {
			count++
		}
	}

	return count
}

// Given a map of cubes, use percolation to build a map of emplacements accessible to water
func percolate(cubes map[helpers.Coord3D]cube, maxPos helpers.Coord3D) map[helpers.Coord3D]bool {
	water := make(map[helpers.Coord3D]bool)

	q := []helpers.Coord3D{}

	init := helpers.Coord3D{X: -1, Y: -1, Z: -1}

	// add 0,0,0 to the queue
	q = append(q, init)
	water[init] = true

	// while the queue is not empty
	for len(q) > 0 {
		// Pop the first element of the queue.
		p := q[0]
		q = q[1:]

		// Add all the neighbors of p to the queue if they are not already in the water map and if they are not a cube.
		for _, pos := range p.NeighborsWithoutDiagonals() {
			// If pos is < 0, do not add to the queue.
			if pos.X < -1 || pos.Y < -1 || pos.Z < -1 {
				continue
			}

			// If pos is > maxPos, do not add to the queue.
			if pos.X > maxPos.X || pos.Y > maxPos.Y || pos.Z > maxPos.Z {
				continue
			}

			// If pos is in the water map, do not add to the queue.
			if _, ok := water[pos]; ok {
				continue
			}

			// If pos is a cube, do not add to the queue.
			if _, ok := cubes[pos]; ok {
				continue
			}

			// Add pos to the queue.
			q = append(q, pos)
			// Add p to the water map
			water[pos] = true
		}
	}

	return water
}

// Count the faces of a cube that are in contact with water
func (c cube) countWaterFaces(water map[helpers.Coord3D]bool) int {
	count := 0
	for _, pos := range c.pos.NeighborsWithoutDiagonals() {
		if _, ok := water[pos]; ok {
			count++
		}
	}

	return count
}
