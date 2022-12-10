package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	highest := -1

	for _, boardingPass := range lines {
		seatId, err := computeSeatId(boardingPass)
		if err != nil {
			return fmt.Errorf("error finding seat ID : %w", err)
		}

		if seatId > highest {
			highest = seatId
		}
	}

	_, err = fmt.Fprintf(answer, "%d", highest)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var usedSeats [864]int
	var mySeat int

	for _, boardingPass := range lines {
		seatId, err := computeSeatId(boardingPass)
		if err != nil {
			return fmt.Errorf("error finding seat ID : %w", err)
		}
		usedSeats[seatId-1] = 1
	}

	for idx := 1; idx < 864; idx++ {
		if usedSeats[idx] == 0 && usedSeats[idx-1] == 1 && usedSeats[idx+1] == 1 {
			mySeat = idx + 1
			break
		}
	}

	_, err = fmt.Fprintf(answer, "%d", mySeat)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func computeSeatId(boarding_pass string) (int, error) {
	binaryString := strings.ReplaceAll(boarding_pass, "F", "0")
	binaryString = strings.ReplaceAll(binaryString, "B", "1")
	binaryString = strings.ReplaceAll(binaryString, "L", "0")
	binaryString = strings.ReplaceAll(binaryString, "R", "1")
	row, err := strconv.ParseInt(binaryString[:7], 2, 8)
	if err != nil {
		return 0, fmt.Errorf("error parsing row %s : %w", binaryString[:7], err)
	}

	column, err := strconv.ParseInt(binaryString[7:], 2, 4)
	if err != nil {
		return 0, fmt.Errorf("error parsing column %s : %w", binaryString[7:], err)
	}

	return int(row*8 + column), nil
}
