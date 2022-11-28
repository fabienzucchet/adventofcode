package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Precedence is the map containing the priority of the operators
	precedence := map[string]int{
		"+": 1,
		"*": 1,
	}

	sum := 0

	for idx, line := range lines {
		res, err := evaluate(line, precedence)
		if err != nil {
			return fmt.Errorf("an error occured while evaluating %s (%d): %w", line, idx, err)
		}
		sum += res
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Precedence is the map containing the priority of the operators
	precedence := map[string]int{
		"+": 2,
		"*": 1,
	}

	sum := 0

	for idx, line := range lines {
		res, err := evaluate(line, precedence)
		if err != nil {
			return fmt.Errorf("an error occured while evaluating %s (%d): %w", line, idx, err)
		}
		sum += res
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Implement stacks to solve this problem
type OpStack []string

func (s *OpStack) Push(elt string) {
	*s = append(*s, elt)
}

func (s *OpStack) IsEmpty() (isEmpty bool) {
	return len(*s) == 0
}

func (s *OpStack) View() (elt string, exists bool) {
	if s.IsEmpty() {
		return "", false
	}

	lastIdx := len(*s) - 1
	elt = (*s)[lastIdx]

	return elt, true
}

func (s *OpStack) Pop() (elt string, exists bool) {
	if s.IsEmpty() {
		return "", false
	}

	lastIdx := len(*s) - 1
	elt = (*s)[lastIdx]
	*s = (*s)[:lastIdx]

	return elt, true
}

type ValStack []int

func (s *ValStack) Push(elt int) {
	*s = append(*s, elt)
}

func (s *ValStack) IsEmpty() (isEmpty bool) {
	return len(*s) == 0
}

func (s *ValStack) Pop() (elt int, exists bool) {
	if s.IsEmpty() {
		return 0, false
	}

	lastIdx := len(*s) - 1
	elt = (*s)[lastIdx]
	*s = (*s)[:lastIdx]

	return elt, true
}

func (s *ValStack) View() (elt int, exists bool) {
	if s.IsEmpty() {
		return 0, false
	}

	lastIdx := len(*s) - 1
	elt = (*s)[lastIdx]

	return elt, true
}

// We use some regexp to determine if a token is composed of integers
var isNumber = regexp.MustCompile(`^[0-9]+$`)

// This function applies the operator the the two operands
func applyOperator(operator string, operandA, operandB int) (res int) {
	switch operator {
	case "*":
		res = operandA * operandB
	case "+":
		res = operandA + operandB
	}

	return res
}

// Display information about the stack for debug purpose
func (s *ValStack) displayStack() {
	elt, _ := s.View()
	helpers.Println("Value stack : length=", len(*s), "element on the top=", elt, "stack=", *s)
}
func (s *OpStack) displayStack() {
	elt, _ := s.View()
	helpers.Println("Operator stack : length=", len(*s), "element on the top=", elt, "stack=", *s)
}

// We use regex to detect if a token contains parenthesis in the next function
var isLeftParenthesis = regexp.MustCompile(`^\(+[0-9]*$`)
var isRightParenthesis = regexp.MustCompile(`^[0-9]*\)+$`)

// Tokenization function handling the parenthesis (Fields doesn't spilt parenthesis)
func tokenize(expression string) (tokens []string) {

	for _, token := range strings.Fields(expression) {
		switch {
		case isLeftParenthesis.MatchString(token):
			for token[0] == '(' {
				tokens = append(tokens, token[:1])
				token = token[1:]
			}
			tokens = append(tokens, token)

		case isRightParenthesis.MatchString(token):
			var appendStack []string
			for token[len(token)-1] == ')' {
				appendStack = append(appendStack, token[len(token)-1:])
				token = token[:len(token)-1]
			}
			tokens = append(tokens, token)
			tokens = append(tokens, appendStack...)

		default:
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// This function evaluates one expression
//We use the algorithm explained in https://www.geeksforgeeks.org/expression-evaluation/
func evaluate(expression string, precedence map[string]int) (result int, err error) {
	tokens := tokenize(expression)

	valStack := ValStack{}
	opStack := OpStack{}

	for _, token := range tokens {

		switch {
		case isNumber.MatchString(token):
			// Push the value in the value stack
			tok, err := strconv.Atoi(token)
			if err != nil {
				return 0, fmt.Errorf("error when parsing token %s : %w", token, err)
			}
			valStack.Push(tok)
		case isLeftParenthesis.MatchString(token):
			// Push the parenthesis in the operator stack
			opStack.Push(token[:1])

			// Push the value just after the parenthesis in the value stack
			if len(token) > 1 {
				tok, err := strconv.Atoi(token[1:])
				if err != nil {
					return 0, fmt.Errorf("error when parsing token %s : %w", token[1:], err)
				}
				valStack.Push(tok)
			}
		case isRightParenthesis.MatchString(token):

			// We need to consider the value just before the parenthesis
			if len(token) > 1 {
				val, err := strconv.Atoi(token[:len(token)-1])
				if err != nil {
					return 0, fmt.Errorf("error when parsing token %s : %w", token[:len(token)-1], err)
				}
				valStack.Push(val)
			}

			// While we haven't encouter the corresponding parenthesis
			op, exists := opStack.Pop()
			for exists && op != "(" {
				// Fetch the two operands
				operandB, exists := valStack.Pop()
				if !exists {
					return 0, fmt.Errorf("unexpected empty value stack")
				}
				operandA, exists := valStack.Pop()
				if !exists {
					return 0, fmt.Errorf("unexpected empty value stack")
				}

				// Apply the operator to the two operands and push the result in the value stack
				res := applyOperator(op, operandA, operandB)
				valStack.Push(res)

				op, exists = opStack.Pop()
			}
		// The default case corresponds to other operators (+, *, etc). Unsupported operators will replace the two operands by 0
		default:
			// While the operator stack is not empty and the operator in the operator stack has a greater or equal precedence than the token
			op, exists := opStack.View()
			for exists && precedence[token] <= precedence[op] {
				operator, exists := opStack.Pop()
				if !exists {
					return 0, fmt.Errorf("could not retrieve operator %s in the operator stack", op)
				}

				// Fetch the two operands
				operandB, exists := valStack.Pop()
				if !exists {
					return 0, fmt.Errorf("unexpected empty value stack")
				}
				operandA, exists := valStack.Pop()
				if !exists {
					return 0, fmt.Errorf("unexpected empty value stack")
				}

				// Apply the operator to the two operands and push the result in the value stack
				res := applyOperator(operator, operandA, operandB)
				valStack.Push(res)

				op, exists = opStack.View()
			}

			// Push the token in the operator stack
			opStack.Push(token)
		}

	}

	isEmpty := opStack.IsEmpty()
	// We now have to empty the operator stack
	for !isEmpty {
		op, exists := opStack.Pop()
		if !exists {
			return 0, fmt.Errorf("unexpected empty operator stack")
		}

		// Fetch the two operands
		operandB, exists := valStack.Pop()
		if !exists {
			return 0, fmt.Errorf("unexpected empty value stack")
		}
		operandA, exists := valStack.Pop()
		if !exists {
			return 0, fmt.Errorf("unexpected empty value stack")
		}

		// Apply the operator to the two operands and push the result in the value stack
		res := applyOperator(op, operandA, operandB)
		valStack.Push(res)

		isEmpty = opStack.IsEmpty()
	}

	// There should be only one element in the value stack : the result
	result, exists := valStack.Pop()

	if !exists {
		return 0, fmt.Errorf("unexpected empty value stack")
	}

	if !valStack.IsEmpty() {
		return 0, fmt.Errorf("value stack should be empty but is not")
	}

	return result, nil
}
