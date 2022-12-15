package fabienz

import (
	"fmt"
	"io"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the paths.
	paths, err := pathsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse paths: %w", err)
	}

	// Create the cave.
	cave := caveFromPaths(paths)

	// Pour the sand until it fells infinitely.
	var counter int
	for stopped := true; stopped; stopped = cave.pourSand(helpers.Coord2D{X: 500, Y: 0}) {
		counter++
	}

	// Remove one sand because we also count the first sand falling infinitely.
	_, err = fmt.Fprintf(answer, "%d", counter-1)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the paths.
	paths, err := pathsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse paths: %w", err)
	}

	// Create the cave.
	cave := caveFromPaths(paths)

	// Pour the sand until there is no sand at the coord 500,0.
	var counter int
	for cave.Map[helpers.Coord2D{X: 500, Y: 0}] != 2 {
		cave.pourSandWithFloor(helpers.Coord2D{X: 500, Y: 0})
		counter++
	}

	_, err = fmt.Fprintf(answer, "%d", counter)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type RockPath []helpers.Coord2D

// 0 will be empty, 1 will be rock, 2 will be sand
type Cave struct {
	Map      helpers.IntGrid2D
	maxDepth int
}

func pathsFromLines(lines []string) ([]RockPath, error) {
	var paths []RockPath
	for _, line := range lines {
		var path RockPath
		for _, c := range strings.Split(line, " -> ") {
			var x, y int
			if _, err := fmt.Sscanf(c, "%d,%d", &x, &y); err != nil {
				return nil, fmt.Errorf("could not parse coord %q: %w", c, err)
			}
			path = append(path, helpers.Coord2D{X: x, Y: y})
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func caveFromPaths(paths []RockPath) Cave {
	Map := make(helpers.IntGrid2D)
	maxDepth := 0
	for _, path := range paths {
		c := path[0]
		Map[c] = 1
		if c.Y > maxDepth {
			maxDepth = c.Y
		}
		for i := 1; i < len(path); i++ {
			newC := path[i]
			switch {
			case newC.X == c.X && newC.Y > c.Y:
				for y := c.Y + 1; y <= newC.Y; y++ {
					Map[helpers.Coord2D{X: c.X, Y: y}] = 1
				}
			case newC.X == c.X && newC.Y < c.Y:
				for y := c.Y - 1; y >= newC.Y; y-- {
					Map[helpers.Coord2D{X: c.X, Y: y}] = 1
				}
			case newC.Y == c.Y && newC.X > c.X:
				for x := c.X + 1; x <= newC.X; x++ {
					Map[helpers.Coord2D{X: x, Y: c.Y}] = 1
				}
			case newC.Y == c.Y && newC.X < c.X:
				for x := c.X - 1; x >= newC.X; x-- {
					Map[helpers.Coord2D{X: x, Y: c.Y}] = 1
				}
			}
			c = newC
			if c.Y > maxDepth {
				maxDepth = c.Y
			}
		}
	}
	return Cave{Map: Map, maxDepth: maxDepth}
}

// Pour one unit of sand from a given coord. Returns true if the sand comes to a standstill
func (cave *Cave) pourSand(c helpers.Coord2D) bool {
	// The sand will fall endlessly if it has depth > maxDepth .
	for c.Y <= cave.maxDepth {
		switch {
		// If the coord below is empty, we fall down
		case cave.Map[helpers.Coord2D{X: c.X, Y: c.Y + 1}] == 0:
			c.Y++
		// Otherwise we try to go bottom-left
		case cave.Map[helpers.Coord2D{X: c.X - 1, Y: c.Y + 1}] == 0:
			c.X--
			c.Y++
		// Otherwise we try to go bottom-right
		case cave.Map[helpers.Coord2D{X: c.X + 1, Y: c.Y + 1}] == 0:
			c.X++
			c.Y++
		default:
			// If we can't go down, we are stuck
			cave.Map[c] = 2
			return true
		}
	}
	return false
}

func (cave *Cave) pourSandWithFloor(c helpers.Coord2D) bool {
	for c.Y < cave.maxDepth+1 {
		switch {
		// If the coord below is empty, we fall down
		case cave.Map[helpers.Coord2D{X: c.X, Y: c.Y + 1}] != 1 && cave.Map[helpers.Coord2D{X: c.X, Y: c.Y + 1}] != 2:
			c.Y++
		// Otherwise we try to go bottom-left
		case cave.Map[helpers.Coord2D{X: c.X - 1, Y: c.Y + 1}] != 1 && cave.Map[helpers.Coord2D{X: c.X - 1, Y: c.Y + 1}] != 2:
			c.X--
			c.Y++
		// Otherwise we try to go bottom-right
		case cave.Map[helpers.Coord2D{X: c.X + 1, Y: c.Y + 1}] != 1 && cave.Map[helpers.Coord2D{X: c.X + 1, Y: c.Y + 1}] != 2:
			c.X++
			c.Y++
		default:
			// If we can't go down, we are stuck
			cave.Map[c] = 2
			return true
		}
	}

	cave.Map[c] = 2

	return cave.Map[helpers.Coord2D{X: 500, Y: 0}] != 2
}
