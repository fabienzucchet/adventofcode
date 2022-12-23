package fabienz

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 19 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the blueprints.
	blueprints, err := blueprintsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse blueprints: %w", err)
	}

	sum := 0

	for _, bp := range blueprints {
		sum += bp.id * maxGeodeProduced(bp, 25)
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 19 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the blueprints.
	blueprints, err := blueprintsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse blueprints: %w", err)
	}

	product := 1

	for _, bp := range blueprints[:3] {
		product *= maxGeodeProduced(bp, 33)
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type oreRobotCost struct {
	ore int
}

type clayRobotCost struct {
	ore int
}

type obsidianRobotCost struct {
	ore  int
	clay int
}

type geodeRobotCost struct {
	ore      int
	obsidian int
}

type maxCost struct {
	ore      int
	clay     int
	obsidian int
}

type blueprint struct {
	id                int
	oreRobotCost      oreRobotCost
	clayRobotCost     clayRobotCost
	obsidianRobotCost obsidianRobotCost
	geodeRobotCost    geodeRobotCost
	maxCost           maxCost
}

var bpRegex = regexp.MustCompile(`^Blueprint ([0-9]+): Each ore robot costs ([0-9]+) ore\. Each clay robot costs ([0-9]+) ore\. Each obsidian robot costs ([0-9]+) ore and ([0-9]+) clay\. Each geode robot costs ([0-9]+) ore and ([0-9]+) obsidian\.$`)

// Parse a blueprint from a line.
func blueprintFromLine(line string) (blueprint, error) {
	matches := bpRegex.FindStringSubmatch(line)
	if len(matches) != 8 {
		return blueprint{}, fmt.Errorf("could not parse blueprint from line %q", line)
	}

	bp := blueprint{}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse blueprint id from line %q: %w", line, err)
	}
	bp.id = id

	oreRobotCostInOre, err := strconv.Atoi(matches[2])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse ore robot cost from line %q: %w", line, err)
	}
	bp.oreRobotCost.ore = oreRobotCostInOre

	clayRobotCostInOre, err := strconv.Atoi(matches[3])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse clay robot cost from line %q: %w", line, err)
	}
	bp.clayRobotCost.ore = clayRobotCostInOre

	obsidianRobotCostInOre, err := strconv.Atoi(matches[4])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse obsidian robot cost from line %q: %w", line, err)
	}
	bp.obsidianRobotCost.ore = obsidianRobotCostInOre

	obsidianRobotCostInClay, err := strconv.Atoi(matches[5])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse obsidian robot cost from line %q: %w", line, err)
	}
	bp.obsidianRobotCost.clay = obsidianRobotCostInClay

	geodeRobotCostInOre, err := strconv.Atoi(matches[6])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse geode robot cost from line %q: %w", line, err)
	}
	bp.geodeRobotCost.ore = geodeRobotCostInOre

	geodeRobotCostInObsidian, err := strconv.Atoi(matches[7])
	if err != nil {
		return blueprint{}, fmt.Errorf("could not parse geode robot cost from line %q: %w", line, err)
	}
	bp.geodeRobotCost.obsidian = geodeRobotCostInObsidian

	maxOreCost := math.MinInt
	for _, v := range []int{oreRobotCostInOre, clayRobotCostInOre, obsidianRobotCostInOre, geodeRobotCostInOre} {
		if v > maxOreCost {
			maxOreCost = v
		}
	}
	bp.maxCost.ore = maxOreCost
	bp.maxCost.clay = obsidianRobotCostInClay
	bp.maxCost.obsidian = geodeRobotCostInObsidian

	return bp, nil
}

// Parse a list of blueprints from a list of lines.
func blueprintsFromLines(lines []string) ([]blueprint, error) {
	bps := []blueprint{}
	for _, line := range lines {
		bp, err := blueprintFromLine(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse blueprint from line %q: %w", line, err)
		}
		bps = append(bps, bp)
	}
	return bps, nil
}

// The state of the factory at a given time.
type factoryState struct {
	time          int
	ore           int
	clay          int
	obsidian      int
	geode         int
	oreRobot      int
	clayRobot     int
	obsidianRobot int
	geodeRobot    int
}

