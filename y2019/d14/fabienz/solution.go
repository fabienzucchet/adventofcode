package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the requirements
	requirements, err := ParseRequirements(lines)
	if err != nil {
		return fmt.Errorf("could not parse requirements: %w", err)
	}

	nbOreNeeded, err := ComputeOre(requirements, 1)
	if err != nil {
		return fmt.Errorf("could not compute ore: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", nbOreNeeded)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the requirements
	requirements, err := ParseRequirements(lines)
	if err != nil {
		return fmt.Errorf("could not parse requirements: %w", err)
	}

	maxFuel, err := FindMaxFuel(requirements, 1000000000000)
	if err != nil {
		return fmt.Errorf("could not find max fuel: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", maxFuel)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Chemical struct {
	Quantity int
	Name     string
}

type Requirement struct {
	produced  int
	chemicals []Chemical
}

func ParseRequirement(line string) (name string, req Requirement, err error) {
	parts := strings.Split(line, " => ")
	if len(parts) != 2 {
		return "", req, fmt.Errorf("invalid requirement: %s", line)
	}

	// Parse the output
	output, err := ParseChemical(parts[1])
	if err != nil {
		return "", req, fmt.Errorf("could not parse output: %w", err)
	}

	// Parse the inputs
	inputs := strings.Split(parts[0], ", ")
	chemicals := make([]Chemical, len(inputs))
	for i, input := range inputs {
		chemical, err := ParseChemical(input)
		if err != nil {
			return "", req, fmt.Errorf("could not parse input: %w", err)
		}
		chemicals[i] = chemical
	}

	return output.Name, Requirement{produced: output.Quantity, chemicals: chemicals}, nil
}

func ParseChemical(line string) (Chemical, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return Chemical{}, fmt.Errorf("invalid chemical: %s", line)
	}

	quantity, err := strconv.Atoi(parts[0])
	if err != nil {
		return Chemical{}, fmt.Errorf("could not parse quantity: %w", err)
	}

	return Chemical{
		Quantity: quantity,
		Name:     parts[1],
	}, nil
}

// Create a requirements map from the input
func ParseRequirements(lines []string) (map[string]Requirement, error) {
	requirements := make(map[string]Requirement)
	for _, line := range lines {
		name, req, err := ParseRequirement(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse requirement: %w", err)
		}
		requirements[name] = req
	}

	return requirements, nil
}

// compute the quantity of ore needed to produce the given quantity of fuel
func ComputeOre(requirements map[string]Requirement, fuelQuantity int) (int, error) {
	// Map to store the quantity of each chemical we need
	needed := make(map[string]int)
	// Map to store the quantity of each chemical we have in excess
	excess := make(map[string]int)

	// Start with the fuel
	needed["FUEL"] = fuelQuantity

	// While we still need some chemicals other than ore in the needed map
	for len(needed) > 1 || needed["ORE"] == 0 {
		// Get the first chemical we need
		var name string
		for name = range needed {
			if name != "ORE" {
				break
			}
		}

		// Get the quantity we need
		quantity := needed[name]

		// If we have enough in excess, use it
		if excess[name] >= quantity {
			excess[name] -= quantity
			delete(needed, name)
			continue
		}

		// Otherwise, use what we have in excess
		quantity -= excess[name]
		excess[name] = 0

		// Get the requirements for this chemical
		req, ok := requirements[name]
		if !ok {
			return 0, fmt.Errorf("could not find requirements for %s", name)
		}

		// Compute the number of time we need to produce the reaction to get the quantity we need
		reactionCount := (quantity-1)/req.produced + 1
		reqQuantity := reactionCount * req.produced

		// If we have excess, store it
		if reqQuantity > quantity {
			excess[name] = reqQuantity - quantity
		}

		// Add the requirements to the needed map
		for _, chemical := range req.chemicals {
			needed[chemical.Name] += reactionCount * chemical.Quantity
		}

		// Remove the current chemical from the needed map
		delete(needed, name)
	}

	return needed["ORE"], nil
}

// Find the amount of fuel we can produce with the given amount of ore
func FindMaxFuel(requirements map[string]Requirement, oreQuantity int) (int, error) {
	// Binary search
	min := 0
	max := oreQuantity
	for min < max {
		mid := (min + max + 1) / 2
		oreNeeded, err := ComputeOre(requirements, mid)
		if err != nil {
			return 0, fmt.Errorf("could not compute ore: %w", err)
		}

		if oreNeeded <= oreQuantity {
			min = mid
		} else {
			max = mid - 1
		}
	}

	return min, nil
}
