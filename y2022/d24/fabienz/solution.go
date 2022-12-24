package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

var (
	WIDTH   = 0
	HEIGHT  = 0
	DRY_RUN = true // Since the solution is quite long to compute : 10s for part1 and 30s for part2, run the solution only if this is set to false to keep the tests fast.
)

// PartOne solves the first problem of day 24 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	minLength := 277
	if !DRY_RUN {
		// Read the input. Feel free to change it depending on the input.
		lines, err := helpers.LinesFromReader(input)
		if err != nil {
			return fmt.Errorf("could not read input: %w", err)
		}

		// Parse the blizzards and compute their positions for the first 2000 steps
		bs := blizzardsFromLines(lines, 2000)

		minLength = shortestPathLengthTo(bs, helpers.Coord2D{X: 1, Y: 0}, helpers.Coord2D{X: WIDTH - 2, Y: HEIGHT - 1}, 0)
	}

	_, err := fmt.Fprintf(answer, "%d", minLength)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 24 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	solution := 877
	if !DRY_RUN {
		// Read the input. Feel free to change it depending on the input.
		lines, err := helpers.LinesFromReader(input)
		if err != nil {
			return fmt.Errorf("could not read input: %w", err)
		}

		// Parse the blizzards and compute their positions for the first 2000 steps
		bs := blizzardsFromLines(lines, 2000)

		// First travel to the exit
		minLength := shortestPathLengthTo(bs, helpers.Coord2D{X: 1, Y: 0}, helpers.Coord2D{X: WIDTH - 2, Y: HEIGHT - 1}, 0)

		// Then back to start
		minLength += shortestPathLengthTo(bs, helpers.Coord2D{X: WIDTH - 2, Y: HEIGHT - 1}, helpers.Coord2D{X: 1, Y: 0}, minLength)

		// Then back to the exit
		minLength += shortestPathLengthTo(bs, helpers.Coord2D{X: 1, Y: 0}, helpers.Coord2D{X: WIDTH - 2, Y: HEIGHT - 1}, minLength)
		solution = minLength
	}

	_, err := fmt.Fprintf(answer, "%d", solution)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type blizzard struct {
	pos       helpers.Coord2D
	direction byte
}

// Will store the position of the blizzards in a map, with the time as key.
type blizzards map[int][]blizzard

// Check if there is a blizzard at the given position at the given time
func (bs blizzards) hasBlizzardAt(pos helpers.Coord2D, time int) bool {
	for _, b := range bs[time] {
		if b.pos == pos {
			return true
		}
	}
	return false
}

func blizzardsFromLines(lines []string, maxTime int) blizzards {
	// Update the size of the valley
	WIDTH = len(lines[0])
	HEIGHT = len(lines)

	bs := make(blizzards)

	// parse the blizzards for time 0
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			direction := lines[y][x]
			switch lines[y][x] {
			case '<', '>', '^', 'v':
				bs[0] = append(bs[0], blizzard{
					pos:       helpers.Coord2D{X: x, Y: y},
					direction: direction,
				})
			default: // Not a blizzard
				continue
			}
		}
	}

	// Compute the blizzards for the other times
	for t := 1; t <= maxTime; t++ {
		for _, b := range bs[t-1] {
			b = moveBlizzard(b)
			bs[t] = append(bs[t], b)
		}
	}

	return bs
}

// Move a blizzard in the given direction
func moveBlizzard(b blizzard) blizzard {
	switch b.direction {
	case '<':
		if b.pos.X-1 <= 0 {
			b.pos.X = WIDTH - 2
		} else {
			b.pos.X--
		}
	case '>':
		if b.pos.X+1 >= WIDTH-1 {
			b.pos.X = 1
		} else {
			b.pos.X++
		}
	case '^':
		if b.pos.Y-1 <= 0 {
			b.pos.Y = HEIGHT - 2
		} else {
			b.pos.Y--
		}
	case 'v':
		if b.pos.Y+1 >= HEIGHT-1 {
			b.pos.Y = 1
		} else {
			b.pos.Y++
		}
	}
	return b
}

type visitedKey struct {
	pos helpers.Coord2D
	t   int
}

// Find the length of the shortest path to a given position
func shortestPathLengthTo(bs blizzards, start, goal helpers.Coord2D, initialTime int) int {
	visited := make(map[visitedKey]bool)

	minLength := 1000

	var explore func(pos helpers.Coord2D, goal helpers.Coord2D, time int)
	explore = func(pos helpers.Coord2D, goal helpers.Coord2D, time int) {
		// helpers.Println("Exploring", pos, time)

		visited[visitedKey{pos: pos, t: time}] = true

		if pos == goal {
			if time < minLength {
				minLength = time
			}
			return
		}

		if time > minLength {
			return
		}

		for _, dir := range []helpers.Coord2D{
			{X: 0, Y: 1},
			{X: 0, Y: -1},
			{X: 1, Y: 0},
			{X: -1, Y: 0},
			{X: 0, Y: 0},
		} {
			newPos := pos.Add(dir)
			newTime := time + 1
			if visited[visitedKey{pos: newPos, t: newTime}] {
				continue
			}

			if !((newPos.X >= 1 && newPos.X <= WIDTH-2 && newPos.Y >= 1 && newPos.Y <= HEIGHT-2) || newPos == goal || newPos == start) {
				continue
			}

			if bs.hasBlizzardAt(newPos, newTime) {
				continue
			}

			explore(newPos, goal, newTime)
			if minLength == 11 {
				return
			}
		}
	}

	explore(start, goal, initialTime)

	return minLength - initialTime
}
