package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

const layerWidth = 25
const layerHeight = 6

// PartOne solves the first problem of day 8 of Advent of Code 2019.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the picture.
	picture, err := ParsePicture(lines[0])
	if err != nil {
		return fmt.Errorf("could not parse picture: %w", err)
	}

	_, err = fmt.Fprintf(answer, "%d", picture.Checksum())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 8 of Advent of Code 2019.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the picture.
	picture, err := ParsePicture(lines[0])
	if err != nil {
		return fmt.Errorf("could not parse picture: %w", err)
	}

	// Decode the picture.
	decodedPicture := picture.Decode()

	// Print the decoded picture.
	// decodedPicture.String()

	// The answer can be read by printing the image with the Print function.
	_, err = fmt.Fprintf(answer, "%v", decodedPicture)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// A layer will be represented by a slice of ints.
type Layer [layerWidth * layerHeight]int

// A picture will be represented by a slice of layers.
type Picture []Layer

// ParsePicture parses a picture from a string.
func ParsePicture(s string) (Picture, error) {
	// Compute the number of layers.
	nbLayers := len(s) / (layerWidth * layerHeight)

	// Create the picture.
	picture := make(Picture, nbLayers)

	// Parse the layers.
	for i := 0; i < nbLayers; i++ {
		// Parse the layer.
		layer, err := ParseLayer(s[i*layerWidth*layerHeight : (i+1)*layerWidth*layerHeight])
		if err != nil {
			return nil, fmt.Errorf("could not parse layer: %w", err)
		}

		// Store the layer in the picture.
		picture[i] = layer
	}

	return picture, nil
}

// ParseLayer parses a layer from a string.
func ParseLayer(s string) (Layer, error) {
	// Create the layer.
	layer := Layer{}

	// Parse the layer.
	for i, c := range s {
		// Convert the rune to an int.
		n, err := helpers.IntFromRune(c)
		if err != nil {
			return Layer{}, fmt.Errorf("could not convert rune to int: %w", err)
		}

		// Store the int in the layer.
		layer[i] = n
	}

	return layer, nil
}

// Compute the occurence of each digit in a layer.
func (l Layer) CountDigits() map[int]int {
	counts := make(map[int]int)

	for _, d := range l {
		counts[d]++
	}

	return counts
}

// Compute the checksum of a picture. The checksum is the number of 1 digits multiplied by the number of 2 digits in the layer with the fewest 0 digits.
func (p Picture) Checksum() int {
	// Find the layer with the fewest 0 digits.
	var layer Layer
	min := layerWidth * layerHeight
	for _, l := range p {
		counts := l.CountDigits()
		if counts[0] < min {
			layer = l
			min = counts[0]
		}
	}

	// Compute the checksum.
	counts := layer.CountDigits()
	return counts[1] * counts[2]
}

// Decode a picture.
func (p Picture) Decode() Layer {
	// Create the decoded layer.
	decoded := Layer{}

	// Decode the layer.
	for i := 0; i < layerWidth*layerHeight; i++ {
		// Find the first non-transparent pixel.
		for _, l := range p {
			if l[i] != 2 {
				decoded[i] = l[i]
				break
			}
		}
	}

	return decoded
}

// Print a layer
func (l Layer) String() {
	for i := 0; i < layerHeight; i++ {
		for j := 0; j < layerWidth; j++ {
			if l[i*layerWidth+j] == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("█")
			}
		}
		fmt.Println()
	}
}
