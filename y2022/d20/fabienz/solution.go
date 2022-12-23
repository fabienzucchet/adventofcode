package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 20 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	encrypted, err := listFromLines(lines, 1)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Copy the encrypted list.
	l := encrypted.Copy()

	// Iterate on each node of the list and mix it
	for n := range encrypted.Iterate() {
		// Find the index of the node.
		idx, err := l.IndexOf(n.value, n.initialIdx)
		if err != nil {
			return fmt.Errorf("could not find node: %w", err)
		}

		// Remove the node at this index.
		l.Remove(idx)

		// Insert the node at the index + n.value.
		l.Add(n.value, n.initialIdx, idx+n.value)
	}

	// Compute the grove coordinates.
	res, err := computeGroveCoordinates(l)
	if err != nil {
		return fmt.Errorf("could not compute grove coordinates: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	encrypted, err := listFromLines(lines, 811589153)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Copy the encrypted list.
	l := encrypted.Copy()

	// Mix 10 times the list.
	for i := 0; i < 10; i++ {
		// Iterate on each node of the list and mix it
		for n := range encrypted.Iterate() {
			// Find the index of the node.
			idx, err := l.IndexOf(n.value, n.initialIdx)
			if err != nil {
				return fmt.Errorf("could not find node: %w", err)
			}

			// Simplify the index by the size of the list.
			idx %= l.size

			// Remove the node at this index.
			l.Remove(idx)

			newIdx := (idx + n.value) % l.size

			// Insert the node at the index + n.value.
			l.Add(n.value, n.initialIdx, newIdx)
		}
	}

	// Compute the grove coordinates.
	res, err := computeGroveCoordinates(l)
	if err != nil {
		return fmt.Errorf("could not compute grove coordinates: %w", err)
	}

	// TODO: Write your solution to Part 2 below.
	_, err = fmt.Fprintf(answer, "%d", res)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Implement a circular linked list.
type CircularLinkedList struct {
	// The current node.
	current *Node

	// The number of nodes in the list.
	size int
}

// A node in the list.
type Node struct {
	// The value of the node.
	value int

	// Make the nodes unique with the initial index.
	initialIdx int

	// The next node in the list.
	next *Node
}

// Get the node at a given index.
func (l *CircularLinkedList) Get(index int) *Node {
	// Make sure that the index is >= 0.
	for index < 0 {
		index += l.size
	}

	// Make sure that the index is < l.size.
	index %= l.size

	// Get the node at the given index.
	node := l.current
	for i := 0; i < index; i++ {
		node = node.next
	}

	return node
}

// Get the index of a node.
func (l *CircularLinkedList) IndexOf(value int, initialIdx int) (int, error) {
	node := l.current
	for i := 0; i < l.size; i++ {
		if node.value == value && node.initialIdx == initialIdx {
			return i, nil
		}
		node = node.next
	}

	return -1, fmt.Errorf("could not find node with value %d and initial index %d", value, initialIdx)
}

// Get the first index of a node with a given value.
func (l *CircularLinkedList) IndexOfValue(value int) (int, error) {
	node := l.current
	for i := 0; i < l.size; i++ {
		if node.value == value {
			return i, nil
		}
		node = node.next
	}

	return -1, fmt.Errorf("could not find node with value %d", value)
}

// Add a new node to the list at a given index.
func (l *CircularLinkedList) Add(value int, initialIdx int, index int) {
	node := &Node{value: value, initialIdx: initialIdx}

	if l.size == 0 {
		node.next = node
		l.current = node
	} else {
		previous := l.Get(index - 1)
		node.next = previous.next
		previous.next = node
	}

	l.size++
}

// Remove a node from the list at a given index.
func (l *CircularLinkedList) Remove(index int) {
	if l.size == 0 {
		return
	}

	previous := l.Get(index - 1)
	previous.next = previous.next.next

	if index == 0 {
		l.current = previous.next
	}

	l.size--
}

// String returns a string representation of the list.
func (l *CircularLinkedList) String() string {
	if l.size == 0 {
		return "[]"
	}

	s := "["
	node := l.current
	for i := 0; i < l.size; i++ {
		s += fmt.Sprintf("%d", node.value)
		if i < l.size-1 {
			s += " "
		}
		node = node.next
	}
	s += "]"

	return s
}

// Copy a list.
func (l *CircularLinkedList) Copy() *CircularLinkedList {
	newList := &CircularLinkedList{}

	node := l.current
	for i := 0; i < l.size; i++ {
		newList.Add(node.value, node.initialIdx, newList.size)
		node = node.next
	}

	return newList
}

// Iterate over the list.
func (l *CircularLinkedList) Iterate() <-chan *Node {
	ch := make(chan *Node)

	go func() {
		defer close(ch)

		node := l.current
		for i := 0; i < l.size; i++ {
			ch <- node
			node = node.next
		}
	}()

	return ch
}

// Create a new circular linked list from the input.
func listFromLines(lines []string, key int) (*CircularLinkedList, error) {
	list := &CircularLinkedList{}

	for i, line := range lines {
		value, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse line %q: %w", line, err)
		}

		list.Add(key*value, i, list.size)
	}

	return list, nil
}

// Compute the grove coordinates from a list.
func computeGroveCoordinates(list *CircularLinkedList) (int, error) {
	// Find the index of the node with value 0.
	idx, err := list.IndexOfValue(0)
	if err != nil {
		return -1, fmt.Errorf("could not find node with value 0: %w", err)
	}

	sum := 0

	offsetsToConsider := []int{1000, 2000, 3000}

	for _, offset := range offsetsToConsider {
		node := list.Get(idx + offset)
		sum += node.value
	}

	return sum, nil
}
