package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 1 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// This slice will store the number of calories for each elf.
	calories := make([]int, 0)

	currentCalories := 0

	// For each line, we extract the number of calories and add it to the correct elf.
	for _, line := range lines {
		// If empty, go to the next elf
		if line == "" {
			calories = append(calories, currentCalories)
			currentCalories = 0
			continue
		}

		// Extract the number of calories
		cal, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("could not parse %q: %w", line, err)
		}

		// Add it to the current elf
		currentCalories += cal
	}

	_, err = fmt.Fprintf(answer, "%d", max(calories))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 1 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// This slice will store the number of calories for each elf.
	calories := make([]int, 0)

	currentCalories := 0

	// For each line, we extract the number of calories and add it to the correct elf.
	for _, line := range lines {
		// If empty, go to the next elf
		if line == "" {
			calories = append(calories, currentCalories)
			currentCalories = 0
			continue
		}

		// Extract the number of calories
		cal, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("could not parse %q: %w", line, err)
		}

		// Add it to the current elf
		currentCalories += cal
	}

	// Find the 3 maximum
	max1, max2, max3 := max3(calories)

	_, err = fmt.Fprintf(answer, "%d", max1+max2+max3)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Returns the maximum of an array of integers.
func max(arr []int) int {
	max := 0

	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}

// Returns the 3 maximum of an array of integers.
func max3(arr []int) (int, int, int) {
	max1 := 0
	max2 := 0
	max3 := 0

	for _, v := range arr {
		if v > max1 {
			max3 = max2
			max2 = max1
			max1 = v
		} else if v > max2 {
			max3 = max2
			max2 = v
		} else if v > max3 {
			max3 = v
		}
	}

	return max1, max2, max3
}
