package fabienz

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	floor := parseLines(lines)

	steps := floor.iterate()

	_, err = fmt.Fprintf(answer, "%d", steps)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type OceanFloor [][]rune

// Stringer for OceanFloor
func (f OceanFloor) String() string {
	toPrint := "\n"

	for _, line := range f {
		for _, char := range line {
			toPrint += string(char)
		}
		toPrint += "\n"
	}

	return toPrint
}

// INPUT PARSING

func parseLines(lines []string) (floor OceanFloor) {

	floor = make(OceanFloor, 0, len(lines))
	for _, line := range lines {
		floorRow := make([]rune, 0, len(line))
		for _, char := range line {
			floorRow = append(floorRow, char)
		}
		floor = append(floor, floorRow)
	}

	return floor
}

// LOGIC FUNCTIONS

// Execute movement for one hord
func (f *OceanFloor) moveHord(hord rune) (countMoved int) {

	// Create a copy of the floor
	newF := make(OceanFloor, len(*f))
	for row, line := range *f {
		floorRow := make([]rune, 0, len(line))
		floorRow = append(floorRow, line...)
		newF[row] = floorRow
	}

	// Use a channel for parallelism
	rows := make(chan int, len(*f))
	for row := range *f {
		rows <- row
	}
	close(rows)

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(hord rune) {
			for row := range rows {
				for col := range (*f)[row] {
					// If the cucumber is part of the hord and can move
					if newF[row][col] == hord && newF.canCucumberMove(row, col) {
						countMoved++
						switch hord {
						case '>':
							(*f)[row][(col+1)%len((*f)[row])] = '>'
							(*f)[row][col] = '.'

						case 'v':
							(*f)[(row+1)%len(*f)][col] = 'v'
							(*f)[row][col] = '.'
						}
					}
				}
			}
			wg.Done()
		}(hord)
	}
	wg.Wait()

	return countMoved
}

// Check of a cucumber can move
func (f OceanFloor) canCucumberMove(row, col int) bool {
	switch f[row][col] {
	case '>':
		return f[row][(col+1)%len(f[row])] == '.'
	case 'v':
		return f[(row+1)%len(f)][col] == '.'
	}

	return true
}

// Move hords alternatively until cucumbers can no longer move
func (f *OceanFloor) iterate() (step int) {

	hords := []rune{'>', 'v'}

	count := 1<<63 - 1

	for count > 0 {
		count = 0
		for _, hord := range hords {
			count += f.moveHord(hord)
		}
		step++
	}

	return step
}
