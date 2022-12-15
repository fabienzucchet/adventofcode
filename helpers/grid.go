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
