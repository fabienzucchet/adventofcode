package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 14 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	mask := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	mem := make(map[int]int)

	for _, line := range lines {
		// If the line is a mask
		if line[1] == 'a' {
			mask = parseMask(line)
		} else {
			addr, val, err := parseAllocation(line)
			if err != nil {
				return fmt.Errorf("error parsing line %s : %w", line, err)
			}

			mem[addr], err = toDecimal(applyMask(toBinary(val), mask))
			if err != nil {
				return fmt.Errorf("error parsing mask : %w", err)
			}
		}
	}

	sum := 0

	for _, val := range mem {
		sum += val
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 14 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	mask := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	mem := make(map[int]int)

	for _, line := range lines {
		// If the line is a mask
		if line[1] == 'a' {
			mask = parseMask(line)
		} else {
			addr, val, err := parseAllocation(line)
			if err != nil {
				return fmt.Errorf("error parsing allocation %s : %w", line, err)
			}

			mem, err = writeInMemory(addr, val, mask, mem)
			if err != nil {
				return fmt.Errorf("error writing in memory : %w", err)
			}
		}
	}

	sum := 0

	for _, val := range mem {
		sum += val
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

func parseMask(line string) string {
	return line[7:]
}

func parseAllocation(line string) (int, int, error) {
	re := regexp.MustCompile(`mem\[([0-9]*)] = ([0-9]*)`)
	g := re.FindStringSubmatch(line)

	if len(g) != 3 {
		return 0, 0, fmt.Errorf("error parsing line %s", line)
	}

	addr, err := strconv.Atoi(g[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing address %s : %w", g[1], err)
	}

	val, err := strconv.Atoi(g[2])
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing value %s : %w", g[2], err)
	}

	return addr, val, nil
}

func toBinary(i int) string {
	s := strconv.FormatInt(int64(i), 2)

	for len(s) < 36 {
		s = "0" + s
	}

	return s
}

func toDecimal(s string) (int, error) {
	i, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing binary %s : %w", s, err)
	}

	return int(i), nil
}

func applyMask(target string, mask string) string {
	var res []byte

	for i := 0; i < len(target); i++ {
		if mask[i] == 'X' {
			res = append(res, target[i])
		} else {
			res = append(res, mask[i])
		}
	}

	return string(res)
}

func applyMask2(target string, mask string) string {
	var res []byte

	for i := 0; i < len(target); i++ {
		if mask[i] == '0' {
			res = append(res, target[i])
		} else {
			res = append(res, mask[i])
		}
	}

	return (string(res))
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func generateAddresses(generated []string, remains string) []string {

	if len(remains) == 0 {
		return generated
	}

	var res []string

	if remains[0] != 'X' {
		for i := 0; i < len(generated); i++ {
			res = append(res, generated[i]+string(remains[0]))
		}
		return generateAddresses(res, remains[1:])
	}

	for i := 0; i < len(generated); i++ {
		res = append(res, generated[i]+"0")
		res = append(res, generated[i]+"1")
	}
	return generateAddresses(res, remains[1:])

}

func writeInMemory(addr int, val int, mask string, mem map[int]int) (map[int]int, error) {

	addr_binary := toBinary(addr)
	addr_masked := applyMask2(addr_binary, mask)

	var target_addresses []int

	for _, a := range generateAddresses([]string{""}, addr_masked) {
		target, err := toDecimal(a)
		if err != nil {
			return nil, fmt.Errorf("error parsing target address %s : %w", a, err)
		}
		target_addresses = append(target_addresses, target)
	}

	for _, a := range target_addresses {
		mem[a] = val
	}

	return mem, nil
}
