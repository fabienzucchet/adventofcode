package helpers

import (
	"fmt"
	"math"
)

// System of 2D coordinates.
type Coord2D struct {
	X int
	Y int
}

// System of 3D coordinates.
type Coord3D struct {
	X int
	Y int
	Z int
}

// Create the coordinates of all the neighbors of a given coordinate.
func (c Coord3D) NeighborsWithoutDiagonals() []Coord3D {
	return []Coord3D{
		{X: c.X - 1, Y: c.Y, Z: c.Z},
		{X: c.X + 1, Y: c.Y, Z: c.Z},
		{X: c.X, Y: c.Y - 1, Z: c.Z},
		{X: c.X, Y: c.Y + 1, Z: c.Z},
		{X: c.X, Y: c.Y, Z: c.Z - 1},
		{X: c.X, Y: c.Y, Z: c.Z + 1},
	}
}

// Manhattan distance between two coordinates.
func (c Coord2D) ManhattanDistance(other Coord2D) int {
	return AbsInt(c.X-other.X) + AbsInt(c.Y-other.Y)
}

// A 2D grid of int.
type IntGrid2D map[Coord2D]int

func (g IntGrid2D) getDimensions() (int, int, int, int) {
	minX := math.MaxInt
	maxX := math.MinInt
	minY := math.MaxInt
	maxY := math.MinInt
	for c := range g {
		if c.X < minX {
			minX = c.X
		}
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y < minY {
			minY = c.Y
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	}
	return minX, maxX, minY, maxY
}

// Print the grid
func (g IntGrid2D) Print() {
	minX, maxX, minY, maxY := g.getDimensions()
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			fmt.Printf("%d ", g[Coord2D{X: x, Y: y}])
		}
		fmt.Println()
	}
}
