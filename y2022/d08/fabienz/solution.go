package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the forest
	forest, err := parseForest(lines)
	if err != nil {
		return fmt.Errorf("could not parse forest: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", forest.CountVisibleTrees())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the forest
	forest, err := parseForest(lines)
	if err != nil {
		return fmt.Errorf("could not parse forest: %w", err)
	}

	// Compute the scenic score for each tree
	scenicScores := []int{}

	for y := 0; y < forest.Height; y++ {
		for x := 0; x < forest.Width; x++ {
			scenicScores = append(scenicScores, forest.ScenicScore(Coord{x, y}))
		}
	}

	_, err = fmt.Fprintf(answer, "%d", helpers.MaxInts(scenicScores))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Cartesian coordinates
type Coord struct {
	X, Y int
}

// We will store the forest in a map
type Forest struct {
	Map           map[Coord]int
	Height, Width int
}

// Parse the input into a Forest
func parseForest(lines []string) (Forest, error) {
	forest := Forest{
		Map:    make(map[Coord]int),
		Height: len(lines),
		Width:  len(lines[0]),
	}

	for y, line := range lines {
		for x, char := range line {
			height, err := strconv.Atoi(string(char))
			if err != nil {
				return Forest{}, fmt.Errorf("could not parse forest: %w", err)
			}
			forest.Map[Coord{x, y}] = height
		}
	}

	return forest, nil
}

// Check if the tree is visible from at least one of the 4 directions up, down, left, right.
func (f Forest) IsVisible(c Coord, directionsToConsider []Coord) bool {
	for _, direction := range directionsToConsider {
		if f.IsVisibleInDirection(c, direction) {
			return true
		}
	}

	return false
}

// Check in a given direction if the tree is visible
func (f Forest) IsVisibleInDirection(c Coord, direction Coord) bool {
	// Fetch the height of the tree in pos c
	height := f.Map[c]

	// Check if there is a tree in the given direction
	for {
		c = Coord{c.X + direction.X, c.Y + direction.Y}
		if c.X < 0 || c.X >= f.Width || c.Y < 0 || c.Y >= f.Height {
			return true
		}
		if f.Map[c] >= height {
			return false
		}
	}
}

// Count the number of visible trees
func (f Forest) CountVisibleTrees() int {
	directionsToConsider := []Coord{
		{1, 0},  // right
		{-1, 0}, // left
		{0, 1},  // down
		{0, -1}, // up
	}
	count := 0
	for c := range f.Map {
		if f.IsVisible(c, directionsToConsider) {
			count++
		}
	}
	return count
}

// Compute the viewing distance in a given direction
func (f Forest) ViewingDistanceInDirection(c Coord, direction Coord) (viewingDistance int) {
	// Fetch the height of the tree in pos c
	height := f.Map[c]

	// Check if there is a tree in the given direction
	for {
		c = Coord{c.X + direction.X, c.Y + direction.Y}
		if c.X < 0 || c.X >= f.Width || c.Y < 0 || c.Y >= f.Height {
			return viewingDistance
		}
		viewingDistance++
		if f.Map[c] >= height {
			return viewingDistance
		}
	}
}

// Compute the scenic score of a given tree
func (f Forest) ScenicScore(c Coord) int {
	directionsToConsider := []Coord{
		{1, 0},  // right
		{-1, 0}, // left
		{0, 1},  // down
		{0, -1}, // up
	}
	scenicScore := 1
	for _, direction := range directionsToConsider {
		scenicScore *= f.ViewingDistanceInDirection(c, direction)
	}
	return scenicScore
}
