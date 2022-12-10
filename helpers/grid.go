package helpers

// System of 2D coordinates.
type Coord2D struct {
	X int
	Y int
}

// A 2D grid of int.
type IntGrid2D map[Coord2D]int
