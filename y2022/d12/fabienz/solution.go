package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input
	m := parseMap(lines)

	// Find the shortest path
	dist := m.findShortestPath(false)

	_, err = fmt.Fprintf(answer, "%d", dist[m.end])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input
	m := parseMap(lines)

	// Reverse end and start
	m.start, m.end = m.end, m.start

	// Find the shortest path
	dist := m.findShortestPath(true)

	// Find the minimum distance to the end
	minDist := 100000
	for c, d := range dist {
		if d < minDist && m.elevation[c] == 0 {
			minDist = d
		}
	}

	_, err = fmt.Fprintf(answer, "%d", minDist)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Coord struct {
	X int
	Y int
}

type Map struct {
	elevation map[Coord]int
	start     Coord
	end       Coord
	width     int
	height    int
}

func parseMap(lines []string) Map {
	m := Map{
		elevation: make(map[Coord]int),
		width:     len(lines[0]),
		height:    len(lines),
	}

	for y, line := range lines {
		for x, c := range line {
			switch c {
			case 'S':
				m.start = Coord{X: x, Y: y}
				m.elevation[Coord{X: x, Y: y}] = 0
			case 'E':
				m.end = Coord{X: x, Y: y}
				m.elevation[Coord{X: x, Y: y}] = 25
			default:
				m.elevation[Coord{X: x, Y: y}] = int(c - 'a')
			}
		}

	}

	return m
}

// Find the shortest path using Dijkstra's algorithm
func (m *Map) findShortestPath(reverse bool) map[Coord]int {
	dist := make(map[Coord]int)
	for c := range m.elevation {
		dist[c] = 10000
	}
	dist[m.start] = 0
	visited := make(map[Coord]bool)

	for len(visited) < len(m.elevation) {
		// Find the closest unvisited node
		var closest Coord
		minDist := 100000
		for c, d := range dist {
			if d < minDist && !visited[c] {
				minDist = d
				closest = c
			}
		}

		// Mark it as visited
		visited[closest] = true

		// Update the distance of its neighbors
		for _, n := range m.neighbors(closest, reverse) {
			if !visited[n] {
				d := dist[closest] + 1
				if d < dist[n] {
					dist[n] = d
				}
			}
		}
	}

	return dist
}

func (m *Map) neighbors(c Coord, reverse bool) []Coord {
	var neighbors []Coord

	if c.X > 0 && m.isAccessible(c, Coord{X: c.X - 1, Y: c.Y}, reverse) {
		neighbors = append(neighbors, Coord{X: c.X - 1, Y: c.Y})
	}
	if c.X < m.width-1 && m.isAccessible(c, Coord{X: c.X + 1, Y: c.Y}, reverse) {
		neighbors = append(neighbors, Coord{X: c.X + 1, Y: c.Y})
	}
	if c.Y > 0 && m.isAccessible(c, Coord{X: c.X, Y: c.Y - 1}, reverse) {
		neighbors = append(neighbors, Coord{X: c.X, Y: c.Y - 1})
	}
	if c.Y < m.height-1 && m.isAccessible(c, Coord{X: c.X, Y: c.Y + 1}, reverse) {
		neighbors = append(neighbors, Coord{X: c.X, Y: c.Y + 1})
	}

	return neighbors
}

func (m *Map) isAccessible(from, to Coord, reverse bool) bool {
	if reverse {
		return m.elevation[from]-m.elevation[to] <= 1
	}

	return m.elevation[to]-m.elevation[from] <= 1
}
