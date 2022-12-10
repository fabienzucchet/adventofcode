package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 17 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	target, err := parseLine(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	var maxY int

	for dx := 1; dx <= target.maxX; dx++ {
		for dy := target.minY; dy <= -target.minY; dy++ {
			probe, success := launch(dx, dy, target)
			if success && probe.maxY > maxY {
				maxY = probe.maxY
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", maxY)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 17 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	target, err := parseLine(lines[0])
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	var count int

	for dx := 1; dx <= target.maxX; dx++ {
		for dy := target.minY; dy <= -target.minY; dy++ {
			_, success := launch(dx, dy, target)
			if success {
				count++
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Target struct {
	minX, maxX int
	minY, maxY int
}

type Probe struct {
	x, y   int
	dx, dy int
	maxY   int
}

// INPUT PARSING

// Use regex to parse the input line
var re = regexp.MustCompile(`^target area: x=(-?[0-9]+)\.\.(-?[0-9]+), y=(-?[0-9]+)..(-?[0-9]+)$`)

func parseLine(line string) (target Target, err error) {

	match := re.FindStringSubmatch(line)
	if len(match) < 5 {
		return target, fmt.Errorf("error parsing input line %s", line)
	}

	minX, err := strconv.Atoi(match[1])
	if err != nil {
		return target, fmt.Errorf("error parsing min X value %s : %w", match[1], err)
	}

	maxX, err := strconv.Atoi(match[2])
	if err != nil {
		return target, fmt.Errorf("error parsing max X value %s : %w", match[2], err)
	}

	minY, err := strconv.Atoi(match[3])
	if err != nil {
		return target, fmt.Errorf("error parsing min Y value %s : %w", match[3], err)
	}

	maxY, err := strconv.Atoi(match[4])
	if err != nil {
		return target, fmt.Errorf("error parsing max Y value %s : %w", match[4], err)
	}

	target.minX = minX
	target.maxX = maxX
	target.minY = minY
	target.maxY = maxY

	return target, nil
}

// Check if a probe is in a target
func isInTarget(probe Probe, target Target) bool {
	return probe.x >= target.minX && probe.x <= target.maxX && probe.y >= target.minY && probe.y <= target.maxY
}

// Simulate a probe launch
func launch(dx, dy int, target Target) (finalProbe Probe, success bool) {

	probe := Probe{0, 0, dx, dy, 0}

	for probe.x <= target.maxX && probe.y >= target.minY {
		if isInTarget(probe, target) {
			return probe, true
		}

		// Increment one step
		probe.x += probe.dx
		probe.y += probe.dy

		if probe.dx > 0 {
			probe.dx--
		}
		probe.dy--

		if probe.y > probe.maxY {
			probe.maxY = probe.y
		}
	}

	return probe, false
}
