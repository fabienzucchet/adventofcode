package fabienz

import (
	"fmt"
	"io"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

const CHAMBERSIZE = 7

// PartOne solves the first problem of day 17 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	height := findChamberHeight(2022, []rune(lines[0]))

	_, err = fmt.Fprintf(answer, "%d", height)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	height := findChamberHeight(1000000000000, []rune(lines[0]))

	_, err = fmt.Fprintf(answer, "%d", height)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type rock [][]rune

var rocksTypes = []rock{
	{
		{'#', '#', '#', '#'},
	},
	{
		{'.', '#', '.'},
		{'#', '#', '#'},
		{'.', '#', '.'},
	},
	{
		{'.', '.', '#'},
		{'.', '.', '#'},
		{'#', '#', '#'},
	},
	{
		{'#'},
		{'#'},
		{'#'},
		{'#'},
	},
	{
		{'#', '#'},
		{'#', '#'},
	},
}

type state struct {
	chamber     [][CHAMBERSIZE]rune
	rockTypeIdx int
	jetPattern  []rune
	jetIdx      int
	cache       map[key]cacheData
	iteration   int
}

type key struct {
	rockTypeIdx int
	jetIdx      int
}

type cacheData struct {
	iteration int
	height    int
}

func findChamberHeight(rockCount int, jetPattern []rune) int {
	// Init the state.
	s := state{
		jetPattern: jetPattern,
		cache:      make(map[key]cacheData),
		chamber:    [][CHAMBERSIZE]rune{},
	}

	// Drop rocks.
	for s.iteration = 0; s.iteration < rockCount; s.iteration++ {
		// Check if we already have the result in the cache.
		k := key{
			rockTypeIdx: s.rockTypeIdx,
			jetIdx:      s.jetIdx,
		}

		if cacheData, ok := s.cache[k]; ok {
			// We already have the result, so we can skip the rest of the iteration.
			multiple := (rockCount - s.iteration) / (s.iteration - cacheData.iteration)
			remainder := (rockCount - s.iteration) % (s.iteration - cacheData.iteration)

			if remainder == 0 {
				return len(s.chamber) + multiple*(len(s.chamber)-cacheData.height)
			}
		}

		// Update the cache.
		s.cache[k] = cacheData{
			iteration: s.iteration,
			height:    len(s.chamber),
		}

		// Drop the next rock.
		s.dropRock()
	}

	return len(s.chamber)
}

func (s *state) dropRock() {
	// Find the rock type to drop.

	currentRock := rocksTypes[s.rockTypeIdx]

	// x, y will be the coordinates of the top left corner of the rock.
	x, y := 2, len(s.chamber)-1+len(currentRock)+3 // -1 is here because coordinates start at 0.

	// While the rock can still fall.
	for {
		// Apply the jet pattern.
		switch s.jetPattern[s.jetIdx] {
		case '<':
			if s.canRockMoveLeft(currentRock, x, y) {
				x--
			}
		case '>':
			if s.canRockMoveRight(currentRock, x, y) {
				x++
			}
		}

		// Move the jet.
		s.jetIdx = (s.jetIdx + 1) % len(s.jetPattern)

		if s.canRockFall(currentRock, x, y) {
			y--
		} else {
			// The rock can't fall anymore, so we add it to the chamber.
			s.addRockToChamber(currentRock, x, y)
			break
		}
	}

	s.rockTypeIdx = (s.rockTypeIdx + 1) % len(rocksTypes)
}

// Add a rock to the chamber.
func (s *state) addRockToChamber(rock rock, x, y int) {

	// Add the rock to the chamber.
	for i, row := range rock {
		// Make sure the chamber is big enough.
		for y-i >= len(s.chamber) {
			s.chamber = append(s.chamber, [CHAMBERSIZE]rune{})
		}

		for j, r := range row {
			if r == '#' {
				s.chamber[y-i][x+j] = '#' // y-i because y is the top left corner of the rock.
			}
		}
	}
}

// Check if a rock can move left.
func (s *state) canRockMoveLeft(rock rock, x, y int) bool {
	// If the rock is already at the left edge of the chamber, it can't move left.
	if x <= 0 {
		return false
	}

	for i, row := range rock {
		if y-i >= len(s.chamber) {
			continue
		}

		for j, r := range row {
			if r == '#' && s.chamber[y-i][x+j-1] == '#' {
				return false
			}
		}
	}

	return true
}

// Check if a rock can move right.
func (s *state) canRockMoveRight(rock rock, x, y int) bool {
	// If the rock is already at the right edge of the chamber, it can't move right.
	if x+len(rock[0]) >= CHAMBERSIZE {
		return false
	}

	for i, row := range rock {
		if y-i >= len(s.chamber) {
			continue
		}

		for j, r := range row {
			if r == '#' && s.chamber[y-i][x+j+1] == '#' {
				return false
			}
		}
	}

	return true
}

// Check if a rock can fall.
func (s *state) canRockFall(rock rock, x, y int) bool {
	// If the rock is already at the bottom of the chamber, it can't fall.
	if y-len(rock) < 0 {
		return false
	}

	for i, row := range rock {
		if y-i > len(s.chamber) {
			// The rock is above the chamber, so it can fall.
			continue
		}

		for j, r := range row {
			if r == '#' && s.chamber[y-i-1][x+j] == '#' {
				return false
			}
		}
	}

	return true
}

// Return the string representation of the chamber.
func (s *state) String() string {
	var b strings.Builder

	b.WriteRune('\n')

	for i := len(s.chamber) - 1; i >= 0; i-- {
		for _, r := range s.chamber[i] {
			if r == '#' {
				b.WriteRune(r)
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}

	return b.String()
}
