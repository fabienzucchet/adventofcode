package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 21 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	g.scoreToWin = 1000

	deterministicRoll := func(lastRoll int) (rollResult int) {
		if lastRoll == 100 {
			return 1
		}

		return lastRoll + 1
	}

	d := Dice{rollFunction: deterministicRoll}

	g.dice = d

	g.play()

	_, err = fmt.Fprintf(answer, "%d", g.result())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 21 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	// Create the memoization map
	mem = make(map[State][2]int)

	winCount := playDirac(g.state)

	max := winCount[0]
	if winCount[1] > max {
		max = winCount[1]
	}

	_, err = fmt.Fprintf(answer, "%d", max)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type diceRoll func(lastRoll int) (rollResult int)

type Dice struct {
	iteration    int
	lastRoll     int
	rollFunction diceRoll
}

type Player struct {
	pos   int
	score int
}

type State struct {
	players       [2]Player
	currentPlayer int
}

type Game struct {
	state      State
	scoreToWin int
	dice       Dice
}

func (g Game) String() string {
	return fmt.Sprintf("Scores : (J1: %d, J2: %d), Positions : (J1: %d, J2: %d), next player is %d", g.state.players[0].score, g.state.players[1].score, g.state.players[0].pos, g.state.players[1].pos, g.state.currentPlayer)
}

// INPUT PARSING

func parseLines(lines []string) (g Game, err error) {

	if len(lines) < 2 {
		return g, fmt.Errorf("invalid input : only %d lines found", len(lines))
	}

	pos1, err := strconv.Atoi(string(lines[0][len(lines[0])-1]))
	if err != nil {
		return g, fmt.Errorf("error parsing starting position of player 1 : %w", err)
	}

	pos2, err := strconv.Atoi(string(lines[1][len(lines[1])-1]))
	if err != nil {
		return g, fmt.Errorf("error parsing starting position of player 2 : %w", err)
	}

	g.state.players[0].pos = pos1
	g.state.players[1].pos = pos2
	g.state.currentPlayer = 1

	return g, nil
}

// Play a game
func (g *Game) play() {
	for g.state.players[0].score < g.scoreToWin && g.state.players[1].score < g.scoreToWin {

		// A player rolls the dice 3 times and sum the result
		var rollSum int
		for i := 0; i < 3; i++ {
			rollResult := g.dice.rollFunction(g.dice.lastRoll)
			g.dice.iteration++
			g.dice.lastRoll = rollResult
			rollSum += rollResult
		}

		switch g.state.currentPlayer {
		case 1:
			g.state.players[0].pos += rollSum
			for g.state.players[0].pos > 10 {
				g.state.players[0].pos -= 10
			}

			g.state.players[0].score += g.state.players[0].pos
			g.state.currentPlayer = 2

		case 2:
			g.state.players[1].pos += rollSum
			for g.state.players[1].pos > 10 {
				g.state.players[1].pos -= 10
			}

			g.state.players[1].score += g.state.players[1].pos

			g.state.currentPlayer = 1
		}
	}
}

// Compute the result of a game
func (g *Game) result() (result int) {
	if g.state.players[0].score >= g.scoreToWin {
		return g.state.players[1].score * g.dice.iteration
	}

	if g.state.players[1].score >= g.scoreToWin {
		return g.state.players[0].score * g.dice.iteration
	}

	return -1
}

// mem will keep track of the already computed results
var mem map[State][2]int

// This array stores the cardinality of the set of way to make a given sum with 3 dirac dice
var diceCardinality = [10]int{0, 0, 0, 1, 3, 6, 7, 6, 3, 1}

// Play with the Dirac die
func playDirac(s State) (winCount [2]int) {

	// If we already have the solution in memory
	if wCount, present := mem[s]; present {
		return wCount
	}

	// Otherwise we compute recursively the solution
	// A sum of dirac dice can only be between 3 and 9
	for sum := 3; sum <= 9; sum++ {
		copyState := copy(s)

		copyState.players[copyState.currentPlayer-1].pos = (copyState.players[copyState.currentPlayer-1].pos+sum-1)%10 + 1
		copyState.players[copyState.currentPlayer-1].score += copyState.players[copyState.currentPlayer-1].pos

		// If the player won, we can increment its winCount
		if copyState.players[copyState.currentPlayer-1].score >= 21 {
			winCount[copyState.currentPlayer-1] += diceCardinality[sum]
		} else {
			// Otherwise we find the result recursively
			copyState.currentPlayer = 3 - copyState.currentPlayer
			subWinCount := playDirac(copyState)
			winCount[0] += subWinCount[0] * diceCardinality[sum]
			winCount[1] += subWinCount[1] * diceCardinality[sum]
		}
	}

	// Don't forget to memoize the solution
	mem[s] = winCount

	return winCount
}

// Copy a state
func copy(s State) (newState State) {

	newState.players[0].pos = s.players[0].pos
	newState.players[1].pos = s.players[1].pos
	newState.players[0].score = s.players[0].score
	newState.players[1].score = s.players[1].score
	newState.currentPlayer = s.currentPlayer

	return newState
}
