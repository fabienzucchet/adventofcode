package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	decimals := snafuToDecimalSlice(lines)

	sumInDecimal := 0

	for _, decimal := range decimals {
		sumInDecimal += decimal
	}

	_, err = fmt.Fprintf(answer, "%s", decimalToSnafu(sumInDecimal))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Map to convert SNAFU digits to decimal digits.
var snafuToDecimalMap = map[byte]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

var decimalToSnafuMap = map[int]byte{
	-2: '=',
	-1: '-',
	0:  '0',
	1:  '1',
	2:  '2',
}

// Convert SNAFU to decimal.
func snafuToDecimal(snafu string) int {

	if snafu == "" {
		return 0
	}

	return 5*snafuToDecimal(snafu[:len(snafu)-1]) + snafuToDecimalMap[snafu[len(snafu)-1]]
}

// Convert a slice of SNAFU to a slice of decimal.
func snafuToDecimalSlice(snafu []string) []int {
	decimal := make([]int, len(snafu))

	for i := 0; i < len(snafu); i++ {
		decimal[i] = snafuToDecimal(snafu[i])
	}

	return decimal
}

// Convert decimal to SNAFU.
func decimalToSnafu(decimal int) string {

	// No negative numbers in SNAFU.
	if decimal > 0 {
		// Offset 2 to have the remainder in the range [0, 5] instead of [-2, 2].
		q, r := (decimal+2)/5, (decimal+2)%5

		// Do not forget to remove the offset to have the remainder back in the range [-2, 2].
		return decimalToSnafu(q) + string(decimalToSnafuMap[r-2])
	}

	return ""
}
