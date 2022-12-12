package fabienz

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 11 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	monkeys, err := parseMonkeys(lines)
	if err != nil {
		return fmt.Errorf("could not parse monkeys: %w", err)
	}

	// Compute 20 rounds.
	for i := 0; i < 20; i++ {
		for idx := range monkeys {
			monkeys.computeRound(idx, false, 1)
		}
	}

	itemsInspected := []int{}
	for _, monkey := range monkeys {
		itemsInspected = append(itemsInspected, monkey.itemsInspected)
	}

	sort.Slice(itemsInspected, func(i, j int) bool {
		return itemsInspected[i] > itemsInspected[j]
	})

	_, err = fmt.Fprintf(answer, "%d", itemsInspected[0]*itemsInspected[1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 11 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	monkeys, err := parseMonkeys(lines)
	if err != nil {
		return fmt.Errorf("could not parse monkeys: %w", err)
	}

	// Find the LCM of tge monkey's test values.
	lcm := 1
	for _, monkey := range monkeys {
		lcm = helpers.LCM([]int{lcm, monkey.TestValue})
	}

	// Compute 10000 rounds.
	for i := 0; i < 10000; i++ {
		for idx := range monkeys {
			monkeys.computeRound(idx, true, lcm)
		}
	}

	itemsInspected := []int{}
	for _, monkey := range monkeys {
		itemsInspected = append(itemsInspected, monkey.itemsInspected)
	}

	sort.Slice(itemsInspected, func(i, j int) bool {
		return itemsInspected[i] > itemsInspected[j]
	})

	_, err = fmt.Fprintf(answer, "%d", itemsInspected[0]*itemsInspected[1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Operation struct {
	Op     string
	Arg    int
	HasArg bool
}

type Monkey struct {
	Id             int
	Items          []int
	Operation      Operation
	TestValue      int
	IfTrue         int
	IfFalse        int
	itemsInspected int
}

type Monkeys []*Monkey

// Parse the lines of the input file and return a slice of monkeys.
func parseMonkeys(lines []string) (Monkeys, error) {
	monkeys := []*Monkey{}

	currentMonkey := &Monkey{}

	for idx, line := range lines {
		switch idx % 7 {
		case 0:
			_, err := fmt.Sscanf(line, "Monkey %d:", &currentMonkey.Id)
			if err != nil {
				return nil, fmt.Errorf("could not parse line %s: %w", line, err)
			}
		case 1:
			itemsStr := strings.Split(line[18:], ", ")
			items := []int{}
			for _, itemStr := range itemsStr {
				item, err := strconv.Atoi(itemStr)
				if err != nil {
					return nil, fmt.Errorf("could not parse line %s: %w", line, err)
				}
				items = append(items, item)
			}
			currentMonkey.Items = items
		case 2:
			Operation := Operation{}
			var argumentStr string
			_, err := fmt.Sscanf(line, "  Operation: new = old %s %s", &Operation.Op, &argumentStr)
			if err != nil {
				return nil, fmt.Errorf("could not parse line %s: %w", line, err)
			}
			// Parse the argument if argument is not old.
			if argumentStr != "old" {
				Operation.HasArg = true
				argument, err := strconv.Atoi(argumentStr)
				if err != nil {
					return nil, fmt.Errorf("could not parse line %s: %w", line, err)
				}
				Operation.Arg = argument
			}
			currentMonkey.Operation = Operation
		case 3:
			_, err := fmt.Sscanf(line, "  Test: divisible by %d", &currentMonkey.TestValue)
			if err != nil {
				return nil, fmt.Errorf("could not parse line %s: %w", line, err)
			}
		case 4:
			_, err := fmt.Sscanf(line, "    If true: throw to monkey %d", &currentMonkey.IfTrue)
			if err != nil {
				return nil, fmt.Errorf("could not parse line %s: %w", line, err)
			}
		case 5:
			_, err := fmt.Sscanf(line, "    If false: throw to monkey %d", &currentMonkey.IfFalse)
			if err != nil {
				return nil, fmt.Errorf("could not parse line %s: %w", line, err)
			}

			monkeys = append(monkeys, currentMonkey)
			currentMonkey = &Monkey{}
		}

	}

	return monkeys, nil
}

// Compute a round for a monkey.
func (m *Monkeys) computeRound(idx int, worried bool, testLCM int) {
	monkey := (*m)[idx]
	// If the monkey has no items, do nothing.
	if len(monkey.Items) == 0 {
		return
	}
	for _, item := range monkey.Items {
		monkey.itemsInspected++
		// Compute the new value of the item.
		newValue := item
		if monkey.Operation.HasArg {
			switch monkey.Operation.Op {
			case "+":
				newValue += monkey.Operation.Arg
			case "*":
				newValue *= monkey.Operation.Arg
			}
		} else {
			switch monkey.Operation.Op {
			case "+":
				newValue = 2 * item
			case "*":
				newValue = item * item
			}
		}
		if !worried {
			// The inspection didn't damage the item, so the worry is divided by 3.
			newValue /= 3
		} else {
			// Else we need to keep low newvalues so we only keep the modulo since we only perform modulo tests.
			newValue %= testLCM
		}
		// If the new value is divisible by the test value, throw the item to the monkey
		// specified in IfTrue. Otherwise, throw the item to the monkey specified in IfFalse.
		if newValue%monkey.TestValue == 0 {
			(*m)[monkey.IfTrue].Items = append((*m)[monkey.IfTrue].Items, newValue)
		} else {
			(*m)[monkey.IfFalse].Items = append((*m)[monkey.IfFalse].Items, newValue)
		}

	}
	// Remove the items from the current monkey.
	monkey.Items = []int{}
}
