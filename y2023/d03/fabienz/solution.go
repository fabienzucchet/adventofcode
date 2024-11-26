package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 3 of Advent of Code 2023.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	s, err := schematicFromInput(lines)
	if err != nil {
		return fmt.Errorf("could not create schematic: %w", err)
	}

	s.UpdateParts()

	sum := 0
	for _, p := range s.parts {
		sum += p.partNumber
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 3 of Advent of Code 2023.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	s, err := schematicFromInput(lines)
	if err != nil {
		return fmt.Errorf("could not create schematic: %w", err)
	}

	s.UpdateParts()
	s.UpdateGears()

	sum := 0
	for _, g := range s.gears {
		sum += g.gearRadius
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Schematic struct {
	elements             []*SchematicElement
	visualRepresentation [][]rune
	parts                []*SchematicElement
	gears                []*SchematicElement
}

type SchematicElement struct {
	coordinates [2]int
	length      int
	value       string

	isSymbol   bool
	partNumber int
	gearRadius int
}

func (s *Schematic) AddElement(e *SchematicElement) {
	s.elements = append(s.elements, e)
}

func (s *Schematic) String() string {
	var sb strings.Builder

	for _, line := range s.visualRepresentation {
		for _, c := range line {
			sb.WriteRune(c)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (e *SchematicElement) String() string {
	return fmt.Sprintf("%s (length: %d) @ (%d,%d)", e.value, e.length, e.coordinates[0], e.coordinates[1])
}

func (s *Schematic) UpdateParts() {
	parts := []*SchematicElement{}
	for _, e := range s.elements {
		if !e.isSymbol && s.isPart(e) {
			e.partNumber = e.getPartNumber()
			parts = append(parts, e)
		}
	}

	s.parts = parts
}

func (e *SchematicElement) getPartNumber() int {
	v, err := strconv.Atoi(e.value)
	if err != nil {
		panic(err)
	}
	return v
}

func (s *Schematic) UpdateGears() {
	gears := []*SchematicElement{}
	for _, e := range s.elements {
		if e.isSymbol && s.isGear(e) {
			gears = append(gears, e)
		}
	}

	s.gears = gears
}

func schematicFromInput(lines []string) (*Schematic, error) {
	s := &Schematic{}

	// Initialize the visual representation of the schematic
	s.visualRepresentation = make([][]rune, len(lines))
	for x, line := range lines {
		s.visualRepresentation[x] = make([]rune, len(line))
		for y, c := range line {
			s.visualRepresentation[x][y] = c
		}
	}

	// Extract the elements from the input
	for x, line := range lines {
		i := 0

		for i < len(line) {
			switch {
			case line[i] >= '0' && line[i] <= '9':
				// Extract the length
				length := 1
				for i+length < len(line) && line[i+length] >= '0' && line[i+length] <= '9' {
					length++
				}

				s.AddElement(&SchematicElement{
					coordinates: [2]int{x, i},
					length:      length,
					value:       line[i : i+length],
				})

				i += length
			case line[i] == '.':
				i++
			default:
				s.AddElement(&SchematicElement{
					coordinates: [2]int{x, i},
					length:      1,
					value:       line[i : i+1],
					isSymbol:    true,
				})
				i++
			}
		}
	}

	return s, nil
}

func (s *Schematic) isPart(e *SchematicElement) bool {
	// Check if an element is adjacent to a symbol
	neighbours := [][2]int{}

	// Check the element above
	if e.coordinates[0] > 0 {
		for i := 0; i < e.length; i++ {
			neighbours = append(neighbours, [2]int{e.coordinates[0] - 1, e.coordinates[1] + i})
		}
	}

	// Check the element below
	if e.coordinates[0] < len(s.visualRepresentation)-1 {
		for i := 0; i < e.length; i++ {
			neighbours = append(neighbours, [2]int{e.coordinates[0] + 1, e.coordinates[1] + i})
		}
	}

	// Check the element to the left
	if e.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{e.coordinates[0], e.coordinates[1] - 1})
	}

	// Check the element to the right
	if e.coordinates[1]+e.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{e.coordinates[0], e.coordinates[1] + e.length})
	}

	// Check the diagonal elements
	if e.coordinates[0] > 0 && e.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{e.coordinates[0] - 1, e.coordinates[1] - 1})
	}
	if e.coordinates[0] > 0 && e.coordinates[1]+e.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{e.coordinates[0] - 1, e.coordinates[1] + e.length})
	}
	if e.coordinates[0] < len(s.visualRepresentation)-1 && e.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{e.coordinates[0] + 1, e.coordinates[1] - 1})
	}
	if e.coordinates[0] < len(s.visualRepresentation)-1 && e.coordinates[1]+e.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{e.coordinates[0] + 1, e.coordinates[1] + e.length})
	}

	for _, n := range neighbours {
		if s.visualRepresentation[n[0]][n[1]] != '.' && !isDigit(s.visualRepresentation[n[0]][n[1]]) {
			return true
		}
	}

	return false
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Schematic) isGear(e *SchematicElement) bool {
	adjacentParts := []*SchematicElement{}
	for _, p := range s.parts {
		if s.areAdjacent(p, e) {
			adjacentParts = append(adjacentParts, p)
		}
	}

	if len(adjacentParts) == 2 {
		e.gearRadius = adjacentParts[0].partNumber * adjacentParts[1].partNumber

		return true
	}

	return false
}

func (s *Schematic) areAdjacent(p *SchematicElement, e *SchematicElement) bool {
	// Check if an element is adjacent to a symbol
	neighbours := [][2]int{}

	// Check the element above
	if p.coordinates[0] > 0 {
		for i := 0; i < p.length; i++ {
			neighbours = append(neighbours, [2]int{p.coordinates[0] - 1, p.coordinates[1] + i})
		}
	}

	// Check the element below
	if p.coordinates[0] < len(s.visualRepresentation)-1 {
		for i := 0; i < p.length; i++ {
			neighbours = append(neighbours, [2]int{p.coordinates[0] + 1, p.coordinates[1] + i})
		}
	}

	// Check the element to the left
	if p.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{p.coordinates[0], p.coordinates[1] - 1})
	}

	// Check the element to the right
	if p.coordinates[1]+p.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{p.coordinates[0], p.coordinates[1] + p.length})
	}

	// Check the diagonal elements
	if p.coordinates[0] > 0 && p.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{p.coordinates[0] - 1, p.coordinates[1] - 1})
	}
	if p.coordinates[0] > 0 && p.coordinates[1]+p.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{p.coordinates[0] - 1, p.coordinates[1] + p.length})
	}
	if p.coordinates[0] < len(s.visualRepresentation)-1 && p.coordinates[1] > 0 {
		neighbours = append(neighbours, [2]int{p.coordinates[0] + 1, p.coordinates[1] - 1})
	}
	if p.coordinates[0] < len(s.visualRepresentation)-1 && p.coordinates[1]+p.length < len(s.visualRepresentation[0])-1 {
		neighbours = append(neighbours, [2]int{p.coordinates[0] + 1, p.coordinates[1] + p.length})
	}

	for _, n := range neighbours {
		if n == e.coordinates {
			return true
		}
	}

	return false
}
