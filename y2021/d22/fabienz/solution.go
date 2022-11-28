package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 22 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	reactor := make(Reactor)

	reactor.applyAll(instructions)

	_, err = fmt.Fprintf(answer, "%d", reactor.countOnInCube(Cube{-50, -50, -50, 50, 50, 50}))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	reactor := make(Reactor)

	reactor.applyAll(instructions)

	_, err = fmt.Fprintf(answer, "%d", reactor.countOn())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Cube struct {
	minX, minY, minZ, maxX, maxY, maxZ int
}

type Instruction struct {
	power bool
	cube  Cube
}

type Reactor map[Cube]bool

// type Reactor map[[3]int]bool

// INPUT PARSING

// We use regex to parse input
var re = regexp.MustCompile(`^([a-z]{2,3}) x=([0-9-]+)..([0-9-]+),y=([0-9-]+)..([0-9-]+),z=([0-9-]+)..([0-9-]+)$`)

func parseLines(lines []string) (instructions []Instruction, err error) {

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if len(match) < 8 {
			return instructions, fmt.Errorf("error parsing line %s : not enough values found", line)
		}

		var instruction Instruction

		instruction.power = match[1] == "on"

		minX, err := strconv.Atoi(match[2])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.minX = minX

		maxX, err := strconv.Atoi(match[3])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.maxX = maxX

		minY, err := strconv.Atoi(match[4])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.minY = minY

		maxY, err := strconv.Atoi(match[5])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.maxY = maxY

		minZ, err := strconv.Atoi(match[6])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.minZ = minZ

		maxZ, err := strconv.Atoi(match[7])
		if err != nil {
			return instructions, fmt.Errorf("error parsing line %s : %w", line, err)
		}
		instruction.cube.maxZ = maxZ

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

// PROCESSING FUNCTIONS

// Apply an instruction to the reactor
// func (r Reactor) apply(i Instruction, maxCoor int) {
// 	if i.cube.minX < -maxCoor {
// 		i.cube.minX = -maxCoor
// 	}

// 	if i.cube.maxX > maxCoor {
// 		i.cube.maxX = maxCoor
// 	}

// 	if i.cube.minY < -maxCoor {
// 		i.cube.minY = -maxCoor
// 	}

// 	if i.cube.maxY > maxCoor {
// 		i.cube.maxY = maxCoor
// 	}

// 	if i.cube.minZ < -maxCoor {
// 		i.cube.minZ = -maxCoor
// 	}

// 	if i.cube.maxZ > maxCoor {
// 		i.cube.maxZ = maxCoor
// 	}

// 	// Switch the state of every cube
// 	for x := i.cube.minX; x<= i.cube.maxX; x++ {
// 		for y := i.cube.minY; y<=i.cube.maxY; y++ {
// 			for z := i.cube.minZ; z<=i.cube.maxZ; z++ {
// 				r[[3]int{x,y,z}] = i.power
// 			}
// 		}
// 	}
// }

// // Count all cubes powered on in a range of maxCoor coordinates
// func (r Reactor) countOn(maxCoor int) (count int) {
// 	for coor, isOn := range r {
// 		if coor[0] >= -maxCoor && coor[0] <= maxCoor && coor[1] >= -maxCoor && coor[1] <= maxCoor && coor[2] >= -maxCoor && coor[2] <= maxCoor {
// 			if isOn {
// 				count++
// 			}
// 		}
// 	}

// 	return count
// }

// Apply an instruction
func (r Reactor) apply(ins Instruction) {
	// Check the previous cubes to see if the cubes intersects
	for cube, isOn := range r {
		// If the two cubes overlap
		if !isNotOverlapping(cube, ins.cube) {
			// Remove the original cube from the map
			delete(r, cube)

			// Compute the new cubes not overlapping
			cubes := split(cube, ins.cube)
			for _, c := range cubes {
				r[c] = isOn
			}
		}
	}

	// Add the new cube to the map of cubes
	r[ins.cube] = ins.power
}

// Apply a slice of instructions
func (r Reactor) applyAll(instructions []Instruction) {
	for _, instruction := range instructions {
		r.apply(instruction)
	}
}

// Checks if two cubes overlaps
func isNotOverlapping(c1, c2 Cube) bool {

	return c1.minX > c2.maxX || c1.maxX < c2.minX || c1.minY > c2.maxY || c1.maxY < c2.minY || c1.minZ > c2.maxZ || c1.maxZ < c2.minZ
}

// Split cubes if needed to have distinct cubes not overlapping with the reference Cube
func split(cube Cube, refCube Cube) (newCubes []Cube) {

	if refCube.minX > cube.minX {
		newCubes = append(newCubes, Cube{cube.minX, cube.minY, cube.minZ, refCube.minX - 1, cube.maxY, cube.maxZ})
	}

	if refCube.maxX < cube.maxX {
		newCubes = append(newCubes, Cube{refCube.maxX + 1, cube.minY, cube.minZ, cube.maxX, cube.maxY, cube.maxZ})
	}

	if refCube.minY > cube.minY {
		newCubes = append(newCubes, Cube{max(refCube.minX, cube.minX), cube.minY, cube.minZ, min(refCube.maxX, cube.maxX), refCube.minY - 1, cube.maxZ})
	}

	if refCube.maxY < cube.maxY {
		newCubes = append(newCubes, Cube{max(refCube.minX, cube.minX), refCube.maxY + 1, cube.minZ, min(refCube.maxX, cube.maxX), cube.maxY, cube.maxZ})
	}

	if refCube.minZ > cube.minZ {
		newCubes = append(newCubes, Cube{max(refCube.minX, cube.minX), max(refCube.minY, cube.minY), cube.minZ, min(refCube.maxX, cube.maxX), min(refCube.maxY, cube.maxY), refCube.minZ - 1})
	}

	if refCube.maxZ < cube.maxZ {
		newCubes = append(newCubes, Cube{max(refCube.minX, cube.minX), max(refCube.minY, cube.minY), refCube.maxZ + 1, min(refCube.maxX, cube.maxX), min(refCube.maxY, cube.maxY), cube.maxZ})
	}

	return newCubes
}

// Find the minimum of two integers
func min(a, b int) (min int) {
	if a > b {
		return b
	}
	return a
}

// Find the maximum of two integers
func max(a, b int) (max int) {
	if a < b {
		return b
	}
	return a
}

// Count the cubes on in a given zone
func (r Reactor) countOnInCube(zone Cube) (count int) {
	for cube, isOn := range r {
		// If the cube is intersecting the zone
		if !isNotOverlapping(cube, zone) {
			if isOn {
				count += cube.volume()
			}
		}
	}

	return count
}

// Count all the cubes powered on
func (r Reactor) countOn() (count int) {
	for cube, isOn := range r {
		if isOn {
			count += cube.volume()
		}
	}

	return count
}

// Compute the volume of a cube
func (c Cube) volume() (vol int) {
	return (c.maxX - c.minX + 1) * (c.maxY - c.minY + 1) * (c.maxZ - c.minZ + 1)
}
