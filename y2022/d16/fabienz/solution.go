package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Start by parsing the input into a map of valves.
	valves, err := valvesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse valves: %w", err)
	}

	// Determine the valves worth investigating (i.e. having a flow > 0).
	interestingValves := getInterestingValves(valves)

	// Compute the distances from any valve to the interesting valves.
	distances := buildDistances(valves, interestingValves)

	maxFlow := 0
	initialPipeSystem := pipeSystem{
		time:         1,
		position:     "AA",
		openedValves: []string{},
		currentFlow:  0,
		totalFlow:    0,
	}
	// Queue to process the pipe systems.
	pipeSystems := []pipeSystem{initialPipeSystem}

	// Process the pipe systems.
	for len(pipeSystems) > 0 {
		s := pipeSystems[0]
		pipeSystems = pipeSystems[1:]
		// If we find a better solution, we save it.
		if f := s.totalFlowIfStays(30); f > maxFlow {
			maxFlow = f
		}
		v := valves[s.position]

		// If we are short on time, we can't explore further.
		if s.time == 30 {
			continue
		}

		// If the valve is not already open and is worth opening, we open it.
		if !isStringInSlice(s.position, s.openedValves) && v.flow > 0 {
			pipeSystems = append(pipeSystems, pipeSystem{
				time:         s.time + 1,
				position:     s.position,
				openedValves: append(s.openedValves, s.position),
				currentFlow:  s.currentFlow + v.flow,
				totalFlow:    s.totalFlow + s.currentFlow,
			})
			continue
		}

		// Else we explore the possible paths.
		for _, t := range interestingValves {
			// If the valve is already open, we can't go there.
			if isStringInSlice(t, s.openedValves) {
				continue
			}
			d := distances[s.position][t]
			// If we don't have enough time to go there, we can't go there.
			if s.time+d > 29 {
				continue
			}
			pipeSystems = append(pipeSystems, pipeSystem{
				time:         s.time + d,
				position:     t,
				openedValves: copySlice(s.openedValves),
				currentFlow:  s.currentFlow,
				totalFlow:    s.totalFlow + s.currentFlow*d,
			})
		}
	}

	_, err = fmt.Fprintf(answer, "%d", maxFlow)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Start by parsing the input into a map of valves.
	valves, err := valvesFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse valves: %w", err)
	}

	// Determine the valves worth investigating (i.e. having a flow > 0).
	interestingValves := getInterestingValves(valves)

	// Compute the distances from any valve to the interesting valves.
	distances := buildDistances(valves, interestingValves)

	// Instead of only keeping the maximum flow, we will store the states that are interesting i.e. the state with the maximum flow for each combination of opened valves.
	cachedStates := make(map[string]pipeSystem)
	cachedMaxFlow := make(map[string]int)

	initialPipeSystem := pipeSystem{
		time:         1,
		position:     "AA",
		openedValves: []string{},
		currentFlow:  0,
		totalFlow:    0,
	}
	// Queue to process the pipe systems.
	pipeSystems := []pipeSystem{initialPipeSystem}

	// Process the pipe systems.
	for len(pipeSystems) > 0 {
		s := pipeSystems[0]
		pipeSystems = pipeSystems[1:]

		v := valves[s.position]

		// If we are short on time, we can't explore further.
		if s.time == 26 {
			continue
		}

		// If the valve is not already open and is worth opening, we open it.
		if !isStringInSlice(s.position, s.openedValves) && v.flow > 0 {
			newS := pipeSystem{
				time:         s.time + 1,
				position:     s.position,
				openedValves: append(s.openedValves, s.position),
				currentFlow:  s.currentFlow + v.flow,
				totalFlow:    s.totalFlow + s.currentFlow,
			}
			pipeSystems = append(pipeSystems, newS)

			// Convert the opened valves to a string to use it as a key in the map.
			openedValvesString := convertSliceToStringKey(newS.openedValves)
			// If we find a better solution, we save it.
			if f := newS.totalFlowIfStays(26); f > cachedMaxFlow[openedValvesString] {
				cachedStates[openedValvesString] = newS
				cachedMaxFlow[openedValvesString] = f
			}
			continue
		}

		// Else we explore the possible paths.
		for _, t := range interestingValves {
			// If the valve is already open, we can't go there.
			if isStringInSlice(t, s.openedValves) {
				continue
			}
			d := distances[s.position][t]
			// If we don't have enough time to go there, we can't go there.
			if s.time+d > 25 {
				continue
			}
			pipeSystems = append(pipeSystems, pipeSystem{
				time:         s.time + d,
				position:     t,
				openedValves: copySlice(s.openedValves),
				currentFlow:  s.currentFlow,
				totalFlow:    s.totalFlow + s.currentFlow*d,
			})
		}
	}

	// Now that we have the best solution for each combination of opened valves, we can dispatch human and elephants to open the valves.
	maxFlow := 0

	for _, hs := range cachedStates {
		for _, es := range cachedStates {
			// Check that the two states are compatible i.e. that no valve is opened twice.
			if !areStatesCompatible(hs, es) {
				continue
			}

			// If the two states are compatible, we can compute the total flow.
			f := hs.totalFlowIfStays(26) + es.totalFlowIfStays(26)
			// If we find a better solution, we save it.
			if f > maxFlow {
				maxFlow = f
			}
		}
	}

	_, err = fmt.Fprintf(answer, "%d", maxFlow)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type valve struct {
	name      string
	flow      int
	neighbors []string
}

