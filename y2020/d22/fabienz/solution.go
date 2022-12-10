package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

const BASE = 100 // Base used for the hash of a deck, assume no cards will have a value > 100

// PartOne solves the first problem of day 22 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing input : %w", err)
	}

	g.playGame()

	_, err = fmt.Fprintf(answer, "%d", g.getFinalScore())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 22 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error when parsing input : %w", err)
	}

	g.playRecursiveGame()

	_, err = fmt.Fprintf(answer, "%d", g.getFinalScore())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Types
type game struct {
	p1 pile
	p2 pile
	// Used for part 2 : will store history of every round to compare and prevent infinite recursion
	previousDecks1 []pile
	previousDecks2 []pile
}

type pile []int

// PARSE INPUT

func parseLines(lines []string) (g game, err error) {

	var player string
	for _, line := range lines {
		switch line {
		case "Player 1:":
			player = line
		case "Player 2:":
			player = line
		case "":
		default:
			number, err := strconv.Atoi(line)
			if err != nil {
				return g, fmt.Errorf("error when parsing the line %s : %w", line, err)
			}

			switch player {
			case "Player 1:":
				g.p1 = append(g.p1, number)
			case "Player 2:":
				g.p2 = append(g.p2, number)
			}
		}
	}

	return g, nil
}

// Get the card on top of a deck
func (p *pile) drawCard() (card int) {
	card = (*p)[0]
	*p = (*p)[1:]

	return card
}

// Add card at the bottom of the deck
func (p *pile) addCard(card int) {
	*p = append(*p, card)
}

// Play a round for of the Combat
func (g *game) playRound() {
	// Draw cards
	card1 := g.p1.drawCard()
	card2 := g.p2.drawCard()

	if card1 > card2 {
		g.p1.addCard(card1)
		g.p1.addCard(card2)
	} else if card1 == card2 {
		g.p1.addCard(card1)
		g.p2.addCard(card2)
	} else {
		g.p2.addCard(card2)
		g.p2.addCard(card1)
	}
}

// Play a round of the recursive Combat
func (g *game) playRecursiveRound() {

	// Check if this deck already happenend before
	for idx := range g.previousDecks1 {
		if compareDecks(g.previousDecks1[idx], g.p1) && compareDecks(g.previousDecks2[idx], g.p2) {
			// Player 1 wins instantly
			g.p2 = pile{}
			return
		}
	}

	// If this configuration is new, add the hash to history
	g.previousDecks1 = append(g.previousDecks1, copyDeck(g.p1))
	g.previousDecks2 = append(g.previousDecks2, copyDeck(g.p2))

	// Draw cards
	card1 := g.p1.drawCard()
	card2 := g.p2.drawCard()

	if card1 <= len(g.p1) && card2 <= len(g.p2) {
		// Generate a recursive sub-game
		var subGame game

		// Copy the deck for the recursive subGame
		for i := 0; i < card1; i++ {
			subGame.p1 = append(subGame.p1, g.p1[i])
		}
		for i := 0; i < card2; i++ {
			subGame.p2 = append(subGame.p2, g.p2[i])
		}

		// Play the subgame
		subGame.playRecursiveGame()

		// If player 2 won
		if len(subGame.p1) == 0 {
			g.p2.addCard(card2)
			g.p2.addCard(card1)
		} else {
			g.p1.addCard(card1)
			g.p1.addCard(card2)
		}

		return
	}
	// Traditional round
	if card1 > card2 {
		g.p1.addCard(card1)
		g.p1.addCard(card2)
	} else if card1 == card2 {
		g.p1.addCard(card1)
		g.p2.addCard(card2)
	} else {
		g.p2.addCard(card2)
		g.p2.addCard(card1)
	}
}

// Play a game
func (g *game) playGame() {
	for !g.hasWon() {
		g.playRound()
	}
}

// Play a recursive game
func (g *game) playRecursiveGame() {
	for !g.hasWon() {
		g.playRecursiveRound()
	}
}

// Check if a player has won
func (g *game) hasWon() bool {
	return len(g.p1) == 0 || len(g.p2) == 0
}

// Compute the score of a deck
func (p pile) getScore() (score int) {

	n := len(p)

	for idx, number := range p {
		score += (n - idx) * number
	}

	return score
}

// Compute the score of the winner, returns 0 if no player has won
func (g *game) getFinalScore() (score int) {
	if !g.hasWon() {
		return score
	}
	if len(g.p1) == 0 {
		return g.p2.getScore()
	}

	return g.p1.getScore()
}

// Compare two decks
func compareDecks(p1, p2 pile) bool {
	if len(p1) != len(p2) {
		return false
	}

	for idx := range p1 {
		if p1[idx] != p2[idx] {
			return false
		}
	}

	return true
}

// Copy a deck
func copyDeck(p pile) (copy pile) {
	for _, number := range p {
		copy = append(copy, number)
	}
	return copy
}
