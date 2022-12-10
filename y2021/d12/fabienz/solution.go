package fabienz

import (
	"fmt"
	"io"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := parseLines(lines)

	paths := g.findPaths(false)

	_, err = fmt.Fprintf(answer, "%d", len(paths))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g := parseLines(lines)

	paths := g.findPaths(true)

	_, err = fmt.Fprintf(answer, "%d", len(paths))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Graph map[string][]string

// INPUT PARSING

func parseLines(lines []string) (g Graph) {

	g = make(map[string][]string)

	for _, line := range lines {
		vertices := strings.SplitN(line, "-", 2)
		if len(vertices) == 2 {
			g[vertices[0]] = append(g[vertices[0]], vertices[1])
			g[vertices[1]] = append(g[vertices[1]], vertices[0])
		}
	}

	return g
}

// Use a graph traversal to find the paths
func (g *Graph) findPaths(allowTwoVisits bool) (paths [][]string) {

	// Track the visited caves
	var visited []string

	var rec func(x string, visited []string) (paths [][]string)

	// Recursive helper finding paths starting from vertice x
	rec = func(x string, visited []string) (paths [][]string) {

		// If we reached the end, we stop
		if x == "end" {
			paths = append(paths, []string{})
			return paths
		}

		var accessibleNeighbors []string

		if !allowTwoVisits {
			accessibleNeighbors = g.getAccessibleNeighbors(x, visited)
		} else {
			accessibleNeighbors = g.getAccessibleNeighborsTwoVisits(x, visited, hadAlreadyTwoVisits(visited))
		}

		// If we are in a dead end, we stop
		if len(accessibleNeighbors) == 0 {
			return paths
		}

		// Else find the paths recursively
		for _, neigh := range accessibleNeighbors {
			subPaths := rec(neigh, append(visited, neigh))

			for _, p := range subPaths {
				p = append(p, neigh)
				paths = append(paths, p)
			}
		}

		return paths
	}

	paths = rec("start", append(visited, "start"))

	return paths
}

// Check if a string is lower case
func isLower(x string) bool {
	return x == strings.ToLower(x)
}

// Check if an element is already in the path
func isInPath(x string, p []string) bool {

	for _, v := range p {
		if x == v {
			return true
		}
	}

	return false
}

// Find the accessible neighbors
func (g *Graph) getAccessibleNeighbors(x string, visited []string) (neighbors []string) {

	for _, neigh := range (*g)[x] {
		if !isLower(neigh) || !isInPath(neigh, visited) {
			neighbors = append(neighbors, neigh)
		}
	}

	return neighbors
}

// Find the accessible neighbors with possibility to visit ONE small cave once
func (g *Graph) getAccessibleNeighborsTwoVisits(x string, visited []string, hadTwoVisits bool) (neighbors []string) {

	for _, neigh := range (*g)[x] {
		if !isLower(neigh) || !isInPath(neigh, visited) || (neigh != "start" && neigh != "end" && !hadTwoVisits) {
			neighbors = append(neighbors, neigh)
		}
	}

	return neighbors
}

func hadAlreadyTwoVisits(p []string) bool {
	if len(p) < 2 {
		return false
	}

	for idx, v := range p[:len(p)-2] {
		// If we find a small cave with two occurence in p -> true
		if isLower(v) && isInPath(v, p[idx+1:]) {
			return true
		}
	}

	return false
}
