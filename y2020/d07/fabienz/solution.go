package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, err := splitLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	known_results := make(map[string]bool)

	count := 0

	for bag := range rules {
		if contains_shiny_bag(bag, rules, known_results) {
			count = count + 1
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count-1)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, err := splitLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	known_results := make(map[string]int)

	_, err = fmt.Fprintf(answer, "%d", count_sub_bags("shiny gold", rules, known_results)-1)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Parse a rule and returns a map ofthe bags contained
func parseRule(rule string) (string, map[string]int, error) {
	bag_re := regexp.MustCompile("^(.*) bags contain")
	sub_bags_re := regexp.MustCompile(`([0-9] [a-z\s]*) bag`)

	bag := bag_re.FindStringSubmatch(rule)[1]
	sub_bags := sub_bags_re.FindAllStringSubmatch(rule, -1)

	if len(sub_bags) == 0 {
		return bag, nil, nil
	}

	contains := make(map[string]int)

	for _, match := range sub_bags {
		b := match[1][2:]
		nb, err := strconv.Atoi(match[1][:1])
		if err != nil {
			return "", nil, fmt.Errorf("error parsing value %s : %w", match[1][:1], err)
		}

		contains[b] = nb
	}

	return bag, contains, nil
}

// Converts the slice of rules into a map of map containing the rules
func splitLines(lines []string) (map[string]map[string]int, error) {
	rules := make(map[string]map[string]int)

	for _, rule := range lines {
		b, contains, err := parseRule(rule)
		if err != nil {
			return nil, fmt.Errorf("error parsing rule %s : %w", rule, err)
		}

		rules[b] = contains
	}

	return rules, nil
}

// Recursive function that checks if the bag contains at least one shiny gold bag
func contains_shiny_bag(bag string, rules map[string]map[string]int, known_results map[string]bool) bool {
	// Check if we already know the result
	if _, ok := known_results[bag]; ok {
		return known_results[bag]
	} else if rules[bag] == nil {
		known_results[bag] = false
		return false
	} else if bag == "shiny gold" {
		known_results[bag] = true
		return true
	}
	for sub_bag, _ := range rules[bag] {
		if contains_shiny_bag(sub_bag, rules, known_results) {
			known_results[bag] = true
			return true
		}
	}
	known_results[bag] = false
	return false
}

// Recursive function that counts the number if bags inside
func count_sub_bags(bag string, rules map[string]map[string]int, known_results map[string]int) int {
	// Check if we already know the result
	if _, ok := known_results[bag]; ok {
		return known_results[bag]
	} else if rules[bag] == nil {
		known_results[bag] = 1
		return 1
	}
	count := 1
	for sub_bag, nb := range rules[bag] {
		count = count + nb*count_sub_bags(sub_bag, rules, known_results)
	}
	known_results[bag] = count
	return count
}
