package fabienz

import (
	"fmt"
	"io"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 8 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	_, outputs, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", countUniqueNumbers(outputs))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	inputs, outputs, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	sum := 0

	for i := range inputs {
		sum += decode(inputs[i], outputs[i])
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PARSE INPUT
func parseLines(lines []string) (inputs [][]string, outputs [][]string, err error) {

	for _, line := range lines {
		splittedLines := strings.Split(line, "|")
		if len(splittedLines) != 2 {
			return nil, nil, fmt.Errorf("error splitting line %s", line)
		}

		inputs = append(inputs, strings.Fields(splittedLines[0]))
		outputs = append(outputs, strings.Fields(splittedLines[1]))
	}

	return inputs, outputs, nil
}

// PROCESSING FUNCTIONS

/*
* The number of segments per number is :
* 0 -> 6 -> abcefg
* 1 -> 2 -> cf
* 2 -> 5 -> acdeg
* 3 -> 5 -> acdfg
* 4 -> 4 -> bcdf
* 5 -> 5 -> abdfg
* 6 -> 6 -> abdefg
* 7 -> 3 -> acf
* 8 -> 7 -> abcdefg
* 9 -> 6 -> abcdfg
*
* To solve the puzzle we need to build a map mapping the observed segments with the corresponding one in the right order
 */

// Count the numbers with a unique number of segment
func countUniqueNumbers(outputs [][]string) (count int) {

	for _, line := range outputs {
		for _, output := range line {
			// If its a 1, 4 , 7 or 8
			if len(output) == 2 || len(output) == 3 || len(output) == 4 || len(output) == 7 {
				count++
			}
		}
	}

	return count
}

// Decode the numbers
func decode(inputs []string, outputs []string) (number int) {

	// Classify the inputs per length
	inputsPerLength := make(map[int][]string)

	for _, input := range inputs {
		inputsPerLength[len(input)] = append(inputsPerLength[len(input)], input)
	}

	// Mapping from theoretical segment to the experimental segment
	segmentsMapping := make(map[rune]rune)

	// Find the segments of 1 (len == 2)
	one := inputsPerLength[2][0]

	if isSegmentPresent(rune(one[0]), inputsPerLength[6]) {
		segmentsMapping[rune(one[0])] = 'f'
		segmentsMapping[rune(one[1])] = 'c'
	} else {
		segmentsMapping[rune(one[1])] = 'f'
		segmentsMapping[rune(one[0])] = 'c'
	}

	// Use the number 7 (len == 3)
	seven := inputsPerLength[3][0]

	for _, segment := range seven {
		if !isSegmentPresent(segment, inputsPerLength[2]) {
			segmentsMapping[segment] = 'a'
		}
	}

	// Use the number 4 (len == 4)
	four := inputsPerLength[4][0]

	for _, segment := range four {
		if !isSegmentPresent(segment, inputsPerLength[2]) {
			if isSegmentPresent(segment, inputsPerLength[5]) {
				segmentsMapping[segment] = 'd'
			} else {
				segmentsMapping[segment] = 'b'
			}
		}
	}

	eight := inputsPerLength[7][0]

	for _, segment := range eight {
		if isSegmentPresent(segment, inputsPerLength[5]) && !isSegmentPresent(segment, inputsPerLength[3]) && isSegmentPresent(segment, inputsPerLength[6]) {
			segmentsMapping[segment] = 'g'
		}
	}

	// The last one is e
	for _, segment := range eight {
		if _, present := segmentsMapping[segment]; !present {
			segmentsMapping[segment] = 'e'
		}
	}

	// Invert the map
	invertedSegmentsMapping := make(map[rune]rune)

	for experimental, theoretical := range segmentsMapping {
		invertedSegmentsMapping[theoretical] = experimental
	}

	for _, output := range outputs {
		number = number*10 + decodeNumber(output, invertedSegmentsMapping)
	}

	return number
}

// Checks if a segment is present in all elements of a list of string
func isSegmentPresent(segment rune, inputs []string) bool {
	for _, input := range inputs {
		isPresent := false
		for _, r := range input {
			if segment == r {
				isPresent = true
			}
		}

		if !isPresent {
			return false
		}
	}

	return true
}

// Convert to a number
func decodeNumber(input string, segmentsMapping map[rune]rune) int {
	switch len(input) {
	case 2:
		return 1
	case 3:
		return 7
	case 4:
		return 4
	case 5:
		if !isSegmentPresent(segmentsMapping['c'], []string{input}) {
			return 5
		} else if !isSegmentPresent(segmentsMapping['e'], []string{input}) {
			return 3
		} else {
			return 2
		}
	case 6:
		if !isSegmentPresent(segmentsMapping['c'], []string{input}) {
			return 6
		} else if !isSegmentPresent(segmentsMapping['d'], []string{input}) {
			return 0
		} else {
			return 9
		}
	case 7:
		return 8
	}

	return 0
}

func printSegmentMapping(segmentMapping map[rune]rune) string {
	var toPrint string

	for key, val := range segmentMapping {
		toPrint += " " + string(key) + " -> " + string(val) + ","
	}

	return toPrint
}