// Represent the state of the system.
type pipeSystem struct {
	time         int
	position     string
	openedValves []string
	currentFlow  int
	totalFlow    int
}

// Return the total flow if the system stays like this until the end.
func (p pipeSystem) totalFlowIfStays(allowedTime int) int {
	return p.totalFlow + p.currentFlow*(allowedTime-p.time+1)
}

// Represent the distance between two valves.
type distances map[string]map[string]int

var parseRegex = regexp.MustCompile(`^Valve ([A-Z]+) has flow rate=([0-9]+); tunnels? leads? to valves? (.+)$`)

// Parse the input lines into a map of valves.
func valvesFromLines(lines []string) (map[string]valve, error) {
	valves := make(map[string]valve, len(lines))

	// Create the valves
	for i, line := range lines {
		matches := parseRegex.FindStringSubmatch(line)
		if len(matches) != 4 {
			return nil, fmt.Errorf("could not parse line %d: %s", i, line)
		}

		name := matches[1]
		flow, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("could not parse flow rate for valve %s: %w", name, err)
		}

		neighbors := matches[3]
		neighborsParsed := strings.Split(neighbors, ", ")

		valves[name] = valve{name: name, flow: flow, neighbors: neighborsParsed}
	}

	return valves, nil
}

// Return a slice of the valves having a flow rate > 0.
func getInterestingValves(valves map[string]valve) []string {
	var interestingValves []string

	for _, v := range valves {
		if v.flow > 0 {
			interestingValves = append(interestingValves, v.name)
		}
	}

	return interestingValves
}

// Build the distances from any valve to the interesting valves.
func buildDistances(valves map[string]valve, interestingValves []string) distances {
	distances := make(distances, len(valves))

	for v := range valves {
		distances[v] = make(map[string]int, len(valves))
		for _, t := range interestingValves {
			distances[v][t] = distBetweenValves(v, t, valves)
		}
	}

	return distances
}

type partialDist struct {
	valve string
	dist  int
}

// Compute the distance between two valves.
func distBetweenValves(from, to string, valves map[string]valve) int {
	if from == to {
		return 0
	}

	visited := make(map[string]bool, len(valves))
	visited[from] = true

	q := []partialDist{{from, 0}}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		for _, n := range valves[p.valve].neighbors {
			if n == to {
				return p.dist + 1
			}

			if !visited[n] {
				visited[n] = true
				q = append(q, partialDist{n, p.dist + 1})
			}
		}
	}

	return -1
}

// Check if a string is contained in a slice of strings.
func isStringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// Copy a slice of strings.
func copySlice(slice []string) []string {
	newSlice := make([]string, len(slice))
	copy(newSlice, slice)
	return newSlice
}

// Convert a slice of strings to a string.
func convertSliceToStringKey(slice []string) string {
	// Sort the strings to have a deterministic key.
	sort.Strings(slice)
	return strings.Join(slice, ",")
}

// Check if two states are compatible i.e. that no valve is opened twice.
func areStatesCompatible(hs, es pipeSystem) bool {
	for _, v := range hs.openedValves {
		if isStringInSlice(v, es.openedValves) {
			return false
		}
	}
	return true
}
