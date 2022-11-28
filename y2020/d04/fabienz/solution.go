package fabienz

import (
	"fmt"
	"io"
	"regexp"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 4 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	passports := parseLines(lines)

	var nbCompletePassports int

	for _, passport := range passports {
		if isComplete(passport) {
			nbCompletePassports++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", nbCompletePassports)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 4 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	passports := parseLines(lines)

	var nbValidPassports int

	for _, passport := range passports {
		if isComplete(passport) && isValid(passport) {
			nbValidPassports++
		}
	}

	_, err = fmt.Fprintf(answer, "%d", nbValidPassports)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func isByrValid(passport string) bool {
	re := regexp.MustCompile(`byr:(\S*)`)
	byr := re.FindStringSubmatch(passport)[1]

	return len(byr) == 4 && byr >= "1920" && byr <= "2002"
}

func isIyrValid(passport string) bool {
	re := regexp.MustCompile(`iyr:(\S*)`)
	iyr := re.FindStringSubmatch(passport)[1]

	return len(iyr) == 4 && iyr >= "2010" && iyr <= "2020"
}

func isEyrValid(passport string) bool {
	re := regexp.MustCompile(`eyr:(\S*)`)
	eyr := re.FindStringSubmatch(passport)[1]

	return len(eyr) == 4 && eyr >= "2020" && eyr <= "2030"
}

func isHgtValid(passport string) bool {
	re := regexp.MustCompile(`hgt:(\S*)`)
	hgt := re.FindStringSubmatch(passport)[1]

	if hgt[len(hgt)-2:] == "cm" {
		return hgt[:len(hgt)-2] >= "150" && hgt[:len(hgt)-2] <= "193"

	} else if hgt[len(hgt)-2:] == "in" {
		return hgt[:len(hgt)-2] >= "59" && hgt[:len(hgt)-2] <= "76"
	}

	return false
}

func isHclValid(passport string) bool {
	re := regexp.MustCompile(`hcl:#[0-9a-f]{6}`)
	return re.MatchString(passport)
}

func isEclValid(passport string) bool {
	re := regexp.MustCompile(`ecl:(\S*)`)
	ecl := re.FindStringSubmatch(passport)[1]

	return ecl == "amb" || ecl == "blu" || ecl == "brn" || ecl == "gry" || ecl == "grn" || ecl == "hzl" || ecl == "oth"
}

func isPidValid(passport string) bool {
	re := regexp.MustCompile(`pid:[0-9]{9}(?:\D|$)`)
	return re.MatchString(passport)
}

func isComplete(passport string) bool {
	// Here we assume that passport is represented by one line only

	// Count the number of attributes
	re := regexp.MustCompile(`(\S*:\S*)`)
	nbAttributes := len(re.FindAllString(passport, -1))

	// Check if cid is present
	reHasCID := regexp.MustCompile(`(cid:\S*)`)
	hasCID := reHasCID.MatchString(passport)

	// Passport is valid iif nbAttributes == 8 || (nbAttributes == 7 && cid != present)
	return nbAttributes == 8 || (nbAttributes == 7 && !hasCID)
}

func isValid(passport string) bool {
	return isByrValid(passport) && isIyrValid(passport) && isEyrValid(passport) && isHgtValid(passport) && isHclValid(passport) && isEclValid(passport) && isPidValid(passport)
}

func parseLines(lines []string) []string {
	parsedLines := make([]string, 1)

	for _, line := range lines {
		if line == "" {
			parsedLines = append(parsedLines, "")
		} else if parsedLines[len(parsedLines)-1] == "" {
			parsedLines[len(parsedLines)-1] = parsedLines[len(parsedLines)-1] + line
		} else {
			parsedLines[len(parsedLines)-1] = parsedLines[len(parsedLines)-1] + " " + line
		}
	}

	return parsedLines
}