// Use a BFS search to find the maximum number of geodes that can be produced with a given amount of time.
func maxGeodeProduced(bp blueprint, time int) int {
	initialState := factoryState{
		time:          0,
		ore:           0,
		clay:          0,
		obsidian:      0,
		geode:         0,
		oreRobot:      1,
		clayRobot:     0,
		obsidianRobot: 0,
		geodeRobot:    0,
	}

	maxGeode := 0

	// Queue to store the states to explore.
	queue := []factoryState{initialState}

	// Iterate as long as there are states to explore.
	for len(queue) > 0 {
		// Pop the first state from the queue.
		state := queue[0]
		queue = queue[1:]

		// If we run out of time, we can't do anything.
		if state.time >= time {
			continue
		}

		// If we have more geodes than the current max, update the max.
		if state.geode > maxGeode {
			maxGeode = state.geode
		}

		// If we can afford a geode robot, buy one.
		if state.ore >= bp.geodeRobotCost.ore && state.obsidian >= bp.geodeRobotCost.obsidian {
			queue = appendIfCouldBeBetter(queue, factoryState{
				time:          state.time + 1,
				ore:           state.ore - bp.geodeRobotCost.ore + state.oreRobot,
				clay:          state.clay + state.clayRobot,
				obsidian:      state.obsidian - bp.geodeRobotCost.obsidian + state.obsidianRobot,
				geode:         state.geode + state.geodeRobot,
				oreRobot:      state.oreRobot,
				clayRobot:     state.clayRobot,
				obsidianRobot: state.obsidianRobot,
				geodeRobot:    state.geodeRobot + 1,
			}, maxGeode)

			// We prune all other branches if we can build a geode robot.
			continue
		}

		// If we can afford an obsidian robot, buy one.
		if state.obsidianRobot < bp.maxCost.obsidian && state.ore >= bp.obsidianRobotCost.ore && state.clay >= bp.obsidianRobotCost.clay {
			queue = appendIfCouldBeBetter(queue, factoryState{
				time:          state.time + 1,
				ore:           state.ore - bp.obsidianRobotCost.ore + state.oreRobot,
				clay:          state.clay - bp.obsidianRobotCost.clay + state.clayRobot,
				obsidian:      state.obsidian + state.obsidianRobot,
				geode:         state.geode + state.geodeRobot,
				oreRobot:      state.oreRobot,
				clayRobot:     state.clayRobot,
				obsidianRobot: state.obsidianRobot + 1,
				geodeRobot:    state.geodeRobot,
			}, maxGeode)
		}

		// If we can afford a clay robot, buy one.
		if state.clayRobot < bp.maxCost.clay && state.ore >= bp.clayRobotCost.ore && state.ore-state.oreRobot < bp.clayRobotCost.ore {
			queue = appendIfCouldBeBetter(queue, factoryState{
				time:          state.time + 1,
				ore:           state.ore - bp.clayRobotCost.ore + state.oreRobot,
				clay:          state.clay + state.clayRobot,
				obsidian:      state.obsidian + state.obsidianRobot,
				geode:         state.geode + state.geodeRobot,
				oreRobot:      state.oreRobot,
				clayRobot:     state.clayRobot + 1,
				obsidianRobot: state.obsidianRobot,
				geodeRobot:    state.geodeRobot,
			}, maxGeode)
		}

		// If we can afford an ore robot, buy one.
		if state.oreRobot < bp.maxCost.ore && state.ore >= bp.oreRobotCost.ore && state.ore-state.oreRobot < bp.oreRobotCost.ore {
			queue = appendIfCouldBeBetter(queue, factoryState{
				time:          state.time + 1,
				ore:           state.ore - bp.oreRobotCost.ore + state.oreRobot,
				clay:          state.clay + state.clayRobot,
				obsidian:      state.obsidian + state.obsidianRobot,
				geode:         state.geode + state.geodeRobot,
				oreRobot:      state.oreRobot + 1,
				clayRobot:     state.clayRobot,
				obsidianRobot: state.obsidianRobot,
				geodeRobot:    state.geodeRobot,
			}, maxGeode)
		}

		// Do not produce any robot.
		queue = appendIfCouldBeBetter(queue, factoryState{
			time:          state.time + 1,
			ore:           state.ore + state.oreRobot,
			clay:          state.clay + state.clayRobot,
			obsidian:      state.obsidian + state.obsidianRobot,
			geode:         state.geode + state.geodeRobot,
			oreRobot:      state.oreRobot,
			clayRobot:     state.clayRobot,
			obsidianRobot: state.obsidianRobot,
			geodeRobot:    state.geodeRobot,
		}, maxGeode)
	}

	return maxGeode
}

func appendIfCouldBeBetter(queue []factoryState, state factoryState, currentMaxGeode int) []factoryState {
	if state.geode+state.geodeRobot >= currentMaxGeode {
		return append(queue, state)
	}

	return queue
}
