package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
	"github.com/fabienzucchet/adventofcode/y2019/opcode"
)

// PartOne solves the first problem of day 13 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Init the game
	game, err := initGame(lines[0], 0)
	if err != nil {
		return fmt.Errorf("could not init game: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", game.countBlockTiles())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Init the game
	game, err := initGame(lines[0], 2)
	if err != nil {
		return fmt.Errorf("could not init game: %w", err)
	}

	// Play the game while there are still block tiles: move the joystick to make the ball bounce on the paddle
	for {
		// Run the program
		_, err := game.software.RunIntcode()
		if err != nil {
			return fmt.Errorf("could not run intcode: %w", err)
		}

		// Update the game board
		game.updateBoard()

		// Check if there are still block tiles
		if game.countBlockTiles() == 0 {
			break
		}

		// Move the joystick
		game.moveJoystick()
	}

	_, err = fmt.Fprintf(answer, "%d", game.score)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Game struct {
	board    helpers.IntGrid2D
	software opcode.Intcode
	score    int
}

// Init the game by running the intcode program
func initGame(program string, quarters int) (game Game, err error) {
	instructions, err := helpers.IntsFromString(program, ",")
	if err != nil {
		return game, fmt.Errorf("could not parse intcode: %w", err)
	}

	// Init the intcode program
	game.software = opcode.Intcode{Program: instructions}

	// Insert the number of quarters
	if quarters != 0 {
		game.software.Program[0] = quarters
	}

	// Run the program
	outputs, err := game.software.RunIntcode()
	if err != nil {
		return game, fmt.Errorf("could not run intcode: %w", err)
	}

	// Check that the number of outputs is a multiple of 3
	if len(outputs)%3 != 0 {
		return game, fmt.Errorf("unexpected number of outputs: %d", len(outputs))
	}

	game.board = make(helpers.IntGrid2D)

	// Update the board
	err = game.updateBoard()
	if err != nil {
		return game, fmt.Errorf("could not update board: %w", err)
	}

	return game, nil
}

// Parse a board
func (g *Game) updateBoard() error {
	// Check that the number of outputs is a multiple of 3
	if len(g.software.Outputs)%3 != 0 {
		return fmt.Errorf("unexpected number of outputs: %d", len(g.software.Outputs))
	}

	// Parse the output into the grid
	for i := 0; i < len(g.software.Outputs); i += 3 {
		// If X=-1 and Y=0, the third output instruction is not a tile; the value instead specifies the new score to show in the segment display.
		if g.software.Outputs[i] == -1 && g.software.Outputs[i+1] == 0 {
			g.score = g.software.Outputs[i+2]
			continue
		}
		g.board[helpers.Coord2D{X: g.software.Outputs[i], Y: g.software.Outputs[i+1]}] = g.software.Outputs[i+2]
	}

	// Reset the outputs
	g.software.Outputs = []int{}

	return nil
}

// Count the number of block tiles
func (g *Game) countBlockTiles() int {
	var nbBlockTiles int
	for _, tile := range g.board {
		if tile == 2 {
			nbBlockTiles++
		}
	}
	return nbBlockTiles
}

// Move the joystick
func (g *Game) moveJoystick() {
	// Find the ball and the paddle
	var ballX, paddleX int
	for coord, tile := range g.board {
		if tile == 4 {
			ballX = coord.X
		}
		if tile == 3 {
			paddleX = coord.X
		}
	}

	// Move the joystick
	if ballX > paddleX {
		g.software.Inputs = append(g.software.Inputs, 1)
	} else if ballX < paddleX {
		g.software.Inputs = append(g.software.Inputs, -1)
	} else {
		g.software.Inputs = append(g.software.Inputs, 0)
	}
}
