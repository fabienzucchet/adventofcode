package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	x, y, err := locationIdsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not get location ids from lines: %w", err)
	}

	totalDistance := getTotalDistance(x, y)

	_, err = fmt.Fprintf(answer, "%d", totalDistance)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	left, right, err := locationIdsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not get location ids from lines: %w", err)
	}

	totalSimilarityScore := getTotalSimilarityScore(left, right)

	_, err = fmt.Fprintf(answer, "%d", totalSimilarityScore)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func locationIdsFromLines(lines []string) ([]int, []int, error) {
	x := make([]int, 0, len(lines))
	y := make([]int, 0, len(lines))

	for _, line := range lines {
		xLocation, yLocation, err := locationIdsFromString(line)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse location ids from line: %w", err)
		}

		x = append(x, xLocation)
		y = append(y, yLocation)
	}

	sortLocationIds(x)
	sortLocationIds(y)

	return x, y, nil
}

func locationIdsFromString(s string) (int, int, error) {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(s, -1)
	if len(matches) != 2 {
		return 0, 0, fmt.Errorf("could not parse location ids from string: %s", s)
	}

	x, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse x location id: %w", err)
	}

	y, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse y location id: %w", err)
	}

	return x, y, nil
}

func sortLocationIds(x []int) {
	for i := 0; i < len(x)-1; i++ {
		for j := i + 1; j < len(x); j++ {
			if x[i] > x[j] {
				x[i], x[j] = x[j], x[i]
			}
		}
	}
}

func getDistance(x, y int) int {
	if x < y {
		return y - x
	} else {
		return x - y
	}
}

func getTotalDistance(x []int, y []int) int {
	totalDistance := 0
	for i := 0; i < len(x); i++ {
		totalDistance += getDistance(x[i], y[i])
	}
	return totalDistance
}

func getSimilarityScore(leftLocationId int, rightLocationIds []int) int {
	occurrenceCount := 0
	for _, rightLocationId := range rightLocationIds {
		if leftLocationId == rightLocationId {
			occurrenceCount++
		}
	}

	return leftLocationId * occurrenceCount
}

func getTotalSimilarityScore(x []int, y []int) int {
	totalSimilarityScore := 0
	for _, xLocationId := range x {
		totalSimilarityScore += getSimilarityScore(xLocationId, y)
	}

	return totalSimilarityScore
}
