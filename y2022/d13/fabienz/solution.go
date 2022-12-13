package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 13 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	packetPairs, err := parsePacketPairs(lines)
	if err != nil {
		return fmt.Errorf("could not parse packet pairs: %w", err)
	}

	sum := 0

	for idx, pair := range packetPairs {
		if pair.compare() < 0 {
			sum += idx + 1
		}
	}

	_, err = fmt.Fprintf(answer, "%d", sum)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 13 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Add the two decoder packets to lines
	lines = append(lines, "[[2]]")
	lines = append(lines, "[[6]]")

	var packets []Packet

	// Parse all packets
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		p, err := packetFromString(line)
		if err != nil {
			return fmt.Errorf("could not parse packet %s: %w", line, err)
		}

		packets = append(packets, p)

	}

	// Sort the packets
	sortedPackets := sortPackets(packets)

	product := 1

	// Find the two indexes
	for idx, packet := range sortedPackets {
		if packet.toString() == "[[2,],]" || packet.toString() == "[[6,],]" {
			product *= (idx + 1)
		}
	}

	_, err = fmt.Fprintf(answer, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

type Packet struct {
	isInt    bool
	value    int
	children []*Packet
	parent   *Packet
}

type PacketPair struct {
	left  Packet
	right Packet
}

func (p Packet) toString() string {
	if p.isInt {
		return strconv.Itoa(p.value)
	}

	s := "["
	for _, child := range p.children {
		s += child.toString()
		s += ","
	}
	s += "]"

	return s
}

// Create the packet pairs from the input.
func parsePacketPairs(lines []string) ([]PacketPair, error) {
	nbPacketPairs := (len(lines)-1)/3 + 1

	pairs := make([]PacketPair, nbPacketPairs)

	for i := 0; i < nbPacketPairs; i++ {
		left, err := packetFromString(lines[3*i])
		if err != nil {
			return nil, fmt.Errorf("could not parse left packet: %w", err)
		}
		right, err := packetFromString(lines[3*i+1])
		if err != nil {
			return nil, fmt.Errorf("could not parse right packet: %w", err)
		}

		pair := PacketPair{
			left:  left,
			right: right,
		}

		pairs[i] = pair
	}

	return pairs, nil
}

// parse a packet from a string
func packetFromString(s string) (Packet, error) {
	// Cursor used to iterate over the string
	idx := 0

	p := Packet{}

	var parseToken func(p *Packet) error
	var parseList func(p *Packet) error

	// Parse the first token
	parseToken = func(p *Packet) error {

		// Read until the end
		for {
			switch {
			case s[idx] == '[':
				// Recursive parsing
				idx++
				p.isInt = false
				return parseList(p)
			case isDigit(s[idx]):
				// Parse the number
				var n int
				for isDigit(s[idx]) {
					n = n*10 + int(s[idx]-'0')
					idx++
				}
				p.isInt = true
				p.value = n
				return nil
			}
		}
	}

	// Recusively parse a list
	parseList = func(p *Packet) error {
		for {
			switch {
			case s[idx] == ']':
				// End of the list
				if p.isInt {
					return fmt.Errorf("unexpected end of list")
				}
				idx++
				return nil
			case s[idx] == ',':
				// End of a token
				if p.isInt {
					return fmt.Errorf("unexpected end of token")
				}
				idx++
			default:
				var child Packet
				child.parent = p
				err := parseToken(&child)
				if err != nil {
					return fmt.Errorf("could not parse token: %w", err)
				}
				p.children = append(p.children, &child)
			}
		}
	}

	if err := parseToken(&p); err != nil {
		return Packet{}, fmt.Errorf("could not parse token: %w", err)
	}

	return p, nil
}

// Check if a character is a digit
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Define the comparison for a packet. Returns < 0 if left < right, > 0 if left > right and 0 if left == right.
func (p Packet) compare(other Packet) int {
	// Both are integers
	if p.isInt && other.isInt {
		return p.value - other.value
	}

	// Convert p to a list
	if p.isInt {
		return Packet{children: []*Packet{&p}}.compare(other)
	}

	// Convert other to a list
	if other.isInt {
		return p.compare(Packet{children: []*Packet{&other}})
	}

	// Otherwise compare the lists
	for idx := range p.children {
		// If the other list is shorter, return false
		if idx >= len(other.children) {
			return 1
		}
		if res := p.children[idx].compare(*other.children[idx]); res != 0 {
			return res
		}
	}

	// If the other list is longer, return true
	if len(p.children) < len(other.children) {
		return -1
	}

	return 0
}

// Compare two packets in a packet pair and return true if they are in the correct order.
func (p PacketPair) compare() int {
	return p.left.compare(p.right)
}

// Pile
type Pile []string

// Push a string to the pile.
func (p *Pile) Push(s string) {
	*p = append(*p, s)
}

// Pop a string from the pile.
func (p *Pile) Pop() string {
	if len(*p) == 0 {
		return ""
	}

	s := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]

	return s
}

// Len returns the length of the pile.
func (p *Pile) Len() int {
	return len(*p)
}

// Sort a list of packets using bubble sort
func sortPackets(packets []Packet) []Packet {
	for i := 0; i < len(packets); i++ {
		for j := 0; j < len(packets); j++ {
			if packets[i].compare(packets[j]) < 0 {
				packets[i], packets[j] = packets[j], packets[i]
			}
		}
	}

	return packets
}
