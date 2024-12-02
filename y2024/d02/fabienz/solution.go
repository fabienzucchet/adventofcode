package fabienz

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 2 of Advent of Code 2024.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	reports, err := reportsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse reports: %w", err)
	}

	count := countValidReports(reports)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 2 of Advent of Code 2024.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	reports, err := reportsFromLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse reports: %w", err)
	}

	count := countValidReportsWhenSkippingIdx(reports)

	_, err = fmt.Fprintf(answer, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Report []int

func reportFromLine(line string) (Report, error) {
	stringLevels := strings.Split(line, " ")

	levels := make(Report, len(stringLevels))

	for i, stringLevel := range stringLevels {
		level, err := strconv.Atoi(stringLevel)
		if err != nil {
			return nil, fmt.Errorf("could not parse level %s: %w", stringLevel, err)
		}

		levels[i] = level
	}

	return levels, nil
}

func reportsFromLines(lines []string) ([]Report, error) {
	reports := make([]Report, len(lines))

	for i, line := range lines {
		report, err := reportFromLine(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse report from line %s: %w", line, err)
		}

		reports[i] = report
	}

	return reports, nil
}

func (r Report) IsAscending() bool {
	for i := 1; i < len(r); i++ {
		levelIncrease := r[i] - r[i-1]
		if levelIncrease <= 0 || levelIncrease > 3 {
			return false
		}
	}

	return true
}

func (r Report) IsDescending() bool {
	for i := 1; i < len(r); i++ {
		levelIncrease := r[i-1] - r[i]
		if levelIncrease <= 0 || levelIncrease > 3 {
			return false
		}
	}

	return true
}

func (r Report) IsValid() bool {
	return r.IsAscending() || r.IsDescending()
}

func countValidReports(reports []Report) int {
	count := 0

	for _, report := range reports {
		if report.IsValid() {
			count++
		}
	}

	return count
}

func (r Report) IsAscendingWhenSkippingIdx(idx int) bool {
	for i := 1; i < len(r); i++ {
		if i == idx {
			continue
		}

		levelIncrease := r[i] - r[i-1]
		if levelIncrease <= 0 || levelIncrease > 3 {
			return false
		}
	}

	return true
}

func (r Report) IsDescendingWhenSkippingIdx(idx int) bool {
	for i := 1; i < len(r); i++ {
		if i == idx {
			continue
		}

		levelIncrease := r[i-1] - r[i]
		if levelIncrease <= 0 || levelIncrease > 3 {
			return false
		}
	}

	return true
}

func (r Report) IsValidWhenSkippingIdx() bool {
	for i := range r {
		subReport := Report{}
		subReport = append(subReport, r[:i]...)
		subReport = append(subReport, r[i+1:]...)
		if subReport.IsValid() {
			return true
		}
	}

	return false
}

func countValidReportsWhenSkippingIdx(reports []Report) int {
	count := 0

	for _, report := range reports {
		if report.IsValidWhenSkippingIdx() {
			count++
		}
	}

	return count
}
