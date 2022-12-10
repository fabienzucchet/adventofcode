package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 16 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	transmission := hexaToBin(lines[0])

	packets, _ := decodeTransmission(transmission, 1)

	_, err = fmt.Fprintf(answer, "%d", sumVersionNumbers(packets))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 16 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	transmission := hexaToBin(lines[0])

	packets, _ := decodeTransmission(transmission, 1)

	_, err = fmt.Fprintf(answer, "%d", evaluate(packets[0]))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Packet struct {
	version    int
	pType      int
	value      int
	subPackets []Packet
}

// INPUT PARSING

var hexaToBinMap = map[rune][]int{
	'0': {0, 0, 0, 0},
	'1': {0, 0, 0, 1},
	'2': {0, 0, 1, 0},
	'3': {0, 0, 1, 1},
	'4': {0, 1, 0, 0},
	'5': {0, 1, 0, 1},
	'6': {0, 1, 1, 0},
	'7': {0, 1, 1, 1},
	'8': {1, 0, 0, 0},
	'9': {1, 0, 0, 1},
	'A': {1, 0, 1, 0},
	'B': {1, 0, 1, 1},
	'C': {1, 1, 0, 0},
	'D': {1, 1, 0, 1},
	'E': {1, 1, 1, 0},
	'F': {1, 1, 1, 1},
}

// Convert Hexadecimal into binary string
func hexaToBin(hexa string) (bin []int) {

	for _, char := range hexa {
		bin = append(bin, hexaToBinMap[char]...)
	}

	return bin
}

func decodeTransmission(transmission []int, maxPackets int) (packets []Packet, remainingTransmission []int) {

	for (maxPackets == -1 || len(packets) < maxPackets) && len(transmission) > 0 {

		p := Packet{}

		version := binToDec(transmission[:3])
		transmission = transmission[3:]

		p.version = version

		pType := binToDec(transmission[:3])
		transmission = transmission[3:]

		p.pType = pType

		switch pType {
		// Case of a value packet
		case 4:
			var value []int
			for {
				group := transmission[:5]
				transmission = transmission[5:]
				value = append(value, group[1:]...)

				if group[0] == 0 {
					break
				}
			}

			p.value = binToDec(value)

		// Otherwise it's an operator packet
		default:
			lengthTypeId := transmission[0]
			transmission = transmission[1:]

			var subPackets []Packet

			switch lengthTypeId {
			// We know the length of the sub packets
			case 0:
				length := binToDec(transmission[:15])
				transmission = transmission[15:]

				subPackets, _ = decodeTransmission(transmission[:length], -1)
				transmission = transmission[length:]

				p.subPackets = append(p.subPackets, subPackets...)

			// We know the number of sub packets
			case 1:
				count := binToDec(transmission[:11])
				transmission = transmission[11:]

				subPackets, transmission = decodeTransmission(transmission, count)

				p.subPackets = append(p.subPackets, subPackets...)
			}
		}

		packets = append(packets, p)
	}

	return packets, transmission
}

// Convert from Binary to decimal
func binToDec(bin []int) (dec int) {

	for _, bit := range bin {
		dec = dec*2 + bit
	}

	return dec
}

// Sum all version numbers of a sub packet
func sumVersionNumbers(packets []Packet) (sum int) {

	for _, p := range packets {

		sum += p.version + sumVersionNumbers(p.subPackets)
	}

	return sum
}

// Evaluate the expression of packet
func evaluate(packet Packet) (res int) {

	switch packet.pType {
	// Value packet
	case 4:
		return packet.value

	// Sum packet
	case 0:
		for _, p := range packet.subPackets {
			res += evaluate(p)
		}

	// Product packet
	case 1:
		res = 1
		for _, p := range packet.subPackets {
			res *= evaluate(p)
		}

	// Minimum packet
	case 2:
		mini := 1 << 62
		for _, p := range packet.subPackets {
			val := evaluate(p)
			if val < mini {
				mini = val
			}
		}
		res = mini

	// Maximum packet
	case 3:
		maxi := -1
		for _, p := range packet.subPackets {
			val := evaluate(p)
			if val > maxi {
				maxi = val
			}
		}
		res = maxi

	// Greater than packet
	case 5:
		if evaluate(packet.subPackets[0]) > evaluate(packet.subPackets[1]) {
			return 1
		} else {
			return 0
		}

	// Less than packet
	case 6:
		if evaluate(packet.subPackets[0]) < evaluate(packet.subPackets[1]) {
			return 1
		} else {
			return 0
		}

	// Equal packet
	case 7:
		if evaluate(packet.subPackets[0]) == evaluate(packet.subPackets[1]) {
			return 1
		} else {
			return 0
		}
	}

	return res
}
