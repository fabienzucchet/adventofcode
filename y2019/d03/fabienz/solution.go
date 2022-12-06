package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Store all points with a pipe in a map.
	pipes := make(map[Point]int)

	// Will contain the intersections bewteen the pipes.
	intersections := make([]Intersection, 0)

	// For each pipe, parse the input and update the map.
	for _, line := range lines {
		// Parse the pipe.
		pipe := parsePipe(line)

		// Update the map.
		pipes, intersections = updatePipes(pipes, pipe, intersections)
	}

	// Compute the manhattan distance between the origin and the closest intersection.
	distances := make([]int, len(intersections))
	for i, intersection := range intersections {
		distances[i] = manhattanDistance(intersection.point, Point{x: 0, y: 0})
	}

	_, err = fmt.Fprintf(answer, "%d", min(distances))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Store all points with a pipe in a map.
	pipes := make(map[Point]int)

	// Will contain the intersections bewteen the pipes.
	intersections := make([]Intersection, 0)

	// For each pipe, parse the input and update the map.
	for _, line := range lines {
		// Parse the pipe.
		pipe := parsePipe(line)

		// Update the map.
		pipes, intersections = updatePipes(pipes, pipe, intersections)
	}

	// Compute the minimum number of steps.
	steps := make([]int, len(intersections))
	for i, intersection := range intersections {
		steps[i] = intersection.steps
	}

	// TODO: Write your solution to Part 2 below.
	_, err = fmt.Fprintf(answer, "%d", min(steps))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Struct to store a point in the grid.
type Point struct {
	x int
	y int
}

// Section of a wire.
type Section struct {
	direction string
	distance  int
}

// Store an intersection between two pipes.
type Intersection struct {
	point Point
	steps int
}

// Parse a pipe and return a slice of sections.
func parsePipe(pipe string) []Section {
	sections := make([]Section, 0)

	// Split the pipe into sections.
	sectionsStrings := strings.Split(pipe, ",")

	// Parse each section.
	for _, sectionString := range sectionsStrings {
		// Get the direction.
		direction := string(sectionString[0])

		// Get the distance.
		distance, err := strconv.Atoi(sectionString[1:])
		if err != nil {
			panic(err)
		}

		// Add the section to the slice.
		sections = append(sections, Section{direction: direction, distance: distance})
	}

	return sections
}

// Update the map of pipes with the new pipe.
func updatePipes(pipes map[Point]int, pipe []Section, intersections []Intersection) (map[Point]int, []Intersection) {
	// Current point.
	currentPoint := Point{x: 0, y: 0}

	// Work on a copy of the map to avoid counting intersections with the same pipe.
	newPipes := make(map[Point]int)
	for k, v := range pipes {
		newPipes[k] = v
	}

	// Keep track of the number of steps.
	steps := 0

	// For each section of the pipe.
	for _, section := range pipe {
		// For each point of the section.
		for i := 0; i < section.distance; i++ {
			steps++
			// Update the current point.
			switch section.direction {
			case "U":
				currentPoint.y++
			case "D":
				currentPoint.y--
			case "L":
				currentPoint.x--
			case "R":
				currentPoint.x++
			}

			// Check if there was a pipe at this point. If yes, add it to the intersections.
			if prevSteps, ok := pipes[currentPoint]; ok {
				intersections = append(intersections, Intersection{point: currentPoint, steps: steps + prevSteps})
			}

			// Add the point to the map only of not already in it.
			if _, ok := newPipes[currentPoint]; !ok {
				newPipes[currentPoint] = steps
			}
		}
	}

	return newPipes, intersections
}

// Compute the absolute value of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Return the minimum of a slice of ints.
func min(slice []int) int {
	min := slice[0]
	for _, value := range slice {
		if value < min {
			min = value
		}
	}
	return min
}

// Compute the manhattan distance between two points.
func manhattanDistance(p1 Point, p2 Point) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}
