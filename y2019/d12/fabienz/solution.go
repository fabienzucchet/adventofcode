package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 12 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the moons.
	moons := parseMoons(lines)

	// Iterate the system 1000 times.
	for i := 0; i < 1000; i++ {
		iterateMoons(moons)
	}

	// Compute the total energy.
	var totalEnergy int
	for _, m := range moons {
		totalEnergy += m.totalEnergy()
	}

	_, err = fmt.Fprintf(answer, "%d", totalEnergy)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 12 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the moons.
	moons := parseMoons(lines)

	// Find the period of the system.
	period := findPeriod(moons)

	_, err = fmt.Fprintf(answer, "%d", period)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// Represent a moon with its position and velocity.
type moon struct {
	x, y, z    int
	vx, vy, vz int
}

// Return the position of a moon.
func (m moon) position() [3]int {
	return [3]int{m.x, m.y, m.z}
}

// Return the velocity of a moon.
func (m moon) velocity() [3]int {
	return [3]int{m.vx, m.vy, m.vz}
}

// Parse a moon from a string.
func parseMoon(s string) moon {
	var m moon
	fmt.Sscanf(s, "<x=%d, y=%d, z=%d>", &m.x, &m.y, &m.z)
	return m
}

// Parse all moons from a list of strings.
func parseMoons(lines []string) []moon {
	moons := make([]moon, len(lines))
	for i, line := range lines {
		moons[i] = parseMoon(line)
	}
	return moons
}

// Update the velocity of a moon based on the position of another moon.
func (m *moon) updateVelocity(other moon) {
	if m.x < other.x {
		m.vx++
	} else if m.x > other.x {
		m.vx--
	}

	if m.y < other.y {
		m.vy++
	} else if m.y > other.y {
		m.vy--
	}

	if m.z < other.z {
		m.vz++
	} else if m.z > other.z {
		m.vz--
	}
}

// Update the position of a moon based on its velocity.
func (m *moon) updatePosition() {
	m.x += m.vx
	m.y += m.vy
	m.z += m.vz
}

// Compute the potential energy of a moon.
func (m moon) potentialEnergy() int {
	return helpers.AbsInt(m.x) + helpers.AbsInt(m.y) + helpers.AbsInt(m.z)
}

// Compute the kinetic energy of a moon.
func (m moon) kineticEnergy() int {
	return helpers.AbsInt(m.vx) + helpers.AbsInt(m.vy) + helpers.AbsInt(m.vz)
}

// Compute the total energy of a moon.
func (m moon) totalEnergy() int {
	return m.potentialEnergy() * m.kineticEnergy()
}

// Iterate the moons for one step.
func iterateMoons(moons []moon) {
	for i := range moons {
		for j := range moons {
			if i == j {
				continue
			}
			moons[i].updateVelocity(moons[j])
		}
	}

	for i := range moons {
		moons[i].updatePosition()
	}
}

// We will consider only one dimension at a time for the second part of the problem.
type body1D struct {
	p, v int
}

type system1D struct {
	bodies []body1D
}

// Iterate the 1D system for one step.
func (s *system1D) iterate() {
	for i := range s.bodies {
		for j := range s.bodies {
			if i == j {
				continue
			}
			if s.bodies[i].p < s.bodies[j].p {
				s.bodies[i].v++
			} else if s.bodies[i].p > s.bodies[j].p {
				s.bodies[i].v--
			}
		}
	}

	for i := range s.bodies {
		s.bodies[i].p += s.bodies[i].v
	}
}

// Find the period of a 3D system by using 3 1D systems.
func findPeriod(moons []moon) (period int) {
	// Init 3 1D systems with the moons.
	systems := [3]system1D{}
	for i := 0; i < 3; i++ {
		bodies := []body1D{}
		for _, m := range moons {
			body := body1D{p: m.position()[i], v: m.velocity()[i]}
			bodies = append(bodies, body)
		}
		systems[i] = system1D{bodies: bodies}
	}

	// Store the state of each 1D system in a map.
	states := [3]map[string]bool{}
	for i := 0; i < 3; i++ {
		states[i] = make(map[string]bool)
	}

	for counter := 0; ; counter++ {
		// Iterate each 1D system.
		for i := 0; i < 3; i++ {
			systems[i].iterate()
		}

		// Stop if we have seen the state of each 1D system.
		if states[0][systems[0].key()] && states[1][systems[1].key()] && states[2][systems[2].key()] {
			break
		}

		// Store the state of each 1D system.
		for i := 0; i < 3; i++ {
			states[i][systems[i].key()] = true
		}
	}

	return helpers.LCM([]int{len(states[0]), len(states[1]), len(states[2])})
}

// Create a key string of a 1D system.
func (s system1D) key() string {
	return fmt.Sprint(s.bodies)
}
