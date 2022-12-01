package fabienz

import (
	"fmt"
	"io"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the orbit map.
	orbitMap := ParseOrbitMap(lines)

	_, err = fmt.Fprintf(answer, "%d", orbitMap.CountOrbits())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the orbit map.
	orbitMap := ParseOrbitMap(lines)

	// Find the common ancestor between YOU and SAN.
	commonAncestor := orbitMap.FindCommonAncestor("YOU", "SAN")

	// Count the number of orbits between YOU and the common ancestor.
	count := 0
	for object := orbitMap["YOU"]; object != commonAncestor; object = orbitMap[object] {
		count++
	}

	// Count the number of orbits between SAN and the common ancestor.
	for object := orbitMap["SAN"]; object != commonAncestor; object = orbitMap[object] {
		count++
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// We store the direct orbit relationships in a map.
type OrbitMap map[string]string

// ParseOrbitMap parses the input and returns the orbit map.
func ParseOrbitMap(input []string) OrbitMap {
	orbitMap := make(OrbitMap)
	for _, orbit := range input {
		// Split the orbit into the two objects.
		objects := strings.Split(orbit, ")")
		orbitMap[objects[1]] = objects[0]
	}
	return orbitMap
}

// Helper function to count the number of indirect orbits for a given object.
func (o OrbitMap) CountOrbitsForObject(object string) int {
	count := 0
	if orbit, ok := o[object]; ok {
		count++
		count += o.CountOrbitsForObject(orbit)
	}
	return count
}

// CountOrbits counts the number of direct and indirect orbits.
func (o OrbitMap) CountOrbits() int {
	count := 0
	for _, object := range o {
		count++
		count += o.CountOrbitsForObject(object)
	}
	return count
}

// Find the common ancestor of two objects.
func (o OrbitMap) FindCommonAncestor(object1, object2 string) string {
	ancestors := make(map[string]bool)
	for object := object1; object != ""; object = o[object] {
		ancestors[object] = true
	}

	for object := object2; object != ""; object = o[object] {
		if ancestors[object] {
			return object
		}
	}

	return ""
}
