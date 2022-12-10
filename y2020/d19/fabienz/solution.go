package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, messages, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("an error occured when parsing the input : %w", err)
	}

	count := 0

	for _, message := range messages {
		indexes := applyRule(0, message, rules)
		for _, idx := range indexes {
			if idx == len(message) {
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	rules, messages, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("an error occured when parsing the input : %w", err)
	}

	var children8 [][]int
	children8 = append(children8, []int{42})
	children8 = append(children8, []int{42, 8})
	rules[8] = Rule{
		children: children8,
	}

	var children11 [][]int
	children11 = append(children11, []int{42, 31})
	children11 = append(children11, []int{42, 11, 31})
	rules[11] = Rule{
		children: children11,
	}

	count := 0

	for _, message := range messages {
		indexes := applyRule(0, message, rules)
		for _, idx := range indexes {
			if idx == len(message) {
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseLines(lines []string) (rules map[int]Rule, messages []string, err error) {

	targetSlice := "rules"
	rules = make(map[int]Rule)

	for _, line := range lines {
		if line == "" {
			targetSlice = "messages"
		} else {
			switch targetSlice {
			case "rules":
				id, rule, err := parseRule(line)
				if err != nil {
					return nil, nil, fmt.Errorf("error when parsing rule %s : %w", line, err)
				}

				rules[id] = rule
			case "messages":
				messages = append(messages, line)
			}
		}
	}

	return rules, messages, nil
}

// We use regex to parse rule
var ruleRegex = regexp.MustCompile(`^([0-9]+): (.*)$`)
var explicitRuleRegex = regexp.MustCompile(`^"[ab]"$`)

// We store data about the rules in the following struct :
// value contains a string if the rule is an explicit rule (example : "a")
// children contains a list of list of indexes. Each list is a possibility due to OR in conditions
// Ex: 0: 4 1 5 -> children = [[4 1 5]]
// Ex: 0: 1 2 | 2 1 -> children = [[1 2] [2 1]]
type Rule struct {
	value    string
	children [][]int
}

// Parse a rule
func parseRule(rule string) (id int, r Rule, err error) {

	match := ruleRegex.FindStringSubmatch(rule)

	if len(match) != 3 {
		return 0, r, fmt.Errorf("error when parsing rule %s", rule)
	}

	id, err = strconv.Atoi(match[1])
	if err != nil {
		return 0, r, fmt.Errorf("error when parsing rule id for rule %s : %w", rule, err)
	}

	content := match[2]

	// If rule is explicit i.e. it has no children
	if explicitRuleRegex.MatchString(content) {
		r.value = content[1:2] // Keep only the letter
	} else {
		// Parse the different possibilities
		possibilities := strings.Split(content, "|")

		for _, possibility := range possibilities {
			var children []int
			for _, child := range strings.Fields(possibility) {
				id, err := strconv.Atoi(child)
				if err != nil {
					return 0, r, fmt.Errorf("error while parsing id %s : %w", child, err)
				}
				children = append(children, id)
			}

			r.children = append(r.children, children)
		}
	}

	return id, r, nil
}

func applyRule(ruleId int, message string, rules map[int]Rule) (indexes []int) {
	rule := rules[ruleId]

	if rule.value != "" {
		if len(message) == 0 {
			return nil
		}
		if message[:1] == rule.value {
			return []int{1}
		}
		return nil
	}

	for _, possibility := range rule.children {
		indexes = append(indexes, applyRemainingRules(possibility, message, rules)...)
	}
	return indexes
}

func applyRemainingRules(possibility []int, message string, rules map[int]Rule) (indexes []int) {
	if len(possibility) == 0 {
		return nil
	}

	if len(possibility) == 1 {
		return applyRule(possibility[0], message, rules)
	}

	for _, firstIndex := range applyRule(possibility[0], message, rules) {
		for _, idxs := range applyRemainingRules(possibility[1:], message[firstIndex:], rules) {
			indexes = append(indexes, idxs+firstIndex)
		}
	}

	return indexes
}
