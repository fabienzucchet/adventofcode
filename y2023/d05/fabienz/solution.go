package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 5 of Advent of Code 2023.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	a, err := parseInput(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	locations := []int{}

	for _, seed := range a.seeds {
		locations = append(locations, a.getLocationForSeed(seed))
	}

	_, err = fmt.Fprintf(answer, "%d", helpers.MinInts(locations))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 5 of Advent of Code 2023.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// TODO: Write code to solve Part 2 here.

	// TODO: Write your solution to Part 2 below.
	_, err = fmt.Fprintf(answer, "%d", len(lines))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

var LINES = map[string][2]int{
	"SEED_TO_SOIL":            {3, 50},
	"SOIL_TO_FERTILIZER":      {53, 88},
	"FERTILIZER_TO_WATER":     {91, 120},
	"WATER_TO_LIGHT":          {123, 148},
	"LIGHT_TO_TEMPERATURE":    {151, 194},
	"TEMPERATURE_TO_HUMIDITY": {197, 203},
	"HUMIDITY_TO_LOCATION":    {206, 247},
}

// var LINES = map[string][2]int{
// 	"SEED_TO_SOIL":            {3, 4},
// 	"SOIL_TO_FERTILIZER":      {7, 9},
// 	"FERTILIZER_TO_WATER":     {12, 15},
// 	"WATER_TO_LIGHT":          {18, 19},
// 	"LIGHT_TO_TEMPERATURE":    {22, 24},
// 	"TEMPERATURE_TO_HUMIDITY": {27, 28},
// 	"HUMIDITY_TO_LOCATION":    {31, 32},
// }

type Range struct {
	sourceStart      int
	destinationStart int
	length           int
}

type Almanac struct {
	seeds                 []int
	seedToSoil            []Range
	soilToFertilizer      []Range
	fertilizerToWater     []Range
	waterToLight          []Range
	lightToTemperature    []Range
	temperatureToHumidity []Range
	humidityToLocation    []Range
}

func parseInput(lines []string) (a Almanac, err error) {
	a = Almanac{}

	// Parse the seeds.
	seeds := strings.Split(lines[0], " ")[1:]
	for _, seed := range seeds {
		seedInt, err := strconv.Atoi(seed)
		if err != nil {
			return a, fmt.Errorf("could not parse seed: %w", err)
		}
		a.seeds = append(a.seeds, seedInt)
	}

	// Parse the seed to soil map.
	for i := LINES["SEED_TO_SOIL"][0]; i <= LINES["SEED_TO_SOIL"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse seed to soil map: %w", err)
		}
		a.seedToSoil = append(a.seedToSoil, r)
	}

	// Parse the soil to fertilizer map.
	for i := LINES["SOIL_TO_FERTILIZER"][0]; i <= LINES["SOIL_TO_FERTILIZER"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse soil to fertilizer map: %w", err)
		}
		a.soilToFertilizer = append(a.soilToFertilizer, r)
	}

	// Parse the fertilizer to water map.
	for i := LINES["FERTILIZER_TO_WATER"][0]; i <= LINES["FERTILIZER_TO_WATER"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse fertilizer to water map: %w", err)
		}
		a.fertilizerToWater = append(a.fertilizerToWater, r)
	}

	// Parse the water to light map.
	for i := LINES["WATER_TO_LIGHT"][0]; i <= LINES["WATER_TO_LIGHT"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse water to light map: %w", err)
		}
		a.waterToLight = append(a.waterToLight, r)
	}

	// Parse the light to temperature map.
	for i := LINES["LIGHT_TO_TEMPERATURE"][0]; i <= LINES["LIGHT_TO_TEMPERATURE"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse light to temperature map: %w", err)
		}
		a.lightToTemperature = append(a.lightToTemperature, r)
	}

	// Parse the temperature to humidity map.
	for i := LINES["TEMPERATURE_TO_HUMIDITY"][0]; i <= LINES["TEMPERATURE_TO_HUMIDITY"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse temperature to humidity map: %w", err)
		}
		a.temperatureToHumidity = append(a.temperatureToHumidity, r)
	}

	// Parse the humidity to location map.
	for i := LINES["HUMIDITY_TO_LOCATION"][0]; i <= LINES["HUMIDITY_TO_LOCATION"][1]; i++ {
		r, err := parseRange(lines[i])
		if err != nil {
			return a, fmt.Errorf("could not parse humidity to location map: %w", err)
		}
		a.humidityToLocation = append(a.humidityToLocation, r)
	}

	return a, nil
}

func parseRange(s string) (r Range, err error) {
	fmt.Sscanf(s, "%d %d %d", &r.destinationStart, &r.sourceStart, &r.length)

	return r, nil
}

func (a *Almanac) getLocationForSeed(seed int) (location int) {
	soil := findInMap(a.seedToSoil, seed)
	fertilizer := findInMap(a.soilToFertilizer, soil)
	water := findInMap(a.fertilizerToWater, fertilizer)
	light := findInMap(a.waterToLight, water)
	temperature := findInMap(a.lightToTemperature, light)
	humidity := findInMap(a.temperatureToHumidity, temperature)
	location = findInMap(a.humidityToLocation, humidity)

	return location
}

func findInMap(m []Range, idx int) (newIdx int) {
	newIdx = idx

	for _, r := range m {
		if idx >= r.sourceStart && idx < r.sourceStart+r.length {
			newIdx = r.destinationStart + (idx - r.sourceStart)
			break
		}
	}

	return newIdx
}
