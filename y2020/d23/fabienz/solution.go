package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 23 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseInput(lines[0], len(lines[0]))
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	g.playGame(100)

	_, err = fmt.Fprintf(answer, "%s", g.generateCupOrder())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 23 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	g, err := parseInput(lines[0], 1e6)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	g.playGame(1e7)

	cup1, cup2 := g.findTwoCupsNext1()

	_, err = fmt.Fprintf(answer, "%d", cup1*cup2)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type cup struct {
	label int
	next  int
}

type game struct {
	cups       []cup
	currentCup int
	nbCups     int
}

// INPUT PARSING

// Parses the input and creates the cups. Returns an array of cups such that cups[i] corresponds to the cup labelled i+1
func parseInput(line string, nbCups int) (g game, err error) {

	cups := make([]cup, nbCups)

	var orderedLabels []int

	// parse the line in slice of labels
	for _, char := range line {
		label, err := strconv.Atoi(string(char))
		if err != nil {
			return g, fmt.Errorf("error parsing label %s : %w", string(char), err)
		}
		orderedLabels = append(orderedLabels, label)
	}

	// Add more cups form len(orderedLabels) + 1 to nbCups
	for i := len(orderedLabels) + 1; i <= nbCups; i++ {
		orderedLabels = append(orderedLabels, i)
	}

	// Create the cups
	for i := 0; i < nbCups; i++ {
		label := orderedLabels[labelToIdx(i, nbCups)]
		cups[labelToIdx(label, nbCups)] = cup{label: label, next: orderedLabels[i]}
	}

	g.cups = cups
	g.currentCup = orderedLabels[0]
	g.nbCups = nbCups

	return g, nil
}

// Play the game (n moves)
func (g *game) playGame(n int) {
	for i := 0; i < n; i++ {
		g.playMove()
	}
}

// Play a move
func (g *game) playMove() {

	// Remove 3 cups from the game
	removedLabels := g.pickCups()

	// Find the target cup
	targetLabel := g.currentCup - 1
	if targetLabel == 0 {
		targetLabel = g.nbCups
	}
	for contains(targetLabel, removedLabels) {
		targetLabel = targetLabel - 1
		if targetLabel == 0 {
			targetLabel = g.nbCups
		}
	}

	// Add the cups again after the target
	g.addCups(targetLabel, removedLabels)

	// Update the current cup
	g.currentCup = g.cups[labelToIdx(g.currentCup, g.nbCups)].next
}

// Check if the list of labels contains a given label
func contains(label int, labels []int) bool {

	for _, lab := range labels {
		if lab == label {
			return true
		}
	}

	return false
}

// Remove the 3 cups next to the current cup
func (g *game) pickCups() (removedLabels []int) {

	// We need to remember the first removed cup to reintegrate it later in the circle
	removedLabels = append(removedLabels, g.cups[labelToIdx(g.currentCup, g.nbCups)].next)
	removedLabels = append(removedLabels, g.cups[labelToIdx(removedLabels[len(removedLabels)-1], g.nbCups)].next)
	removedLabels = append(removedLabels, g.cups[labelToIdx(removedLabels[len(removedLabels)-1], g.nbCups)].next)

	// Change the links between cups to remove a 3-chain
	g.cups[labelToIdx(g.currentCup, g.nbCups)].next = g.cups[labelToIdx(removedLabels[len(removedLabels)-1], g.nbCups)].next
	g.cups[labelToIdx(removedLabels[len(removedLabels)-1], g.nbCups)].next = -1

	return removedLabels
}

// Add the removed cups after the target
func (g *game) addCups(target int, removedLabels []int) {
	g.cups[labelToIdx(removedLabels[len(removedLabels)-1], g.nbCups)].next = g.cups[labelToIdx(target, g.nbCups)].next
	g.cups[labelToIdx(target, g.nbCups)].next = removedLabels[0]
}

// Convert label of a cub to its index in the array of cups
func labelToIdx(label int, offset int) (idx int) {
	idx = label - 1

	for idx < 0 {
		idx += offset
	}

	return idx
}

// Generate the cups order
func (g *game) generateCupOrder() (order string) {

	cup := g.cups[labelToIdx(1, g.nbCups)]

	for i := 1; i < g.nbCups; i++ {
		cup = g.cups[labelToIdx(cup.next, g.nbCups)]
		order += strconv.Itoa(cup.label)
	}

	return order
}

// Find the two cups directly after 1
func (g *game) findTwoCupsNext1() (label1, label2 int) {
	cup := g.cups[labelToIdx(1, g.nbCups)]

	cup1 := g.cups[labelToIdx(cup.next, g.nbCups)]
	cup2 := g.cups[labelToIdx(cup1.next, g.nbCups)]

	return cup1.label, cup2.label
}
