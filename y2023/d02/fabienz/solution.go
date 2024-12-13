package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2023.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	maxCubesInGame := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	sum := 0

	for _, line := range lines {
		gameId, maxCounts, err := parseGame(line)
		if err != nil {
			return fmt.Errorf("could not parse game: %w", err)
		}

		if isGameValid(maxCounts, maxCubesInGame) {
			sum += gameId
		}
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2023.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	sum := 0

	for _, line := range lines {
		_, maxCounts, err := parseGame(line)
		if err != nil {
			return fmt.Errorf("could not parse game: %w", err)
		}

		sum += powerOfACube(maxCounts)
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseGame(line string) (gameId int, maxCounts map[string]int, err error) {
	splitLine := strings.Split(line, ": ")
	if len(splitLine) != 2 {
		return 0, maxCounts, fmt.Errorf("invalid line: %s", line)
	}

	gameId, err = parseGameID(splitLine[0])
	if err != nil {
		return 0, maxCounts, fmt.Errorf("failed to parse game id: %w", err)
	}

	maxCounts, err = findMaxCubesByColor(strings.Split(splitLine[1], "; "))
	if err != nil {
		return 0, maxCounts, fmt.Errorf("failed to find max cubes by color: %w", err)
	}

	return gameId, maxCounts, nil
}

// Parse the game id from the prefix of the line.
func parseGameID(line string) (int, error) {
	id, err := strconv.Atoi((strings.Split(line, " "))[1])
	if err != nil {
		return 0, fmt.Errorf("invalid game id: %s", line)
	}

	return id, nil
}

// return the max number of cubes shown by color
func findMaxCubesByColor(subsets []string) (maxCounts map[string]int, err error) {
	maxCounts = make(map[string]int)

	for _, subset := range subsets {
		for _, cubes := range strings.Split(subset, ", ") {
			count, color, err := parseCubeCountAndColor(cubes)
			if err != nil {
				return maxCounts, fmt.Errorf("failed to parse cube count and color: %w", err)
			}

			if _, ok := maxCounts[color]; !ok {
				maxCounts[color] = count
			} else {
				if count > maxCounts[color] {
					maxCounts[color] = count
				}
			}
		}
	}

	return maxCounts, nil
}

func parseCubeCountAndColor(s string) (count int, color string, err error) {
	split := strings.Split(s, " ")
	if len(split) != 2 {
		return 0, "", fmt.Errorf("invalid cube count and color: %s", s)
	}

	count, err = strconv.Atoi(split[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid cube count: %s", s)
	}

	return count, split[1], nil
}

func isGameValid(maxCounts map[string]int, maxCubesInGame map[string]int) bool {
	for color, count := range maxCounts {
		if count > maxCubesInGame[color] {
			return false
		}
	}

	return true
}

func powerOfACube(maxCounts map[string]int) int {
	return maxCounts["red"] * maxCounts["green"] * maxCounts["blue"]
}
