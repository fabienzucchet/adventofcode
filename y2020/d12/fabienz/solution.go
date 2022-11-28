package fabienz

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	pos := Pos{0, 0, 'E'}

	navigate(&pos, lines)

	_, err = fmt.Fprintf(answer, "%d", Abs(pos.E)+Abs(pos.N))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	wayPos := Pos2{10, 1}
	pos2 := Pos2{}

	navigate2(&wayPos, &pos2, lines)

	_, err = fmt.Fprintf(answer, "%d", Abs(pos2.E)+Abs(pos2.N))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Pos struct {
	E      int
	N      int
	Facing rune
}

type Pos2 struct {
	E int
	N int
}

func ReadLines(r io.Reader) ([]string, error) {
	var lines []string

	// Read content of the file
	scanner := bufio.NewScanner(r)

	// Print file content line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseLine(line string) (rune, int, error) {

	n, err := strconv.Atoi(line[1:])
	if err != nil {
		return ' ', 0, fmt.Errorf("error parsing %v : %w", line[1:], err)
	}

	return rune(line[0]), n, nil
}

var directions = [4]rune{'N', 'E', 'S', 'W'}
var directionsIndex = map[rune]int{
	'N': 0,
	'E': 1,
	'S': 2,
	'W': 3,
}

func turnRight(pos *Pos, value int, direction [4]rune, directionsIndex map[rune]int) *Pos {
	delta := value / 90

	idx := directionsIndex[pos.Facing]

	pos.Facing = directions[(idx+delta)%4]

	return pos
}

func navigate(pos *Pos, instructions []string) (*Pos, error) {

	// If no more instructions
	if len(instructions) == 0 {
		return pos, nil
	}

	instruction := instructions[0]

	letter, value, err := parseLine(instruction)
	if err != nil {
		return nil, fmt.Errorf("error parsing instruction %s : %w", instruction, err)
	}

	if letter == 'F' {
		letter = pos.Facing
	}

	switch letter {
	case 'N':
		pos.N = pos.N + value
	case 'S':
		pos.N = pos.N - value
	case 'E':
		pos.E = pos.E + value
	case 'W':
		pos.E = pos.E - value
	case 'R':
		pos = turnRight(pos, value, directions, directionsIndex)
	case 'L':
		value = -value
		for value < 0 {
			value = value + 360
		}
		pos = turnRight(pos, value, directions, directionsIndex)
	}

	return navigate(pos, instructions[1:])
}

func Abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Turn right of 90 degrees
func turnRight2(wayPos *Pos2) *Pos2 {
	return &Pos2{wayPos.N, -wayPos.E}
}

func navigate2(wayPos *Pos2, pos *Pos2, instructions []string) (*Pos2, *Pos2, error) {
	// If no more instructions
	if len(instructions) == 0 {
		return wayPos, pos, nil
	}

	instruction := instructions[0]

	letter, value, err := parseLine(instruction)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing instruction %s : %w", instruction, err)
	}

	switch letter {
	case 'N':
		wayPos.N = wayPos.N + value
	case 'S':
		wayPos.N = wayPos.N - value
	case 'E':
		wayPos.E = wayPos.E + value
	case 'W':
		wayPos.E = wayPos.E - value
	case 'R':
		nb_turn := value / 90
		for i := 0; i < nb_turn; i++ {
			wayPos = turnRight2(wayPos)
		}
	case 'L':
		value = -value
		for value < 0 {
			value = value + 360
		}
		nb_turn := value / 90
		for i := 0; i < nb_turn; i++ {
			wayPos = turnRight2(wayPos)
		}
	case 'F':
		pos.E = pos.E + value*wayPos.E
		pos.N = pos.N + value*wayPos.N
	}

	return navigate2(wayPos, pos, instructions[1:])
}
