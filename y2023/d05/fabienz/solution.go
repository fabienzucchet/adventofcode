package fabienz

import (
	"fmt"
	"io"
	"regexp"
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

	a, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse lines: %w", err)
	}

	transformedSeeds, err := a.transformSeeds()
	if err != nil {
		return fmt.Errorf("could not transform seeds: %w", err)
	}

	lowestSeed := getLowestSeed(transformedSeeds)

	_, err = fmt.Fprintf(answer, "%d", lowestSeed)
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

	_, err = fmt.Fprintf(answer, "%d", len(lines))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// 50 98 2 translates into MappedInterval{98, 2, 50-98}
type MappedInterval struct {
	Start  int
	Length int
	Offset int
}

func (m *MappedInterval) String() string {
	return fmt.Sprintf("[%d,%d] -> [%d, %d]", m.Start, m.Start+m.Length-1, m.Start+m.Offset, m.Start+m.Length-1+m.Offset)
}

func (m *MappedInterval) Copy() MappedInterval {
	return MappedInterval{m.Start, m.Length, m.Offset}
}

func (m *MappedInterval) Contains(value int) bool {
	return value >= m.Start && value < m.Start+m.Length
}

func (m *MappedInterval) Map(value int) (int, error) {
	if !m.Contains(value) {
		return 0, fmt.Errorf("value %d is not in interval %v", value, m)
	}
	return value + m.Offset, nil
}

func parseInterval(s string) (MappedInterval, error) {
	var targetStart, sourceStart, length int
	_, err := fmt.Sscanf(s, "%d %d %d", &targetStart, &sourceStart, &length)
	if err != nil {
		return MappedInterval{}, fmt.Errorf("could not parse interval %q: %w", s, err)
	}
	return MappedInterval{sourceStart, length, targetStart - sourceStart}, nil
}

type MappedIntervals []MappedInterval

func (mm MappedIntervals) String() string {
	var s string
	for _, m := range mm {
		s += m.String() + " "
	}
	return s
}

// Given a list of intervals, sort them by their start value in ascending order
func (mm MappedIntervals) Sort() {
	for i := 0; i < len(mm); i++ {
		for j := i + 1; j < len(mm); j++ {
			if mm[i].Start > mm[j].Start {
				mm[i], mm[j] = mm[j], mm[i]
			}
		}
	}
}

func (mm MappedIntervals) Map(value int) (int, error) {
	for _, m := range mm {
		if m.Contains(value) {
			mappedValue, err := m.Map(value)
			if err != nil {
				return 0, fmt.Errorf("could not map value %d: %w", value, err)
			}
			return mappedValue, nil
		}
	}
	return value, nil
}

func parseSeeds(s string) (seeds MappedIntervals, err error) {

	re := regexp.MustCompile(`\d+`)
	rawSeeds := re.FindAllString(s, -1)

	seeds = make([]MappedInterval, len(rawSeeds))

	for i, rawSeed := range rawSeeds {
		seed, err := strconv.Atoi(rawSeed)
		if err != nil {
			return nil, fmt.Errorf("could not parse seed %q: %w", seed, err)
		}

		seeds[i] = MappedInterval{seed, 1, 0}
	}

	// Sort the seeds by their start value
	seeds.Sort()

	return seeds, nil
}

type Almanac struct {
	Seeds            MappedIntervals
	AlmanacIntervals []MappedIntervals
}

func (a *Almanac) transformSeed(seed int) (int, error) {
	for _, interval := range a.AlmanacIntervals {
		mappedSeed, err := interval.Map(seed)
		if err != nil {
			return 0, fmt.Errorf("could not map seed %d: %w", seed, err)
		}
		seed = mappedSeed
	}
	return seed, nil
}

func (a *Almanac) transformSeeds() ([]int, error) {
	transformedSeeds := make([]int, len(a.Seeds))
	for i, seed := range a.Seeds {
		transformedSeed, err := a.transformSeed(seed.Start)
		if err != nil {
			return nil, fmt.Errorf("could not transform seed %d: %w", seed, err)
		}
		transformedSeeds[i] = transformedSeed
	}
	return transformedSeeds, nil
}

// Merge two MappedIntervals using the following logic:
// - We assume that both MappedIntervals are sorted by their start value in ascending order
// - We follow a traditional merge algorithm but we need to keep the values of Offset distincts
// - Any interval that is not overlapping will be added to the result with an Offset of 0
func mergeIntervals(a, b MappedIntervals) MappedIntervals {
	merged := MappedIntervals{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		// Case 1: a[i] is before b[j] - keep a[i] the same as we merge with an offset of 0
		if a[i].Start+a[i].Length < b[j].Start {
			merged = append(merged, a[i])
			i++
		} else if a[i].Start > b[j].Start {
			merged = append(merged, b[j])
			j++
		} else {
			// We have an overlap
			// We need to keep the Offset distinct
			if a[i].Offset < b[j].Offset {
				merged = append(merged, a[i])
			} else {
				merged = append(merged, b[j])
			}
			i++
			j++
		}
	}

	// Add the remaining intervals
	for i < len(a) {
		merged = append(merged, a[i])
		i++
	}
	for j < len(b) {
		merged = append(merged, b[j])
		j++
	}

	return merged
}

// We use the fact that the lowest seed will always be at the start of an interval
func (a *Almanac) getLowestLocation() int {
	lowestLocation := a.Seeds[0].Start
	for _, interval := range a.AlmanacIntervals {
		if interval[0].Start < lowestLocation {
			lowestLocation = interval[0].Start
		}
	}
	return lowestLocation
}

func getLowestSeed(seeds []int) int {
	lowestSeed := seeds[0]
	for _, seed := range seeds {
		if seed < lowestSeed {
			lowestSeed = seed
		}
	}
	return lowestSeed
}

func parseLines(lines []string) (Almanac, error) {

	a := Almanac{}
	mapIdx := -1

	for i, line := range lines {
		if i == 0 {
			seeds, err := parseSeeds(line)
			if err != nil {
				return Almanac{}, fmt.Errorf("could not parse seeds: %w", err)
			}

			a.Seeds = seeds
			continue
		}

		if strings.Contains(line, "map:") {
			mapIdx += 1
			a.AlmanacIntervals = append(a.AlmanacIntervals, MappedIntervals{})
			continue
		}

		if line != "" {
			interval, err := parseInterval(line)
			if err != nil {
				return Almanac{}, fmt.Errorf("could not parse interval: %w", err)
			}

			a.AlmanacIntervals[mapIdx] = append(a.AlmanacIntervals[mapIdx], interval)
			continue
		}
	}

	// Sort the intervals by their start value
	for i := range a.AlmanacIntervals {
		a.AlmanacIntervals[i].Sort()
	}

	return a, nil
}
