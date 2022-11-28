package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 18 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := parseLines(lines)

	sum := numbers[0]

	for i := 1; i < len(numbers); i++ {
		sum = add(sum, numbers[i])
	}

	_, err = fmt.Fprintf(answer, "%d", sum.magnitude())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 18 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	numbers := parseLines(lines)

	maxMagnitude := -1

	for i := 0; i < len(numbers)-1; i++ {
		for j := i + 1; j < len(numbers); j++ {
			mag := add(numbers[i], numbers[j]).magnitude()
			if mag > maxMagnitude {
				maxMagnitude = mag
			}
			mag = add(numbers[j], numbers[i]).magnitude()
			if mag > maxMagnitude {
				maxMagnitude = mag
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", maxMagnitude)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Number struct {
	value       int
	left, right *Number
	parent      *Number
}

func (n Number) String() string {
	if n.value != -1 {
		return strconv.Itoa(n.value)
	}

	return fmt.Sprintf("[%s,%s]", n.left.String(), n.right.String())
}

// INPUT PARSING

func parseLines(lines []string) (numbers []*Number) {

	for _, line := range lines {
		number := parseNumber(line, nil)
		numbers = append(numbers, number)
	}

	return numbers
}

// Recursive function to parse a number
func parseNumber(s string, parent *Number) (n *Number) {

	// Check if the string to parse is a number
	val, err := strconv.Atoi(s)
	if err == nil {
		n = &Number{value: val, parent: parent}
		return n
	}

	n = &Number{value: -1, parent: parent}

	// Otherwise parse recursively the number
	var left, right string

	// Find the middle comma index of a number.
	brackets := 1
	for i := 1; i < len(s)-1; i++ {
		switch s[i] {
		case '[':
			brackets++
		case ']':
			brackets--
		case ',':
			if brackets == 1 {
				left = s[1:i]
				right = s[i+1 : len(s)-1]
			}
		}
	}

	n.left = parseNumber(left, n)
	n.right = parseNumber(right, n)

	return n
}

// Compute the magnitude
func (n *Number) magnitude() (mag int) {

	if n == nil {
		return 0
	}

	if n.value != -1 {
		return n.value
	}

	return 3*n.left.magnitude() + 2*n.right.magnitude()
}

// Create a copy of a number
func copyNumber(n *Number, parent *Number) (copy *Number) {

	copy = &Number{value: -1, parent: parent}

	if n.value != -1 {
		copy.value = n.value
		return copy
	}

	copy.left = copyNumber(n.left, copy)
	copy.right = copyNumber(n.right, copy)

	return copy
}

// Split a subnumber
func (n *Number) split() {
	n.left = &Number{value: n.value / 2, parent: n}
	n.right = &Number{value: n.value/2 + n.value%2, parent: n}
	n.value = -1
}

// Explode a number
func (n *Number) explode() {
	var upwardNumber *Number

	// Let's find the common ancester with the first left number.
	upwardNumber = n

	for upwardNumber.parent != nil && upwardNumber.parent.right != upwardNumber {
		upwardNumber = upwardNumber.parent
	}

	upwardNumber = upwardNumber.parent

	if upwardNumber != nil {
		// We go down until	we reach the first left number.
		upwardNumber = upwardNumber.left
		for upwardNumber.value == -1 {
			upwardNumber = upwardNumber.right
		}
		upwardNumber.value += n.left.value
	}

	// Do the same process for the first right number.
	upwardNumber = n

	for upwardNumber.parent != nil && upwardNumber.parent.left != upwardNumber {
		upwardNumber = upwardNumber.parent
	}

	upwardNumber = upwardNumber.parent

	if upwardNumber != nil {
		// We go down until	we reach the first left number.
		upwardNumber = upwardNumber.right
		for upwardNumber.value == -1 {
			upwardNumber = upwardNumber.left
		}
		upwardNumber.value += n.right.value
	}

	// n becomes a terminal node with value 0
	n.value = 0
	n.left = nil
	n.right = nil
}

// Search if a number should explode and explode it.
func (n *Number) shouldExplode(depth int) {
	if n.value == -1 {
		if depth == 4 {
			n.explode()
			return
		}

		n.left.shouldExplode(depth + 1)
		n.right.shouldExplode(depth + 1)
	}
}

// Check if a number should splt and split it if needed
func (n *Number) shouldSplit() (hasSplit bool) {
	if n.value > 9 {
		n.split()
		return true
	}

	if n.value == -1 {
		return n.left.shouldSplit() || n.right.shouldSplit()
	}

	return false
}

// Reduce a number
func (n *Number) reduce() {
	// Start by exploding
	n.shouldExplode(0)
	// Split and explode as long as it is needed
	for n.shouldSplit() {
		n.shouldExplode(0)
	}
}

// Add two numbers
func add(n1, n2 *Number) (sum *Number) {

	sum = &Number{value: -1}

	sum.left = copyNumber(n1, sum)
	sum.right = copyNumber(n2, sum)

	sum.reduce()

	return sum
}
