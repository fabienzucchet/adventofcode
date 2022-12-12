package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	field_lines, _, tickets := parseLines(lines)

	var constraints []Constraint

	for _, field := range field_lines {
		parsed_field, err := parseField(field)
		if err != nil {
			return fmt.Errorf("error parsing field %s : %w", field, err)
		}

		for i := 1; i < len(parsed_field); i++ {
			parsed_constraint := strings.Split(parsed_field[i], "-")

			min, err := strconv.Atoi(parsed_constraint[0])
			if err != nil {
				return fmt.Errorf("error parsing min constraint %s : %w", parsed_constraint[0], err)
			}

			max, err := strconv.Atoi(parsed_constraint[1])
			if err != nil {
				return fmt.Errorf("error parsing max constraint %s : %w", parsed_constraint[1], err)
			}

			constraints = append(constraints, Constraint{min, max})
		}
	}

	var values []int

	for _, ticket := range tickets {
		for _, v := range strings.Split(ticket, ",") {
			val, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("error parsing ticket value %s : %w", v, err)
			}
			values = append(values, val)
		}
	}

	sum := 0
	for _, val := range values {
		if !isValueValid(val, constraints) {
			sum += val
		}
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	field_lines, your_ticket, tickets := parseLines(lines)
	var constraints []Constraint

	for _, field := range field_lines {
		parsed_field, err := parseField(field)
		if err != nil {
			return fmt.Errorf("error parsing field %s : %w", field, err)
		}

		for i := 1; i < len(parsed_field); i++ {
			parsed_constraint := strings.Split(parsed_field[i], "-")

			min, err := strconv.Atoi(parsed_constraint[0])
			if err != nil {
				return fmt.Errorf("error parsing min constraint %s : %w", parsed_constraint[0], err)
			}

			max, err := strconv.Atoi(parsed_constraint[1])
			if err != nil {
				return fmt.Errorf("error parsing max constraint %s : %w", parsed_constraint[1], err)
			}

			constraints = append(constraints, Constraint{min, max})
		}
	}

	var values []int

	for _, ticket := range tickets {
		for _, v := range strings.Split(ticket, ",") {
			val, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("error parsing ticket value %s : %w", v, err)
			}
			values = append(values, val)
		}
	}

	// Selecting the valid tickets
	var valid_tickets [][]int
	for _, ticket := range tickets {
		valid, err := isTicketValid(ticket, constraints)
		if err != nil {
			return fmt.Errorf("error check ticket %s : %w", ticket, err)
		}

		if valid {
			var ticket_values []int
			for _, v := range strings.Split(ticket, ",") {
				val, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("error parsing value %s : %w", v, err)
				}
				ticket_values = append(ticket_values, val)
			}
			valid_tickets = append(valid_tickets, ticket_values)
		}
	}

	// Create list of Fields
	var fields []Field

	for idx, field := range field_lines {
		parsed_field, err := parseField(field)
		if err != nil {
			return fmt.Errorf("error parsing field %s : %w", field, err)
		}

		var constraints []Constraint

		for i := 1; i < len(parsed_field); i++ {
			parsed_constraint := strings.Split(parsed_field[i], "-")

			min, err := strconv.Atoi(parsed_constraint[0])
			if err != nil {
				return fmt.Errorf("error parsing min constraint %s : %w", parsed_constraint[0], err)
			}

			max, err := strconv.Atoi(parsed_constraint[1])
			if err != nil {
				return fmt.Errorf("error parsing max constraint %s : %w", parsed_constraint[1], err)
			}

			constraints = append(constraints, Constraint{min, max})
		}

		fields = append(fields, Field{constraints, parsed_field[0], idx})
	}

	// Find possible fields for each columns
	possibleFields := make([][]Field, len(fields))

	for i := 0; i < len(fields); i++ {
		for _, field := range fields {
			if isFieldMatchingColumn(field, i, valid_tickets) {
				possibleFields[i] = append(possibleFields[i], field)
			}
		}
	}

	// Try to make a mapping
	mapping := make([]Field, len(fields))
	affected := make([]bool, len(fields))

	var recursive_helper func(int) bool
	recursive_helper = func(idx int) bool {

		// If end of the field list => we can affect fields
		if idx == len(fields) {
			return true
		}

		// For each possible fields try to make an affectation
		for _, f := range possibleFields[idx] {

			// If already affected, try next one
			if affected[f.position] {
				continue
			}

			// Try to affect and see what happens for the resmaining indexes
			mapping[idx] = f
			affected[f.position] = true

			// If the mapping is working for the remaining indexes
			if recursive_helper(idx + 1) {
				return true
			}

			// Else try another one
			mapping[idx] = Field{}
			affected[f.position] = false
		}

		return false
	}

	if !recursive_helper(0) {
		return fmt.Errorf("could not find mapping")
	}

	//Make the product of all fields starting with departure
	product := 1

	for idx, value := range strings.Split(your_ticket, ",") {
		if len(mapping[idx].name) >= 9 && mapping[idx].name[:9] == "departure" {

			val, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing value %s : %w", value, err)
			}

			product *= val
		}
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Constraint struct {
	min int
	max int
}

type Field struct {
	constraints []Constraint
	name        string
	position    int
}

func parseLines(lines []string) ([]string, string, []string) {
	var fields []string
	var your_ticket string
	var tickets []string

	section := "fields"

	for _, line := range lines {
		if line == "" {
			continue
		} else if line == "your ticket:" {
			section = "yourticket"
		} else if line == "nearby tickets:" {
			section = "nearbytickets"
		} else {
			if section == "fields" {
				fields = append(fields, line)
			} else if section == "yourticket" {
				your_ticket = line
			} else {
				tickets = append(tickets, line)
			}
		}
	}

	return fields, your_ticket, tickets
}

func parseField(field string) ([]string, error) {
	re := regexp.MustCompile(`(.*): (.*) or (.*)`)
	match := re.FindStringSubmatch(field)

	if len(match) == 0 {
		return nil, fmt.Errorf("error parsing field %s", field)
	}

	return match[1:], nil
}

func isValueValid(val int, constraints []Constraint) bool {

	checkConstraint := false
	for _, c := range constraints {
		if val <= c.max && val >= c.min {
			checkConstraint = true
		}
	}

	return checkConstraint
}

func isTicketValid(ticket string, constraints []Constraint) (bool, error) {
	valid := true

	for _, v := range strings.Split(ticket, ",") {
		val, err := strconv.Atoi(v)
		if err != nil {
			return false, fmt.Errorf("error parsing value %s : %w", v, err)
		}

		if !isValueValid(val, constraints) {
			valid = false
		}
	}

	return valid, nil
}

func isValueMatchingOneConstraints(value int, constraints []Constraint) bool {
	res := false

	for _, constraint := range constraints {
		if value >= constraint.min && value <= constraint.max {
			res = true
		}
	}

	return res
}

// Check if the field matches the column
func isFieldMatchingColumn(field Field, idx int, tickets [][]int) bool {
	res := true

	for _, ticket := range tickets {
		if !isValueMatchingOneConstraints(ticket[idx], field.constraints) {
			res = false
		}
	}

	return res
}
