package fabienz

import (
	"fmt"
	"io"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	for _, group := range parseLines(lines) {
		count += countGroupAnswersWrongInstructions(group)
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	count := 0

	for _, group := range parseLines(lines) {
		count += countGroupAnswersCorrectInstructions(group)
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseLines(lines []string) []string {
	parsedLines := make([]string, 1)

	for _, line := range lines {
		if line == "" {
			parsedLines = append(parsedLines, "")
		} else if parsedLines[len(parsedLines)-1] == "" {
			parsedLines[len(parsedLines)-1] = parsedLines[len(parsedLines)-1] + line
		} else {
			parsedLines[len(parsedLines)-1] = parsedLines[len(parsedLines)-1] + " " + line
		}
	}

	return parsedLines
}

func countGroupAnswersWrongInstructions(answers string) int {
	distinctAnswers := make(map[rune]int)

	for _, answer := range answers {
		_, ok := distinctAnswers[answer]
		if answer != ' ' && !ok {
			distinctAnswers[answer] = 1
		}
	}

	return len(distinctAnswers)
}

func countGroupAnswersCorrectInstructions(answers string) int {
	groupSize := len(strings.Split(answers, " "))

	answersCount := make(map[rune]int)

	for _, answer := range answers {
		_, ok := answersCount[answer]
		if answer != ' ' && !ok {
			answersCount[answer] = 1
		} else if answer != ' ' {
			answersCount[answer] = answersCount[answer] + 1
		}
	}

	count := 0

	for _, nb := range answersCount {
		if nb == groupSize {
			count = count + 1
		}
	}

	return count
}
