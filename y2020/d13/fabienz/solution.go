package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	timestamp, err := strconv.Atoi(lines[0])

	buses, err := parseBuses(lines[1])
	if err != nil {
		return fmt.Errorf("error parsing buses %s : %w", lines[1], err)
	}

	best_bus := buses[0]
	min_wait := waitBus(timestamp, buses[0])

	for _, bus := range buses {
		wait := waitBus(timestamp, bus)

		if wait < min_wait {
			min_wait = wait
			best_bus = bus
		}
	}

	_, err = fmt.Fprintf(answer, "%d", best_bus*min_wait)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	buses, indexes, err := parseBuses2(lines[1])
	if err != nil {
		return fmt.Errorf("error parsing buses %s : %w", lines[1], err)
	}

	_, err = fmt.Fprintf(answer, "%d", findClosestTimestamp(findIdxOfMax(buses), buses, indexes))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseBuses(input string) ([]int, error) {
	var buses []int

	for _, bus := range strings.Split(input, ",") {
		if bus != "x" {
			id, err := strconv.Atoi(bus)
			if err != nil {
				return nil, fmt.Errorf("error parsing bus %s : %w", bus, err)
			}
			buses = append(buses, id)
		}
	}

	return buses, nil
}

func waitBus(timestamp int, bus int) int {
	depart := 0

	for depart < timestamp {
		depart = depart + bus
	}

	return depart - timestamp
}

func parseBuses2(input string) ([]int, []int, error) {
	var buses []int
	var indexes []int

	for idx, bus := range strings.Split(input, ",") {
		if bus != "x" {
			id, err := strconv.Atoi(bus)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing bus %s : %w", bus, err)
			}
			buses = append(buses, id)
			indexes = append(indexes, idx)
		}
	}

	return buses, indexes, nil
}

func checkRequirements(timestamp int, buses []int, indexes []int) bool {
	fmt.Println("Checking timestamp :", timestamp)

	for i := 0; i < len(buses); i++ {
		if buses[i] != -1 && (timestamp+indexes[i])%buses[i] != 0 {
			return false
		}
	}

	return true
}

func findIdxOfMax(l []int) int {
	idx := 0
	max := l[0]

	for i := 1; i < len(l); i++ {
		if max < l[i] {
			idx = i
			max = l[i]
		}
	}

	return idx
}

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a int, b int) int {
	if a == 0 || b == 0 {
		return 0
	}

	if a >= b {
		return a * b / gcd(a, b)
	}

	return a * b / gcd(b, a)
}

func findClosestTimestamp(max_idx int, buses []int, indexes []int) int {

	timestamp := 0

	delta := 1

	for idx, bus := range buses {
		for (timestamp+indexes[idx])%bus != 0 {
			timestamp += delta
		}
		delta = lcm(delta, bus)
	}

	return timestamp
}
