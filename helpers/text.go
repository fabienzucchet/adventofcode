package helpers

import (
	"bufio"
	"fmt"
	"io"
)

// LinesFromReader returns a slice of all lines in r.
func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("failed to scan reader: %w", s.Err())
	}

	return lines, nil
}

// Convert a rune in an int.
func IntFromRune(r rune) (int, error) {
	switch r {
	case '0':
		return 0, nil
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	default:
		return 0, fmt.Errorf("invalid rune: %c", r)
	}
}
