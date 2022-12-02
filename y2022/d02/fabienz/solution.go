package fabienz

import (
	"fmt"
	"io"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	score := 0

	for _, line := range lines {
		play := strings.Split(line, " ")
		score += computeScore(play[1], play[0])
	}

	_, err = fmt.Fprintf(answer, "%d", score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	score := 0

	for _, line := range lines {
		play := strings.Split(line, " ")
		figureToPlay := findFigureToPlay(play[1], play[0])
		score += computeScore(figureToPlay, play[0])
	}

	_, err = fmt.Fprintf(answer, "%d", score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Compute the score for a round
func computeScore(you string, opponent string) int {

	idxYou := int(you[0]) - int('X')
	idxOpponent := int(opponent[0]) - int('A')

	// The score of the figure can be computed by substracting 'X' to the figure letter

	// We can check if we win or loose by comparing the figure letter to the opponent letter (normalized to get a number between 0 and 2).
	// If the numbers are the same, it's a draw.
	// If the numbers are different, we can check if the difference is 1 or 2 (mod 3). If it's 1, we win, if it's 2, we loose.

	// Don't forget to add 1 to the figure score since the figure score is between 1 and 3.

	switch (idxYou - idxOpponent + 3) % 3 {
	case 0:
		// Draw
		return 3 + idxYou + 1
	case 1:
		// Win
		return 6 + idxYou + 1
	case 2:
		// Loose
		return idxYou + 1
	}

	return 0
}

// Find the figure to play in part 2
func findFigureToPlay(result string, opponent string) string {
	// Compute an offset depending on the outcome we want.
	// If we want to win, the offset is 1, if we want to loose, the offset is -1. If we want to draw, the offset is 0.
	offset := int(result[0]) - int('Y')

	// Compute the index of the figure to play by adding the offset to the opponent figure index.
	idx := (int(opponent[0]) - int('A') + offset + 3) % 3

	return string(rune(idx + int('X')))
}
