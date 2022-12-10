package fabienz

import (
	"fmt"
	"io"

	"github.com/fabienzucchet/adventofcode/helpers"
)

// PartOne solves the first problem of day 20 of Advent of Code 2021.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	image, scale := parseLines(lines, 4)

	for i := 0; i < 2; i++ {
		image = enhance(image, scale)
	}

	_, err = fmt.Fprintf(answer, "%d", countLitPixels(image))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2021.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	image, scale := parseLines(lines, 52)

	for i := 0; i < 50; i++ {
		image = enhance(image, scale)
	}

	_, err = fmt.Fprintf(answer, "%d", countLitPixels(image))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Image [][]string

type Scale [512]string

// INPUT PARSING

// Padding is the number of elements around the initial image to add. They are initialized to 0
func parseLines(lines []string, padding int) (image Image, scale Scale) {
	// The first line corresponds to the scale
	for idx, char := range lines[0] {
		scale[idx] = string(char)
	}

	size := len(lines[2])

	for j := -padding; j < size+padding; j++ {
		var row []string
		for i := -padding; i < size+padding; i++ {
			if j < 0 || j >= size {
				row = append(row, ".")
			} else {
				if i < 0 || i >= size {
					row = append(row, ".")
				} else {
					row = append(row, string(lines[2+j][i]))
				}
			}
		}
		image = append(image, row)
	}

	return image, scale
}

func enhance(image Image, scale Scale) (newImage Image) {

	// Enhance each pixel of the image
	for row, line := range image {
		var newRow []string
		for col, pixel := range line {
			// If the pixel is at the border of the image, it will alway be "." if scale[0] == . and will blink if scale[0] == #
			if row == 0 || col == 0 || row == len(image)-1 || col == len(image)-1 {
				if scale[0] == "#" {
					switch pixel {
					case ".":
						newRow = append(newRow, "#")
					case "#":
						newRow = append(newRow, ".")
					}
				} else {
					newRow = append(newRow, ".")
				}
			} else {
				// Otherwise we compute the new value of the pixel by looking at its neighbors
				newRow = append(newRow, scale[getScaleIdx(row, col, image)])
			}
		}
		newImage = append(newImage, newRow)
	}

	return newImage
}

// Get the scale index for a given pixel
func getScaleIdx(row, col int, image Image) (idx int) {
	var bin string

	for j := row + 1; j >= row-1; j-- {
		for i := col + 1; i >= col-1; i-- {
			bin += image[j][i]
		}
	}

	return toDec(bin)
}

// Find the binary value of a binary string
func toDec(bin string) (decimal int) {
	if len(bin) == 0 {
		return 0
	}

	var current int
	switch bin[0] {
	case '.':
		current = 0
	case '#':
		current = 1
	}

	return 2*toDec(bin[1:]) + current
}

// Count the lit pixels
func countLitPixels(image Image) (count int) {

	for _, line := range image {
		for _, pixel := range line {
			if pixel == "#" {
				count++
			}
		}
	}

	return count
}
