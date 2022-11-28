package fabienz

import (
	"fmt"
	"io"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	grid := lines
	newgrid := iterate(grid)

	for hasGridChanged(grid, newgrid) {
		grid = copySlice(newgrid)
		newgrid = iterate(grid)
	}

	_, err = fmt.Fprintf(answer, "%d", countOccupied(newgrid))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	grid := lines
	newgrid := iterate2(grid)

	for hasGridChanged(grid, newgrid) {
		grid = copySlice(newgrid)
		newgrid = iterate2(grid)
	}

	_, err = fmt.Fprintf(answer, "%d", countOccupied(newgrid))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func copySlice(s []string) []string {
	var newS []string

	// Copy the initial grid
	for _, row := range s {
		newS = append(newS, row)
	}

	return newS
}

func iterate(grid []string) []string {

	newgrid := copySlice(grid)

	nrow := len(grid)
	ncol := len(grid[0])

	for i := 0; i < nrow; i++ {
		for j := 0; j < ncol; j++ {

			noccupied := 0

			if i > 0 {
				if grid[i-1][j] == '#' {
					noccupied++
				}

				if j > 0 {
					if grid[i-1][j-1] == '#' {
						noccupied++
					}
				}

				if j < ncol-1 {
					if grid[i-1][j+1] == '#' {
						noccupied++
					}
				}
			}

			if i < nrow-1 {
				if grid[i+1][j] == '#' {
					noccupied++
				}

				if j > 0 {
					if grid[i+1][j-1] == '#' {
						noccupied++
					}
				}

				if j < ncol-1 {
					if grid[i+1][j+1] == '#' {
						noccupied++
					}
				}
			}

			if j > 0 {
				if grid[i][j-1] == '#' {
					noccupied++
				}
			}

			if j < ncol-1 {
				if grid[i][j+1] == '#' {
					noccupied++
				}
			}

			// Update the new state of the grid
			if grid[i][j] == 'L' && noccupied == 0 {
				newgrid[i] = replaceAtIndex(newgrid[i], '#', j)
			} else if grid[i][j] == '#' && noccupied >= 4 {
				newgrid[i] = replaceAtIndex(newgrid[i], 'L', j)
			}
		}
	}

	return newgrid
}

func iterate2(grid []string) []string {

	newgrid := copySlice(grid)

	nrow := len(grid)
	ncol := len(grid[0])

	for i := 0; i < nrow; i++ {
		for j := 0; j < ncol; j++ {

			noccupied := 0

			k := 0
			l := 0

			// Check up
			k = i - 1
			for k > 0 && grid[k][j] == '.' {
				k--
			}

			if k >= 0 && grid[k][j] == '#' {
				noccupied++
			}

			// Check left
			k = j - 1
			for k > 0 && grid[i][k] == '.' {
				k--
			}

			if k >= 0 && grid[i][k] == '#' {
				noccupied++
			}

			// Check bottom
			k = i + 1
			for k < nrow && grid[k][j] == '.' {
				k++
			}

			if k < nrow && grid[k][j] == '#' {
				noccupied++
			}

			// Check right
			k = j + 1
			for k < ncol && grid[i][k] == '.' {
				k++
			}

			if k < ncol && grid[i][k] == '#' {
				noccupied++
			}

			// Check top-left
			k = i - 1
			l = j - 1
			for k > 0 && l > 0 && grid[k][l] == '.' {
				k--
				l--

			}

			if k >= 0 && l >= 0 && grid[k][l] == '#' {
				noccupied++
			}

			// Check top-right
			k = i - 1
			l = j + 1
			for k > 0 && l < ncol && grid[k][l] == '.' {
				k--
				l++

			}

			if k >= 0 && l < ncol && grid[k][l] == '#' {
				noccupied++
			}

			// Check bottom-right
			k = i + 1
			l = j + 1
			for k < nrow && l < ncol && grid[k][l] == '.' {
				k++
				l++

			}

			if k < nrow && l < ncol && grid[k][l] == '#' {
				noccupied++
			}

			// Check bottom-left
			k = i + 1
			l = j - 1
			for k < nrow && l > 0 && grid[k][l] == '.' {
				k++
				l--

			}

			if k < nrow && l >= 0 && grid[k][l] == '#' {
				noccupied++
			}

			// Update the new state of the grid
			if grid[i][j] == 'L' && noccupied == 0 {
				newgrid[i] = replaceAtIndex(newgrid[i], '#', j)
			} else if grid[i][j] == '#' && noccupied >= 5 {
				newgrid[i] = replaceAtIndex(newgrid[i], 'L', j)
			}
		}
	}

	return newgrid
}

func countOccupied(grid []string) int {
	count := 0

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == '#' {
				count++
			}
		}
	}

	return count
}

func hasGridChanged(grid []string, newgrid []string) bool {
	res := false

	for i := 0; i < len(grid); i++ {
		if grid[i] != newgrid[i] {
			res = true
		}
	}

	return res
}
