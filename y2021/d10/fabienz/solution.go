package fabienz

import (
	"fmt"
	"io"
	"sort"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions := parseLines(lines)

	sum := 0

	for _, instruction := range instructions {
		illegalChar, isIllegal, _ := findIllegalCharacter(instruction)

		if isIllegal {
			sum += errorScoreMap[illegalChar]
		}
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	instructions := parseLines(lines)

	var scores []int

	for _, instruction := range instructions {
		_, isIllegal, remainder := findIllegalCharacter(instruction)

		if !isIllegal {
			scores = append(scores, completeInstruction(remainder))
		}
	}

	sort.Ints(scores)

	_, err = fmt.Fprintf(answer, "%d", scores[len(scores)/2])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Pile []rune

var errorScoreMap = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var completeScoreMap = map[rune]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

// INPUT PARSING
func parseLines(lines []string) (instructions [][]rune) {

	for _, line := range lines {
		var instruction []rune
		for _, char := range line {
			instruction = append(instruction, char)
		}
		instructions = append(instructions, instruction)
	}

	return instructions
}

// Find the illegal character
func findIllegalCharacter(instruction []rune) (illegalChar rune, isIllegal bool, remainder Pile) {

	var p Pile

	for _, char := range instruction {

		if char == '(' || char == '{' || char == '[' || char == '<' {
			p = p.push(char)
		} else {

			if len(p) == 0 {
				return illegalChar, false, p
			}

			var previous rune

			previous, p = p.pop()

			switch char {
			case ')':
				if previous != '(' {
					return char, true, p
				}
			case '}':
				if previous != '{' {
					return char, true, p
				}
			case ']':
				if previous != '[' {
					return char, true, p
				}
			case '>':
				if previous != '<' {
					return char, true, p
				}
			}
		}
	}

	return illegalChar, false, p
}

// Push on the pile
func (p Pile) push(char rune) Pile {
	p = append(p, char)

	return p
}

// Pop from the pile
func (p Pile) pop() (rune, Pile) {
	char := p[len(p)-1]

	return char, p[:len(p)-1]
}

// Complete an instruction and compute the autocomplete score
func completeInstruction(remainder Pile) (score int) {
	var char rune
	for len(remainder) > 0 {
		char, remainder = remainder.pop()
		score *= 5
		score += completeScoreMap[char]
	}

	return score
}
