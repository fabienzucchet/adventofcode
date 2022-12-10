package fabienz

import (
	"fmt"
	"io"
	"regexp"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	polymer, rules, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	polymer.serialInsertion(10, rules)

	leastCount, mostCount := polymer.leastMostCommon()

	_, err = fmt.Fprintf(answer, "%d", mostCount-leastCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	polymer, rules, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	polymer.serialInsertion(40, rules)

	leastCount, mostCount := polymer.leastMostCommon()

	_, err = fmt.Fprintf(answer, "%d", mostCount-leastCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Polymer struct {
	pairs    map[string]int
	elements map[string]int
}

// INPUT PARSING

// Use regex to parse rules
var rulesRe = regexp.MustCompile(`^([A-Z]{2}) -> ([A-Z])$`)

func parseLines(lines []string) (polymer Polymer, rules map[string]string, err error) {

	if len(lines) < 3 {
		return Polymer{}, nil, fmt.Errorf("error parsing input : not enough lines")
	}

	// Initialize the polymer
	pairs := make(map[string]int)
	elements := make(map[string]int)

	template := lines[0]

	for i := 0; i < len(template)-1; i++ {
		pairs[template[i:i+2]]++
		elements[template[i:i+1]]++
	}
	// Don't forget the last element
	elements[string(template[len(template)-1])]++

	polymer.pairs = pairs
	polymer.elements = elements

	// Parse all rules
	rules = make(map[string]string)

	for _, line := range lines {
		if rulesRe.MatchString(line) {
			match := rulesRe.FindStringSubmatch(line)
			if len(match) < 3 {
				return Polymer{}, nil, fmt.Errorf("could not parse rule %s", line)
			}

			rules[match[1]] = match[2]
		}
	}

	return polymer, rules, nil
}

// Insertion step
// func insertion(template string, rules map[string]string) (newTemplate string) {

// 	for i:=0; i<len(template)-1; i++ {
// 		newTemplate += string(template[i]) + rules[template[i:i+2]]
// 	}

// 	// Don't forget the last element
// 	newTemplate += string(template[len(template)-1])

// 	return newTemplate
// }

// // Make n insertions
// func serialInsertion(insertionCount int, template string, rules map[string]string) (newTemplate string) {
// 	for i:=0; i<insertionCount; i++ {
// 		template = insertion(template, rules)
// 	}

// 	return template
// }

// We will be using two maps to keep track of the polymer composition : 1 counting the pairs (so that we will be able to process all identical pairs at the same time) and 1 map to count single elements, use for finding the answer

// Iterate polymer insertion
func (p *Polymer) insertion(rules map[string]string) {

	// We must create a new map for the pairs at every iteration
	newPairs := make(map[string]int)

	// For each pair, two pairs are created (ex: NN -> NC and CN i.e. Nrules[NN] and rules[NN]N)
	// For each pair, add one to the count of elements rules[NN]
	for pair, count := range p.pairs {
		newElement := rules[pair]

		newPairs[string(pair[0])+newElement] += count
		newPairs[newElement+string(pair[1])] += count

		p.elements[newElement] += count
	}

	p.pairs = newPairs
}

// Serial insertion
func (p *Polymer) serialInsertion(insertions int, rules map[string]string) {
	for i := 0; i < insertions; i++ {
		p.insertion(rules)
	}
}

// Find least and most common elements counts
func (p *Polymer) leastMostCommon() (leastCount int, mostCount int) {

	leastCount = 1 << 62
	mostCount = -1

	for _, count := range p.elements {
		if count > mostCount {
			mostCount = count
		}

		if count < leastCount {
			leastCount = count
		}
	}

	return leastCount, mostCount
}
