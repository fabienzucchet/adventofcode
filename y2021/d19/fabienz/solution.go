package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	scanners, err := scannersFromLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	beacons, _ := findBeaconsAbsCoordinates(scanners)

	_, err = fmt.Fprintf(answer, "%d", countBeacons(beacons))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	scanners, err := scannersFromLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	_, scanners = findBeaconsAbsCoordinates(scanners)

	_, err = fmt.Fprintf(answer, "%d", maxManhattanDistance(scanners))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Vector struct {
	x, y, z int
}

func (v Vector) String() string {
	return fmt.Sprintf("(%d %d %d)", v.x, v.y, v.z)
}

type Matrix [3]Vector

func (m Matrix) String() string {
	return fmt.Sprintf("\n|%d %d %d|\n|%d %d %d|\n|%d %d %d|", m[0].x, m[1].x, m[2].x, m[0].y, m[1].y, m[2].y, m[0].z, m[1].z, m[2].z)
}

type Scanner struct {
	id         int
	reoriented bool
	pos        Vector
	beacons    []Vector
}

type absBeaconCoor map[Vector]bool

// INPUT PARSING

var scanRegex = regexp.MustCompile(`^--- scanner ([0-9]+) ---$`)
var beaconRegex = regexp.MustCompile(`^([0-9-]+),([0-9-]+),([0-9-]+)$`)

// Transform the raw input lines into a slice of scanners.
func scannersFromLines(lines []string) (scanners []Scanner, err error) {

	var scannerId int

	for _, line := range lines {
		switch {
		// If the line is a scanner header, we parse the scanner id
		case scanRegex.MatchString(line):
			match := scanRegex.FindStringSubmatch(line)
			if len(match) < 2 {
				return scanners, fmt.Errorf("could not parse scanner id in line %s", line)
			}

			id, err := strconv.Atoi(match[1])
			if err != nil {
				return scanners, fmt.Errorf("could not parse scanner id in %s : %w", match[1], err)
			}

			scannerId = id

			scanners = append(scanners, Scanner{id: scannerId})

		// If the line contains the position of a beacon, we parse the position and add it to the corresponding scanner
		case beaconRegex.MatchString(line):
			match := beaconRegex.FindStringSubmatch(line)
			if len(match) < 4 {
				return scanners, fmt.Errorf("could not parse beacon positions in line %s", line)
			}

			x, err := strconv.Atoi(match[1])
			if err != nil {
				return scanners, fmt.Errorf("error parsing x coordinates in %s : %w", match[1], err)
			}

			y, err := strconv.Atoi(match[2])
			if err != nil {
				return scanners, fmt.Errorf("error parsing y coordinates in %s : %w", match[2], err)
			}

			z, err := strconv.Atoi(match[3])
			if err != nil {
				return scanners, fmt.Errorf("error parsing z coordinates in %s : %w", match[3], err)
			}

			scanners[scannerId].beacons = append(scanners[scannerId].beacons, Vector{x, y, z})

		}
	}

	return scanners, nil
}

// VECTORS OPERATIONS

