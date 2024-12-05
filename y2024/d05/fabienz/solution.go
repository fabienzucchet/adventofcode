package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, updates, err := rulesAndUpdatesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse rules and updates: %w", err)
	}

	validUpdates := filterInvalidUpdates(rules, updates)

	res := sumMiddleValuesOfValidUpdates(validUpdates)

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, updates, err := rulesAndUpdatesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse rules and updates: %w", err)
	}

	invalidUpdates := filterValidUpdates(updates, rules)

	reorderUpdates(rules, invalidUpdates)

	res := sumMiddleValuesOfValidUpdates(invalidUpdates)

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Store the rules in a map indexed by the two integers (previous, next)
type Rules map[[2]int]bool

func (r Rules) addRule(prev, next int) {
	r[[2]int{prev, next}] = true
}

func (r Rules) hasRule(prev, next int) bool {
	_, exists := r[[2]int{prev, next}]
	return exists
}

type Update []int

func rulesAndUpdatesFromLines(lines []string) (Rules, []Update, error) {
	r := make(Rules)
	u := make([]Update, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}

		var prev, next int
		if _, err := fmt.Sscanf(line, "%d|%d", &prev, &next); err == nil {
			r.addRule(prev, next)
		} else {
			update, err := helpers.IntsFromString(line, ",")
			if err != nil {
				return nil, nil, fmt.Errorf("could not parse line %q: %w", line, err)
			}

			u = append(u, update)
		}
	}
	return r, u, nil
}

// An update is in the correct order if for every current value current, the rule (current, previous) does NOT exist for each previous value in the update.
func (u Update) isCorrectOrder(r Rules) bool {
	for i, current := range u {
		if i == 0 {
			continue
		}

		for j := 0; j < i; j++ {
			if r.hasRule(current, u[j]) {
				return false
			}
		}
	}

	return true
}

// Get the middle value of the update
func (u Update) middleValue() int {
	return u[len(u)/2]
}

func filterInvalidUpdates(r Rules, updates []Update) []Update {
	valid := make([]Update, 0)
	for _, u := range updates {
		if u.isCorrectOrder(r) {
			valid = append(valid, u)
		}
	}
	return valid
}

// Add the middle value of all valid updates
func sumMiddleValuesOfValidUpdates(updates []Update) int {
	sum := 0
	for _, u := range updates {
		sum += u.middleValue()
	}
	return sum
}

// We can use a bubble sorting algorithm to reorder the update using the rules
func (u Update) reorder(r Rules) {
	for i := 0; i < len(u); i++ {
		for j := i + 1; j < len(u); j++ {
			if r.hasRule(u[j], u[i]) {
				u[i], u[j] = u[j], u[i]
			}
		}
	}
}

func filterValidUpdates(updates []Update, r Rules) []Update {
	invalid := make([]Update, 0)
	for _, u := range updates {
		if !u.isCorrectOrder(r) {
			invalid = append(invalid, u)
		}
	}
	return invalid
}

func reorderUpdates(r Rules, updates []Update) {
	for i := range updates {
		updates[i].reorder(r)
	}
}
