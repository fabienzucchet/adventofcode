package fabienz

import (
	"fmt"
	"io"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0

	for _, line := range lines {
		items1, items2 := parseLine(line)
		commonItems := findCommonItems(items1, items2)
		sum += priorities[commonItems[0]]
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Group the lines by 3.
	groups := groupByThree(lines)

	sum := 0

	// Find the badge in each group.
	for _, group := range groups {
		badge := findCommonCharacter(group)
		sum += priorities[badge]
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Parse the input to return the two half of a slice of string.
func parseLine(line string) ([]string, []string) {
	items := []string{}
	for _, c := range line {
		items = append(items, string(c))
	}

	return items[:len(items)/2], items[len(items)/2:]
}

// Find the common item between two slices of string.
func findCommonItems(items1, items2 []string) []string {
	commonItems := []string{}
	for _, item1 := range items1 {
		for _, item2 := range items2 {
			if item1 == item2 {
				commonItems = append(commonItems, item1)
			}
		}
	}

	return commonItems
}

// Map containing the priorities of each item type.
var priorities = map[string]int{
	"a": 1,
	"b": 2,
	"c": 3,
	"d": 4,
	"e": 5,
	"f": 6,
	"g": 7,
	"h": 8,
	"i": 9,
	"j": 10,
	"k": 11,
	"l": 12,
	"m": 13,
	"n": 14,
	"o": 15,
	"p": 16,
	"q": 17,
	"r": 18,
	"s": 19,
	"t": 20,
	"u": 21,
	"v": 22,
	"w": 23,
	"x": 24,
	"y": 25,
	"z": 26,
	"A": 27,
	"B": 28,
	"C": 29,
	"D": 30,
	"E": 31,
	"F": 32,
	"G": 33,
	"H": 34,
	"I": 35,
	"J": 36,
	"K": 37,
	"L": 38,
	"M": 39,
	"N": 40,
	"O": 41,
	"P": 42,
	"Q": 43,
	"R": 44,
	"S": 45,
	"T": 46,
	"U": 47,
	"V": 48,
	"W": 49,
	"X": 50,
	"Y": 51,
	"Z": 52,
}

// Group the lines by 3.
func groupByThree(lines []string) [][]string {
	groups := [][]string{}
	for i := 0; i < len(lines); i += 3 {
		groups = append(groups, lines[i:i+3])
	}
	return groups
}

// Find the character in common in a list of strings.
func findCommonCharacter(lines []string) string {
	for _, c := range lines[0] {
		found := true
		for _, line := range lines {
			if !strings.ContainsRune(line, c) {
				found = false
				break
			}
		}
		if found {
			return string(c)
		}
	}
	return ""
}
