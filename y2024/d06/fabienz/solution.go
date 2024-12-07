package fabienz

import (
	"fmt"
	"io"
	"log"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 6 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	lab := labFromLines(lines)

	lab.moveGuardUntilLeavingLabOrLoop()

	count := lab.countVisitedCells()

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 6 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	lab := labFromLines(lines)
	visitedPosWithoutExtraObstacle := lab.getAllVisitedPosWithoutExtraObstacle()

	log.Println(visitedPosWithoutExtraObstacle)

	count := lab.countInfiniteLoopsAfterAddingOneObstacle(visitedPosWithoutExtraObstacle)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Coord struct {
	X, Y int
}

type Direction int

const (
	Up    Direction = iota
	Down  Direction = iota
	Left  Direction = iota
	Right Direction = iota
)

type Lab struct {
	LabWidth, LabHeight int
	GuardPos            Coord
	Dir                 Direction
	LabMap              Grid
	Visited             PositionHistory
}

type Grid map[Coord]bool

type PositionHistory map[Coord]Direction

func labFromLines(lines []string) Lab {
	labMap := make(Grid)
	visited := make(PositionHistory)
	guardPos := Coord{0, 0}
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				labMap[Coord{x, y}] = true
			}

			if char == '^' {
				guardPos = Coord{x, y}
				visited[Coord{x, y}] = Up
			}
		}
	}

	return Lab{
		LabWidth:  len(lines[0]),
		LabHeight: len(lines),
		GuardPos:  guardPos,
		Dir:       Up,
		LabMap:    labMap,
		Visited:   visited,
	}
}

func (l *Lab) String() string {
	str := "\n"
	for y := 0; y < l.LabHeight; y++ {
		for x := 0; x < l.LabWidth; x++ {
			if l.GuardPos.X == x && l.GuardPos.Y == y {
				str += "^"
			} else if _, ok := l.LabMap[Coord{x, y}]; ok {
				str += "#"
			} else if _, ok := l.Visited[Coord{x, y}]; ok {
				str += "."
			} else {
				str += " "
			}
		}
		str += "\n"
	}

	return str
}

func (l *Lab) canGuardMoveForward() bool {
	switch l.Dir {
	case Up:
		_, ok := l.LabMap[Coord{l.GuardPos.X, l.GuardPos.Y - 1}]

		return !ok
	case Down:
		_, ok := l.LabMap[Coord{l.GuardPos.X, l.GuardPos.Y + 1}]

		return !ok
	case Left:
		_, ok := l.LabMap[Coord{l.GuardPos.X - 1, l.GuardPos.Y}]

		return !ok
	case Right:
		_, ok := l.LabMap[Coord{l.GuardPos.X + 1, l.GuardPos.Y}]

		return !ok
	}

	return false
}

func (l *Lab) isGuardStillInTheLab() bool {
	return l.GuardPos.X >= 0 && l.GuardPos.X < l.LabWidth && l.GuardPos.Y >= 0 && l.GuardPos.Y < l.LabHeight
}

func (l *Lab) moveGuardForward() {
	switch l.Dir {
	case Up:
		l.GuardPos.Y--
	case Down:
		l.GuardPos.Y++
	case Left:
		l.GuardPos.X--
	case Right:
		l.GuardPos.X++
	}

	if l.isGuardStillInTheLab() {
		l.Visited[l.GuardPos] = l.Dir
	}
}

func (l *Lab) turnGuardRight() {
	switch l.Dir {
	case Up:
		l.Dir = Right
	case Down:
		l.Dir = Left
	case Left:
		l.Dir = Up
	case Right:
		l.Dir = Down
	}
}

// Returns true if the guard is still in the lab (loop), false otherwise.
func (l *Lab) moveGuardUntilLeavingLabOrLoop() bool {
	for l.isGuardStillInTheLab() {
		if l.isGuardInLoop() {
			return true
		}

		if l.canGuardMoveForward() {
			l.moveGuardForward()
		} else {
			l.turnGuardRight()
		}
	}

	return false
}

func (l *Lab) countVisitedCells() int {
	return len(l.Visited)
}

// A guard is in an infinite loop if the next position has already been visited with the same direction.
func (l *Lab) isGuardInLoop() bool {
	nextPos := l.GuardPos
	switch l.Dir {
	case Up:
		nextPos.Y--
	case Down:
		nextPos.Y++
	case Left:
		nextPos.X--
	case Right:
		nextPos.X++
	}

	if dir, ok := l.Visited[nextPos]; ok {
		return dir == l.Dir
	}

	return false
}

// Create a deep copy of the lab.
func (l *Lab) copyLab() Lab {
	newLab := Lab{
		LabWidth:  l.LabWidth,
		LabHeight: l.LabHeight,
		GuardPos:  l.GuardPos,
		Dir:       l.Dir,
		LabMap:    make(Grid),
		Visited:   make(PositionHistory),
	}

	for k, v := range l.LabMap {
		newLab.LabMap[k] = v
	}

	for k, v := range l.Visited {
		newLab.Visited[k] = v
	}

	return newLab
}

func (l *Lab) addObstacle(pos Coord) {
	l.LabMap[pos] = true
}

func (l *Lab) getAllVisitedPosWithoutExtraObstacle() []Coord {
	visitedPos := make([]Coord, 0, len(l.Visited))

	labWithoutExtraObstacle := l.copyLab()
	labWithoutExtraObstacle.moveGuardUntilLeavingLabOrLoop()

	for pos := range labWithoutExtraObstacle.Visited {
		visitedPos = append(visitedPos, pos)
	}

	return visitedPos
}

func (l *Lab) countInfiniteLoopsAfterAddingOneObstacle(positionsWhereAnObstacleCanBeAdded []Coord) int {
	count := 0
	for _, pos := range positionsWhereAnObstacleCanBeAdded {
		newLab := l.copyLab()
		newLab.addObstacle(pos)

		if newLab.moveGuardUntilLeavingLabOrLoop() {
			count++
		}
	}

	return count
}