// Scalar product of two vectors
func (v1 *Vector) scalarProduct(v2 Vector) int {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

// Apply a matrix product to a Vector
func (v *Vector) matrixProduct(m Matrix) {
	var w Vector

	w.x = v.scalarProduct(m[0])
	w.y = v.scalarProduct(m[1])
	w.z = v.scalarProduct(m[2])
	*v = w
}

// Create a copy of a vector
func (v Vector) copy() (newV Vector) {

	newV.x = v.x
	newV.y = v.y
	newV.z = v.z

	return newV
}

// Substract two vectors
func sub(v1, v2 Vector) (sub Vector) {

	sub.x = v1.x - v2.x
	sub.y = v1.y - v2.y
	sub.z = v1.z - v2.z

	return sub
}

// Sum two vectors
func add(v1, v2 Vector) (sum Vector) {

	sum.x = v1.x + v2.x
	sum.y = v1.y + v2.y
	sum.z = v1.z + v2.z

	return sum
}

// Compute the manhattant distance between two vectors
func manhattanDistance(v1, v2 Vector) int {
	return abs(v2.x-v1.x) + abs(v2.y-v1.y) + abs(v2.z-v1.z)
}

// ROTATIONS

var Rx = Matrix{
	{1, 0, 0},
	{0, 0, 1},
	{0, -1, 0},
}

var Ry = Matrix{
	{0, 0, -1},
	{0, 1, 0},
	{1, 0, 0},
}

var Rz = Matrix{
	{0, 1, 0},
	{-1, 0, 0},
	{0, 0, 1},
}

var matrixMap = map[string]Matrix{
	"x": Rx,
	"y": Ry,
	"z": Rz,
}

// An array containing all the possible 24 orientations of a scanner.
var possibleOrientations = [24][]string{
	{},
	{"x"},
	{"y"},
	{"x", "x"},
	{"x", "y"},
	{"y", "x"},
	{"y", "y"},
	{"x", "x", "x"},
	{"x", "x", "y"},
	{"x", "y", "x"},
	{"x", "y", "y"},
	{"y", "x", "x"},
	{"y", "y", "x"},
	{"y", "y", "y"},
	{"x", "x", "x", "y"},
	{"x", "x", "y", "x"},
	{"x", "x", "y", "y"},
	{"x", "y", "x", "x"},
	{"x", "y", "y", "y"},
	{"y", "x", "x", "x"},
	{"y", "y", "y", "x"},
	{"x", "x", "x", "y", "x"},
	{"x", "y", "x", "x", "x"},
	{"x", "y", "y", "y", "x"},
}

// Rotate n times 90 degres counterclockwise around the x axis
func (v *Vector) rotate(rotAxis []string) {
	for _, axis := range rotAxis {
		v.matrixProduct(matrixMap[axis])
	}
}

// SEARCH FOR SCANNER POSITIONS

// Find the absolute coordinates of all beacons
func findBeaconsAbsCoordinates(scanners []Scanner) (coor absBeaconCoor, updatedScanners []Scanner) {

	coor = make(absBeaconCoor)

	// Scanner 0 defines absolute coordinates so we can add beacons seen by scanner 0 to the map.
	scanners[0].reoriented = true
	for _, beaconCoor := range scanners[0].beacons {
		coor[beaconCoor] = true
	}

	for !areAllOriented(scanners) {
		for i, referenceScanner := range scanners {
			if referenceScanner.reoriented {
				// Find the direction and coordinates of other scanners to determine the absolute coordinates of the beacons they detect
				for j, scanner := range scanners {
					if !scanner.reoriented && j != i {
						// For each scanner, we try every orientation until we find the absolute coordinates of the scanner
						for _, orientation := range possibleOrientations {
							var rotatedBeacons []Vector

							// Rotate all beacons
							for _, beacon := range scanner.beacons {
								rotatedBeacon := beacon.copy()
								rotatedBeacon.rotate(orientation)
								rotatedBeacons = append(rotatedBeacons, rotatedBeacon)
							}

							// See if a translation vector is the same for several beacons
							translationsMap := make(map[Vector]int)

							for _, rotatedBeacon := range rotatedBeacons {
								for _, refBeacon := range referenceScanner.beacons {
									translationsMap[sub(add(referenceScanner.pos, refBeacon), rotatedBeacon)]++
								}
							}

							for transVector, count := range translationsMap {
								// Count >= 12 means that we identified enough beacons detected by both scanners to find the scanner's absolute coordinates
								if count >= 12 {
									scanner.pos = transVector
									// Orient correctly the scanner
									scanner.beacons = rotatedBeacons
									scanner.reoriented = true

									scanners[j] = scanner

									// Append the absolute positions of the beacons of the scanner in the global map
									for _, beacon := range rotatedBeacons {
										coor[add(beacon, scanner.pos)] = true
									}

								}
							}
						}
					}
				}
			}
		}
	}

	return coor, scanners
}

// Count the number of distinct beacons
func countBeacons(beacons absBeaconCoor) (count int) {

	for _, isBeacon := range beacons {
		if isBeacon {
			count++
		}
	}

	return count
}

// Check of all scanners are reoriented
func areAllOriented(scanners []Scanner) bool {
	for _, scanner := range scanners {
		if !scanner.reoriented {
			return false
		}
	}

	return true
}

// Find the largest Manhattan distance between any pair of scanners
func maxManhattanDistance(scanners []Scanner) (maxDist int) {

	for _, scanner1 := range scanners {
		for _, scanner2 := range scanners {
			dist := manhattanDistance(scanner1.pos, scanner2.pos)
			if dist > maxDist {
				maxDist = dist
			}
		}
	}

	return maxDist
}

// HELPER FUNCTIONS

// return the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
