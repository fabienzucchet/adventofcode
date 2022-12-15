package fabienz

import (
	"fmt"
	"io"
	"math"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// const SEACHROW = 10
// const MAXX_SEARCHZONE = 20

const SEACHROW = 2000000
const MAXX_SEARCHZONE = 4000000

// PartOne solves the first problem of day 15 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input.
	sensors, beacons, _, _, _, _, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", checkRow(sensors, beacons, SEACHROW, math.MinInt, math.MaxInt))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 15 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse input.
	sensors, beacons, _, _, _, _, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Find the first possible solution
	pos, err := findFirstPossiblePosition(sensors, beacons, 0, 0, MAXX_SEARCHZONE, MAXX_SEARCHZONE)
	if err != nil {
		return fmt.Errorf("could not find first possible position: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", computeTuningFrequency(pos))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type sensor struct {
	pos                   helpers.Coord2D
	closestBeaconDistance int
}

type beacon struct {
	pos helpers.Coord2D
}

// Parse the input.
func parseLines(lines []string) ([]sensor, []beacon, int, int, int, int, error) {
	sensors := make([]sensor, len(lines))
	beacons := make([]beacon, 0)
	// Determine the size of the grid.
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for i, line := range lines {
		var sensX, sensY, beaconX, beaconY int

		_, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensX, &sensY, &beaconX, &beaconY)
		if err != nil {
			return nil, nil, 0, 0, 0, 0, fmt.Errorf("could not parse line %q: %w", line, err)
		}

		sensorPos := helpers.Coord2D{
			X: sensX,
			Y: sensY,
		}
		beaconPos := helpers.Coord2D{
			X: beaconX,
			Y: beaconY,
		}
		closestBeaconDistance := sensorPos.ManhattanDistance(beaconPos)

		sensors[i] = sensor{
			pos:                   sensorPos,
			closestBeaconDistance: closestBeaconDistance,
		}

		// Check if the beacon isn't already in the list.
		found := false
		for _, b := range beacons {
			if b.pos == beaconPos {
				found = true
				break
			}
		}
		if !found {
			beacons = append(beacons, beacon{
				pos: beaconPos,
			})
		}

		// Update the grid size.
		if sensX-closestBeaconDistance < minX {
			minX = sensX
		}
		if sensY-closestBeaconDistance < minY {
			minY = sensY
		}
		if sensX+closestBeaconDistance > maxX {
			maxX = sensX
		}
		if sensY+closestBeaconDistance > maxY {
			maxY = sensY
		}
	}
	return sensors, beacons, minX, minY, maxX, maxY, nil
}

// Check a given row and count the number of position where the beacon cannot be.
func checkRow(sensors []sensor, beacons []beacon, row int, minX int, maxX int) int {
	// Use a map to store the positions where the beacon cannot be.
	posMap := make(map[helpers.Coord2D]bool)

	for _, s := range sensors {
		for i := 0; i <= s.closestBeaconDistance; i++ {
			// Check the left side.
			if s.pos.X-i >= minX && s.pos.X-i <= maxX {
				pos := helpers.Coord2D{
					X: s.pos.X - i,
					Y: row,
				}
				if s.pos.ManhattanDistance(pos) <= s.closestBeaconDistance {
					posMap[pos] = false
				}
			}

			// Do not check the right side if i == 0.
			if i == 0 {
				continue
			}

			if s.pos.X+i >= minX && s.pos.X+i <= maxX {
				pos := helpers.Coord2D{
					X: s.pos.X + i,
					Y: row,
				}
				if s.pos.ManhattanDistance(pos) <= s.closestBeaconDistance {
					posMap[pos] = false
				}
			}
		}
	}

	// Remove the beacons from the map.
	for _, b := range beacons {
		delete(posMap, b.pos)
	}

	count := 0
	for _, v := range posMap {
		if !v {
			count++
		}
	}

	return count
}

// Find the first possible postion for a beacon in a given zone.
func findFirstPossiblePosition(sensors []sensor, beacons []beacon, minX int, minY int, maxX int, maxY int) (helpers.Coord2D, error) {
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			isPointCovered := false
			pos := helpers.Coord2D{
				X: x,
				Y: y,
			}

			// For each sensor, determine if the point is covered.
			for _, s := range sensors {
				dist := s.pos.ManhattanDistance(pos)
				if dist <= s.closestBeaconDistance {
					isPointCovered = true
					// We can skip the next x values that are still covered by the sensor.
					x += s.closestBeaconDistance - dist
					break
				}
			}

			// If we didn't find a sensor covering the point, then it's out candidate.
			if !isPointCovered {
				return pos, nil
			}
		}
	}

	return helpers.Coord2D{}, fmt.Errorf("could not find a position")
}

// Compute a tuning frequency for a given position.
func computeTuningFrequency(pos helpers.Coord2D) int {
	return pos.X*4000000 + pos.Y
}
