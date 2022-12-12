package fabienz

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 10 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse input into an asteroid map.
	asteroidMap, err := ParseAsteroidMap(lines)
	if err != nil {
		return fmt.Errorf("could not parse asteroid map: %w", err)
	}

	maxCount := -1

	// For each asteroid, count the number of asteroids that can be seen.
	for coord, isAsteroid := range asteroidMap.Map {
		if !isAsteroid {
			continue
		}

		// Count the number of asteroids that can be seen from the current asteroid.
		count := asteroidMap.countVisibleAsteroids(coord)

		if count > maxCount {
			maxCount = count
		}

	}

	_, err = fmt.Fprintf(answer, "%d", maxCount)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 10 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse input into an asteroid map.
	asteroidMap, err := ParseAsteroidMap(lines)
	if err != nil {
		return fmt.Errorf("could not parse asteroid map: %w", err)
	}

	// Find the asteroid with the most visible asteroids.
	maxCount := -1
	maxCoord := Coord{}

	// For each asteroid, count the number of asteroids that can be seen.
	for coord, isAsteroid := range asteroidMap.Map {
		if !isAsteroid {
			continue
		}

		// Count the number of asteroids that can be seen from the current asteroid.
		count := asteroidMap.countVisibleAsteroids(coord)

		if count > maxCount {
			maxCount = count
			maxCoord = coord
		}
	}

	// Find the 200th asteroid to be vaporized.
	asteroid := asteroidMap.Vaporise(maxCoord, 200)

	_, err = fmt.Fprintf(answer, "%d", asteroid.X*100+asteroid.Y)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Coord struct {
	X int
	Y int
}

// The asteroid map is a map of coordinates to a boolean indicating if there is an asteroid at that coordinate.
type AsteroidMap struct {
	Map    map[Coord]bool
	Width  int
	Height int
}

// Parse the input into an asteroid map.
func ParseAsteroidMap(lines []string) (AsteroidMap, error) {
	Map := make(map[Coord]bool)

	for y, line := range lines {
		for x, char := range line {
			Map[Coord{X: x, Y: y}] = char == '#'
		}
	}

	return AsteroidMap{Map: Map, Width: len(lines[0]), Height: len(lines)}, nil
}

// Check if an asteroid can be seen in a given direction from the coordinate.
func (a AsteroidMap) CanSeeAsteroid(from Coord, asteroid Coord) bool {
	// The asteroid is not visible if an asteroid is on the line from the asteroid to the from coordinate.
	// Compute the direction from the asteroid to the from coordinate.
	direction := Coord{X: from.X - asteroid.X, Y: from.Y - asteroid.Y}

	// Compute the greatest common divisor of the direction.
	gcd := gcd(direction.X, direction.Y)

	// Divide the direction by the greatest common divisor to get the normalized direction.
	normalizedDirection := Coord{X: direction.X / gcd, Y: direction.Y / gcd}

	asteroid.X += normalizedDirection.X
	asteroid.Y += normalizedDirection.Y

	// Check if there is an asteroid on the line from the asteroid to the from coordinate.
	for asteroid.X != from.X || asteroid.Y != from.Y {

		if a.Map[asteroid] {
			return false
		}

		asteroid.X += normalizedDirection.X
		asteroid.Y += normalizedDirection.Y
	}

	return true
}

// Check for each asteroid if it can be seen from the given coordinate.
func (a AsteroidMap) countVisibleAsteroids(from Coord) int {
	count := 0

	for coord, isAsteroid := range a.Map {
		// If the considered position is not an asteroid, skip it.
		if !isAsteroid {
			continue
		}

		// If the considered position is the same as the from coordinate, skip it.
		if coord == from {
			continue
		}

		// If the asteroid can be seen from the from coordinate, increment the count.
		if a.CanSeeAsteroid(from, coord) {
			count++
		}
	}

	return count
}

// Compute the greatest common divisor of two numbers.
func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}

	if b < 0 {
		b = -b
	}

	for b != 0 {
		a, b = b, a%b
	}

	return a
}

// Use Polar coordinates to represent the asteroids. The angle is the angle between the asteroid and the vertical direction pointing up).
type PolarCoord struct {
	Angle float64
	Dist  float64
}

// Compute the polar coordinates of the asteroid from the given coordinate.
func (a AsteroidMap) PolarCoord(from Coord, asteroid Coord) PolarCoord {
	// Compute the angle between the asteroid and the vertical direction pointing up.
	angle := math.Atan2(float64(asteroid.X-from.X), float64(from.Y-asteroid.Y))

	// If the angle is negative, add 2π to it.
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Compute the distance between the asteroid and the from coordinate.
	dist := math.Sqrt(math.Pow(float64(asteroid.X-from.X), 2) + math.Pow(float64(asteroid.Y-from.Y), 2))

	return PolarCoord{Angle: angle, Dist: dist}
}

// Store all the angles in a map. The key is the angle and the value is a list of asteroids with that angle.
type PolarMap map[float64][]Coord

// Vaporise the n first asteroids from a given position.
func (a AsteroidMap) Vaporise(from Coord, n int) Coord {
	// Create a map of polar coordinates.
	polarMap := make(PolarMap)

	// For each asteroid, compute the polar coordinates and add it to the map.
	for coord, isAsteroid := range a.Map {
		if !isAsteroid {
			continue
		}

		if coord == from {
			continue
		}

		polarCoord := a.PolarCoord(from, coord)
		polarMap[polarCoord.Angle] = append(polarMap[polarCoord.Angle], coord)
	}

	// Sort the asteroids by dist for the same angle.
	for _, asteroids := range polarMap {
		sort.Slice(asteroids, func(i, j int) bool {
			return a.PolarCoord(from, asteroids[i]).Dist < a.PolarCoord(from, asteroids[j]).Dist
		})
	}

	// Sort the angles.
	angles := make([]float64, 0, len(polarMap))
	for angle := range polarMap {
		angles = append(angles, angle)
	}

	sort.Float64s(angles)

	// Current angle.
	currentAngle := -1.0

	// Vaporise the asteroids.
	for i := 0; i < n; i++ {
		// Find the first angle with an asteroid.
		var angle float64
		for _, angle = range angles {
			if len(polarMap[angle]) > 0 && angle > currentAngle {
				break
			}
		}

		// Remove the first asteroid from the list of asteroids with that angle.
		asteroid := polarMap[angle][0]
		currentAngle = angle
		polarMap[angle] = polarMap[angle][1:]

		// If there are no more asteroids with that angle, remove the angle from the list of angles.
		if len(polarMap[angle]) == 0 {
			angles = angles[1:]
		}

		// If this is the n-th asteroid to be vaporized, return it.
		if i == n-1 {
			return asteroid
		}
	}

	return Coord{}
}
