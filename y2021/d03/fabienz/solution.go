package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

const DIAGSIZE = 12

// PartOne solves the first problem of day 3 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	var gamma [DIAGSIZE]int
	var epsilon [DIAGSIZE]int

	bytesCount := countBytesOne(lines)

	for i := range bytesCount {
		if bytesCount[i] > len(lines)/2 {
			gamma[i] = 1
		} else {
			epsilon[i] = 1
		}
	}

	_, err = fmt.Fprintf(answer, "%d", toDecimal(gamma)*toDecimal(epsilon))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	candidates := make([]string, len(lines))

	copy(candidates, lines)
	oxygen, err := filterOxValues(candidates, 0)
	if err != nil {
		return fmt.Errorf("error finding the correct values of oxygen and co2: %w", err)
	}
	oxygenValue, err := strconv.ParseInt(oxygen, 2, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse oxygen value %s : %w", oxygen, err)
	}

	copy(candidates, lines)
	co2, err := filterCo2Values(candidates, 0)
	if err != nil {
		return fmt.Errorf("error finding the correct values of oxygen and co2: %w", err)
	}
	co2Value, err := strconv.ParseInt(co2, 2, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse CO2 value %s : %w", co2, err)
	}

	_, err = fmt.Fprintf(answer, "%d", oxygenValue*co2Value)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Counts the number of bytes to 1 for every column
func countBytesOne(diagnostics []string) (count [DIAGSIZE]int) {

	for _, diagnostic := range diagnostics {
		for i, b := range diagnostic {
			if b == '1' {
				count[i]++
			}
		}
	}

	return count
}

func toDecimal(bin [DIAGSIZE]int) (dec int) {

	for i := range bin {
		dec += 1 << (len(bin) - i - 1) * bin[i]
	}

	return dec
}

func filterOxValues(candidates []string, offset int) (value string, err error) {

	// We stopped when we have only one candidate
	if len(candidates) == 1 {
		return candidates[0], nil
	}

	// Return an error if no valid candidate can be found
	if len(candidates) == 0 || offset >= len(candidates[0]) {
		return "", fmt.Errorf("no valid candidates were found")
	}

	// If we have more than one candidate, we do the filtering before the recursive call
	var filteredCandidates []string
	bytesCount := countBytesOne(candidates)
	for _, candidate := range candidates {
		switch {
		case bytesCount[offset] > len(candidates)-bytesCount[offset] && candidate[offset] == '1':
			filteredCandidates = append(filteredCandidates, candidate)
		case bytesCount[offset] == len(candidates)-bytesCount[offset] && candidate[offset] == '1':
			filteredCandidates = append(filteredCandidates, candidate)
		case bytesCount[offset] < len(candidates)-bytesCount[offset] && candidate[offset] == '0':
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	return filterOxValues(filteredCandidates, offset+1)

}

func filterCo2Values(candidates []string, offset int) (value string, err error) {

	// We stopped when we have only one candidate
	if len(candidates) == 1 {
		return candidates[0], nil
	}

	// Return an error if no valid candidate can be found
	if len(candidates) == 0 || offset >= len(candidates[0]) {
		return "", fmt.Errorf("no valid candidates were found")
	}

	// If we have more than one candidate, we do the filtering before the recursive call
	var filteredCandidates []string
	bytesCount := countBytesOne(candidates)
	for _, candidate := range candidates {
		switch {
		case bytesCount[offset] > len(candidates)-bytesCount[offset] && candidate[offset] == '0':
			filteredCandidates = append(filteredCandidates, candidate)
		case bytesCount[offset] == len(candidates)-bytesCount[offset] && candidate[offset] == '0':
			filteredCandidates = append(filteredCandidates, candidate)
		case bytesCount[offset] < len(candidates)-bytesCount[offset] && candidate[offset] == '1':
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	return filterCo2Values(filteredCandidates, offset+1)

}
