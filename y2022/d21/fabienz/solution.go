package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the expressions in the file.
	tree, err := expressionTreeFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse expressions: %w", err)
	}

	res, err := tree.evaluate("root")
	if err != nil {
		return fmt.Errorf("could not evaluate expression: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the expressions in the file.
	tree, err := expressionTreeFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse expressions: %w", err)
	}

	// Remove the "humn" expression.
	delete(tree, "humn")

	// Create a new tree to store the inverted expressions.
	invertedTree := make(expressionTree)

	// In the inverted tree, the two operands of the root expression are equal.

	// Find the "root" expression.
	root, ok := tree["root"]
	if !ok {
		return fmt.Errorf("could not find root expression")
	}

	// Evaluate both operands of the root expression if possible.
	loperand, err := tree.evaluate(root.loperand)
	if err != nil && !containsHumn(err) {
		return fmt.Errorf("could not evaluate left operand: %w", err)
	}
	if !containsHumn(err) {
		invertedTree[root.roperand] = expression{
			isNumber: true,
			value:    loperand,
		}
	}

	roperand, err := tree.evaluate(root.roperand)
	if err != nil && !containsHumn(err) {
		return fmt.Errorf("could not evaluate right operand: %w", err)
	}
	if !containsHumn(err) {
		invertedTree[root.loperand] = expression{
			isNumber: true,
			value:    roperand,
		}
	}

	// For each expression in the tree, invert it.
	for name, expr := range tree {
		// Skip the root expression.
		if name == "root" {
			continue
		}

		// Skip the expressions that have already been evaluated.
		if expr.isNumber {
			// invertedTree[name] = expr
			continue
		}

		// Evaluate the left operand if possible.
		loperand, err := tree.evaluate(expr.loperand)
		if err != nil && !containsHumn(err) {
			return fmt.Errorf("could not evaluate left operand: %w", err)
		}

		// Evaluate the right operand if possible.
		roperand, err := tree.evaluate(expr.roperand)
		if err != nil && !containsHumn(err) {
			return fmt.Errorf("could not evaluate right operand: %w", err)
		}
		// If the right operand is a number, add the expression where the left operand is the unknown to the inverted tree.
		if !containsHumn(err) {
			// Add the evaluated value of the right operand to the expression tree.
			invertedTree[expr.roperand] = expression{
				isNumber: true,
				value:    roperand,
			}
			newName, e := invertUnknownLeft(name, expr)
			invertedTree[newName] = e
		} else {
			// If the right operand is not a number, add the expression where the right operand is the unknown to the inverted tree.
			// Add the evaluated value of the left operand to the expression tree.
			invertedTree[expr.loperand] = expression{
				isNumber: true,
				value:    loperand,
			}
			newName, e := invertUnknownRight(name, expr)
			invertedTree[newName] = e
		}
	}

	// Evaluate the human expression.
	res, err := invertedTree.evaluate("humn")
	if err != nil {
		return fmt.Errorf("could not evaluate expression: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type expression struct {
	isNumber bool
	value    int
	operator string
	loperand string
	roperand string
}

type expressionTree map[string]expression

var numberRegex = regexp.MustCompile(`^([a-z]+): ([0-9]+)$`)
var operatorRegex = regexp.MustCompile(`^([a-z]+): ([a-z]+) ([+-\/*]) ([a-z]+)$`)

// Parse the lines to build an expression tree.
func expressionTreeFromLines(lines []string) (expressionTree, error) {
	tree := make(expressionTree)

	for _, line := range lines {
		switch {
		case numberRegex.MatchString(line):
			matches := numberRegex.FindStringSubmatch(line)
			if len(matches) != 3 {
				return nil, fmt.Errorf("could not parse line %q", line)
			}
			name := matches[1]
			value, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, fmt.Errorf("could not parse line %q: %w", line, err)
			}
			tree[name] = expression{
				isNumber: true,
				value:    value,
			}
		case operatorRegex.MatchString(line):
			matches := operatorRegex.FindStringSubmatch(line)
			if len(matches) != 5 {
				return nil, fmt.Errorf("could not parse line %q", line)
			}
			name := matches[1]
			loperand := matches[2]
			operator := matches[3]
			roperand := matches[4]
			tree[name] = expression{
				isNumber: false,
				operator: operator,
				loperand: loperand,
				roperand: roperand,
			}
		}
	}

	return tree, nil
}

// Recursive function to evaluate an expression.
func (tree expressionTree) evaluate(name string) (int, error) {
	expr, ok := tree[name]
	if !ok {
		return 0, fmt.Errorf("unknown expression %q", name)
	}

	if expr.isNumber {
		return expr.value, nil
	}

	loperand, err := tree.evaluate(expr.loperand)
	if err != nil {
		return 0, fmt.Errorf("could not evaluate %q: %w", expr.loperand, err)
	}

	roperand, err := tree.evaluate(expr.roperand)
	if err != nil {
		return 0, fmt.Errorf("could not evaluate %q: %w", expr.roperand, err)
	}

	switch expr.operator {
	case "+":
		return loperand + roperand, nil
	case "-":
		return loperand - roperand, nil
	case "*":
		return loperand * roperand, nil
	case "/":
		return loperand / roperand, nil
	default:
		return 0, fmt.Errorf("unknown operator %q", expr.operator)
	}
}

// Check if the error contains the string "humn".
func containsHumn(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "humn")
}

// Invert an expression where the left operand is the unknown value.
func invertUnknownLeft(name string, expr expression) (string, expression) {
	e := expression{
		isNumber: false,
	}

	switch expr.operator {
	case "+":
		e.operator = "-"
		e.loperand = name
		e.roperand = expr.roperand
	case "-":
		e.operator = "+"
		e.loperand = name
		e.roperand = expr.roperand
	case "*":
		e.operator = "/"
		e.loperand = name
		e.roperand = expr.roperand
	case "/":
		e.operator = "*"
		e.loperand = name
		e.roperand = expr.roperand
	}

	return expr.loperand, e
}

// Invert an expression where the right operand is the unknown value.
func invertUnknownRight(name string, expr expression) (string, expression) {
	e := expression{
		isNumber: false,
	}

	switch expr.operator {
	case "+":
		e.operator = "-"
		e.loperand = name
		e.roperand = expr.loperand
	case "-":
		e.operator = "-"
		e.loperand = expr.loperand
		e.roperand = name
	case "*":
		e.operator = "/"
		e.loperand = name
		e.roperand = expr.loperand
	case "/":
		e.operator = "/"
		e.loperand = expr.loperand
		e.roperand = name
	}

	return expr.roperand, e
}
