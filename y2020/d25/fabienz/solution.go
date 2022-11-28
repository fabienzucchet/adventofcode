package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 25 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	pbk1, pbk2, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", findEncyptionKey(pbk1, pbk2))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Parse the input
func parseLines(lines []string) (pbk1, pbk2 int, err error) {

	if len(lines) < 2 {
		return pbk1, pbk2, fmt.Errorf("wrong input format : found %d lines", len(lines))
	}

	pbk1, err = strconv.Atoi(lines[0])
	if err != nil {
		return pbk1, pbk2, fmt.Errorf("error parsing public key %s : %w", lines[0], err)
	}

	pbk2, err = strconv.Atoi(lines[1])
	if err != nil {
		return pbk1, pbk2, fmt.Errorf("error parsing public key %s : %w", lines[1], err)
	}

	return pbk1, pbk2, nil
}

// transforms a secret number with the loops size
func transform(subjectNumber int, loopSize int) (publicKey int) {

	// Start with a value of 1
	publicKey = 1

	for i := 0; i < loopSize; i++ {
		publicKey = (subjectNumber * publicKey) % 20201227
	}

	return publicKey
}

// Determine the loop size with subject number and public key
func bruteForce(publicKey int) (loopSize int) {

	value := 1

	for value != publicKey {
		value = (7 * value) % 20201227
		loopSize++
	}

	return loopSize
}

// Find the encryption key
func findEncyptionKey(pbk1 int, pbk2 int) (encryptionKey int) {

	loopSize1 := bruteForce(pbk1)

	encryptionKey = transform(pbk2, loopSize1)

	return encryptionKey
}
