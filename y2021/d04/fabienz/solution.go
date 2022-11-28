package fabienz

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const BOARDSIDE = 5

// PartOne solves the first problem of day 4 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	draw, boards, _, err := parseInput(lines)
	if err != nil {
		return fmt.Errorf("could not parse input : %w", err)
	}

	score := -1

	// Let's play the Bingo
	for i := 0; i < len(draw) && score == -1; i++ {
		for _, board := range boards {
			board.mark(draw[i])
			if board.won {
				score = draw[i] * board.sumUnmarked()
				break
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	draw, boards, _, err := parseInput(lines)
	if err != nil {
		return fmt.Errorf("could not parse input : %w", err)
	}

	score := -1
	countWonBoards := 0
	// This array is used to know if it's the first time that the board wins
	hasBoardWon := make([]bool, len(boards))

	// Let's play the Bingo
	for i := 0; i < len(draw) && countWonBoards < len(boards); i++ {
		for idx, board := range boards {
			board.mark(draw[i])
			if board.won {
				if !hasBoardWon[idx] {
					countWonBoards++
					hasBoardWon[idx] = true
					if countWonBoards == len(boards) {
						score = draw[i] * board.sumUnmarked()
					}
				}
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Coordinates struct {
	row int
	col int
}

type Board struct {
	numbers map[int]Coordinates
	marks   *[BOARDSIDE][BOARDSIDE]bool
	won     bool
}

// Returns a slice of non-empty lists
func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		if s.Text() != "" {
			lines = append(lines, s.Text())
		}
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("failed to scan reader: %w", s.Err())
	}

	return lines, nil
}

func parseInput(lines []string) (drawnNumbers []int, boards []Board, marks [][BOARDSIDE][BOARDSIDE]bool, err error) {

	// Parse the first line (drawn numbers)
	for _, num := range strings.Split(lines[0], ",") {
		number, err := strconv.Atoi(num)
		if err != nil {
			return nil, nil, nil, err
		}
		drawnNumbers = append(drawnNumbers, number)
	}

	// Parse the boards
	for i := 1; i+BOARDSIDE-1 < len(lines); i += BOARDSIDE {
		numbers := make(map[int]Coordinates)
		for j := 0; j < BOARDSIDE; j++ {
			for k, num := range strings.Fields(lines[i+j]) {
				number, err := strconv.Atoi(num)
				if err != nil {
					return nil, nil, nil, err
				}

				numbers[number] = Coordinates{row: j, col: k}
			}
		}

		var mark [BOARDSIDE][BOARDSIDE]bool
		marks = append(marks, mark)
		boards = append(boards, Board{numbers: numbers, marks: &mark})
	}

	return drawnNumbers, boards, marks, nil
}

// Marks a given number on the board, does nothing if the number is not in the board.
func (b *Board) mark(draw int) {
	if coor, present := b.numbers[draw]; present {
		b.marks[coor.row][coor.col] = true

		// We check if the board has won every time we draw a number because we only have to check 1 row and 1 column

		isWinningCol := true
		for i := 0; i < BOARDSIDE; i++ {
			if !b.marks[i][coor.col] {
				isWinningCol = false
			}
		}
		isWinningRow := true
		for i := 0; i < BOARDSIDE; i++ {
			if !b.marks[coor.row][i] {
				isWinningRow = false
			}
		}

		b.won = isWinningCol || isWinningRow
	}
}

func (b *Board) sumUnmarked() (sum int) {
	for num, coor := range b.numbers {
		if !b.marks[coor.row][coor.col] {
			sum += num
		}
	}

	return sum
}

// This functions prints a board for debug purpose
func (b *Board) printBoard() {

	var matrix [BOARDSIDE][BOARDSIDE]string

	for num, coor := range b.numbers {
		cell := fmt.Sprint(num)

		if b.marks[coor.row][coor.col] {
			cell = "(" + cell + ")"
		}

		matrix[coor.row][coor.col] = cell
	}

	for i := 0; i < BOARDSIDE; i++ {
		for j := 0; j < BOARDSIDE; j++ {
			fmt.Print(matrix[i][j], "     ")
		}
		fmt.Print("\n")
	}

}
