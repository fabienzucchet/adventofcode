package fabienz

import (
	"fmt"
	"io"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

const IMAGESIZE = 12
const TILESIZE = 10

// PartOne solves the first problem of day 20 of Advent of Code 2020.
func PartOne(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	tiles, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	image := reassembleTiles(tiles)

	ids := image.getIds()

	_, err = fmt.Fprintf(answer, "%d", ids[0][0]*ids[0][len(ids)-1]*ids[len(ids)-1][0]*ids[len(ids)-1][len(ids)-1])
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2020.
func PartTwo(input io.Reader, answer io.Writer) error {
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	tiles, err := parseLines(lines)
	if err != nil {
		return fmt.Errorf("error parsing input : %w", err)
	}

	image := reassembleTiles(tiles)

	assembledImage := image.assemble()

	var monstersCount int

	for i := 0; i < 4; i++ {
		monstersCount = assembledImage.findSeaMonsters()
		if monstersCount > 0 {
			break
		}
		assembledImage.flip()
		monstersCount = assembledImage.findSeaMonsters()
		if monstersCount > 0 {
			break
		}
		assembledImage.flip()
		assembledImage.rotate()
	}

	_, err = fmt.Fprintf(answer, "%d", assembledImage.findRoughness())
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// TYPES

type Tile struct {
	id  int
	img [TILESIZE][TILESIZE]string
}

type Transformation struct {
	tile   Tile
	flip   bool
	rotate int
}

type Image [IMAGESIZE][IMAGESIZE]Tile

type ImageIds [IMAGESIZE][IMAGESIZE]int

type AssembledImage [IMAGESIZE * (TILESIZE - 2)][IMAGESIZE * (TILESIZE - 2)]string

// INPUT PARSING
func parseLines(lines []string) (tiles [IMAGESIZE * IMAGESIZE]Tile, err error) {

	for i := 0; i < len(lines); i += TILESIZE + 2 {
		t := Tile{}

		// Fetch the tile ID
		tileId, err := strconv.Atoi(lines[i][5:9])
		if err != nil {
			return tiles, fmt.Errorf("error parsing tile ID %s : %w", lines[i][5:9], err)
		}
		t.id = tileId

		// Fetch the tile image
		for row := 0; row < TILESIZE; row++ {
			for col, char := range lines[i+1+row] {
				t.img[row][col] = string(char)
			}
		}

		tiles[i/(TILESIZE+2)] = t
	}

	return tiles, nil
}

// Display a tile
func (t *Tile) display() {
	helpers.Println("Tile", t.id)
	for _, row := range t.img {
		rowString := ""
		for _, char := range row {
			rowString += char
		}
		helpers.Println(rowString)
	}
}

// Flip a tile
func (t *Tile) flip() {

	n := len(t.img)

	for row := range t.img {
		for col := 0; col < n/2; col++ {
			t.img[row][col], t.img[row][n-1-col] = t.img[row][n-1-col], t.img[row][col]
		}
	}
}

// Flip an image
func (ai *AssembledImage) flip() {

	n := len(ai)

	for row := range ai {
		for col := 0; col < n/2; col++ {
			ai[row][col], ai[row][n-1-col] = ai[row][n-1-col], ai[row][col]
		}
	}
}

// Creates a counter-clockwise 90deg rotation on the tile
func (t *Tile) rotate() {

	// First flip the tile
	t.flip()

	// Then transpose it
	for row := range t.img {
		for col := 0; col < row; col++ {
			t.img[row][col], t.img[col][row] = t.img[col][row], t.img[row][col]
		}
	}
}

// Creates a counter-clockwise 90deg rotation on the image
func (ai *AssembledImage) rotate() {

	// First flip the tile
	ai.flip()

	// Then transpose it
	for row := range ai {
		for col := 0; col < row; col++ {
			ai[row][col], ai[col][row] = ai[col][row], ai[row][col]
		}
	}
}

// Copy a tile
func copy(t Tile) (newTile Tile) {
	// Copy the ID
	newTile.id = t.id

	for row := range t.img {
		for col := range t.img[row] {
			newTile.img[row][col] = t.img[row][col]
		}
	}

	return newTile
}

// Converts a transformation struct into a tile
func transform(trans Transformation) (transformed Tile) {

	transformed = copy(trans.tile)

	// Rotate
	for i := 0; i < trans.rotate; i++ {
		transformed.rotate()
	}

	// Flip if needed
	if trans.flip {
		transformed.flip()
	}

	return transformed
}

// Check if an Image has no empty tiles
func (i *Image) isFull() bool {

	for row := range i {
		for col := range i[row] {
			// If a tile of the image has id 0
			if i[row][col].id == 0 {
				return false
			}
		}
	}

	return true
}

// Reassemble the tiles to create the image
func reassembleTiles(tiles [IMAGESIZE * IMAGESIZE]Tile) (image Image) {

	// Map to keep track of the affected tiles. Key = tile ID, value = if the tile is affected
	affected := make(map[int]bool)

	// Compute the possibile transformation for a given cell
	getPossibilites := func() (possibilities []Transformation) {
		// We check each tile to see if it's a possiblity
		for _, tile := range tiles {

			// If the tile is not already affected, its a possibility
			if !affected[tile.id] {

				// Each tile can be present with 8 transformations (combinations of rotation and flip)
				possibilities = append(possibilities, Transformation{tile: tile, flip: false, rotate: 0})
				possibilities = append(possibilities, Transformation{tile: tile, flip: false, rotate: 1})
				possibilities = append(possibilities, Transformation{tile: tile, flip: false, rotate: 2})
				possibilities = append(possibilities, Transformation{tile: tile, flip: false, rotate: 3})
				possibilities = append(possibilities, Transformation{tile: tile, flip: true, rotate: 0})
				possibilities = append(possibilities, Transformation{tile: tile, flip: true, rotate: 1})
				possibilities = append(possibilities, Transformation{tile: tile, flip: true, rotate: 2})
				possibilities = append(possibilities, Transformation{tile: tile, flip: true, rotate: 3})
			}
		}

		return possibilities
	}

	// Check if a tile at a given idx is valid
	isValid := func(idx int) bool {

		// If the tile is not on the left, the left border has to match
		if idx%IMAGESIZE > 0 {
			if image[idx/IMAGESIZE][idx%IMAGESIZE-1].id != 0 && image[idx/IMAGESIZE][idx%IMAGESIZE-1].getBorder("right") != image[idx/IMAGESIZE][idx%IMAGESIZE].getBorder("left") {
				return false
			}
		}

		// If the tile is not on the right, the right border has to match
		if idx%IMAGESIZE < IMAGESIZE-1 {
			if image[idx/IMAGESIZE][idx%IMAGESIZE+1].id != 0 && image[idx/IMAGESIZE][idx%IMAGESIZE+1].getBorder("left") != image[idx/IMAGESIZE][idx%IMAGESIZE].getBorder("right") {
				return false
			}
		}

		// If the tile is not on the top, the top border has to match
		if idx/IMAGESIZE > 0 {
			if image[idx/IMAGESIZE-1][idx%IMAGESIZE].id != 0 && image[idx/IMAGESIZE-1][idx%IMAGESIZE].getBorder("bottom") != image[idx/IMAGESIZE][idx%IMAGESIZE].getBorder("top") {
				return false
			}
		}

		// If the tile is not on the bottom, the bottom border has to match
		if idx/IMAGESIZE < IMAGESIZE-1 {
			if image[idx/IMAGESIZE+1][idx%IMAGESIZE].id != 0 && image[idx/IMAGESIZE+1][idx%IMAGESIZE].getBorder("top") != image[idx/IMAGESIZE][idx%IMAGESIZE].getBorder("bottom") {
				return false
			}
		}

		return true
	}

	var rec func(idx int) bool

	// Idx is the index of a tile i.e. the position of the tile if we count the tiles line by line
	rec = func(idx int) bool {
		// If all tiles on the image are affected, we found a working configuration
		if idx >= IMAGESIZE*IMAGESIZE {
			return true
		}

		// For each possibility, we make a recursive call
		for _, transformation := range getPossibilites() {
			// Let's affect the tile
			image[idx/IMAGESIZE][idx%IMAGESIZE] = transform(transformation)
			affected[transformation.tile.id] = true

			// See if we can complete the image recursively
			if isValid(idx) && rec(idx+1) {
				// It it's validn we won
				return true
			}

			// Otherwise we must remove the current tile and try other possibilities
			image[idx/IMAGESIZE][idx%IMAGESIZE] = Tile{}
			affected[transformation.tile.id] = false
		}

		return false
	}

	rec(0)

	return image
}

// Turn the border of a tile into a string for comparing
func (t *Tile) getBorder(flag string) (border string) {

	n := len(t.img)

	switch flag {
	case "top":
		for i := 0; i < n; i++ {
			border += t.img[0][i]
		}

	case "bottom":
		for i := 0; i < n; i++ {
			border += t.img[n-1][i]
		}

	case "left":
		for i := 0; i < n; i++ {
			border += t.img[i][0]
		}

	case "right":
		for i := 0; i < n; i++ {
			border += t.img[i][n-1]
		}
	}

	return border
}

// Display an image
func (i *Image) display() {
	for imageRow := range i {
		for tileRow := range i[imageRow][0].img {
			var row string
			for imageCol := range i {
				for tileCol := range i[imageRow][imageCol].img {
					row += i[imageRow][imageCol].img[tileRow][tileCol]
				}
				row += "|"
			}
			helpers.Println(row)
		}
		var sep string
		for j := 0; j < IMAGESIZE*(TILESIZE+1); j++ {
			sep += "-"
		}
		helpers.Println(sep)
	}
}

// Display the tileIds of an image
func (i *Image) getIds() (ids ImageIds) {
	for row := range i {
		for col := range i[row] {
			ids[row][col] = i[row][col].id
		}
	}

	return ids
}

// Returns the assembled image without the borders
func (i *Image) assemble() (assembled AssembledImage) {

	for row := range i {
		for col := range i[row] {
			for j := 0; j < TILESIZE-2; j++ {
				for k := 0; k < TILESIZE-2; k++ {
					assembled[row*(TILESIZE-2)+j][col*(TILESIZE-2)+k] = i[row][col].img[1+j][1+k]
				}
			}
		}
	}

	return assembled
}

// This is the sea monster
var seaMonster = [3][20]string{
	{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "#", " "},
	{"#", " ", " ", " ", " ", "#", "#", " ", " ", " ", " ", "#", "#", " ", " ", " ", " ", "#", "#", "#"},
	{" ", "#", " ", " ", "#", " ", " ", "#", " ", " ", "#", " ", " ", "#", " ", " ", "#", " ", " ", " "},
}

// Find sea monsters
func (ai *AssembledImage) findSeaMonsters() (monstersCount int) {

	// Try all possibile positions for the sea monster
	for row := 0; row < len(ai)-1-len(seaMonster); row++ {
		for col := 0; col < len(ai[row])-1-len(seaMonster[0]); col++ {

			// Check if the delimited zone has a sea monster
			if ai.isSeaMonster(row, col) {
				ai.markSeaMonster(row, col)
				monstersCount++
			}
		}
	}

	return monstersCount
}

// Check if a given zone is a sea monster
func (ai *AssembledImage) isSeaMonster(row, col int) bool {

	for j := range seaMonster {
		for i := range seaMonster[j] {
			if seaMonster[j][i] == "#" && ai[row+j][col+i] != "#" {
				return false
			}
		}
	}

	return true
}

// Marks a sea monster
func (ai *AssembledImage) markSeaMonster(row, col int) {
	for j := range seaMonster {
		for i := range seaMonster[j] {
			if seaMonster[j][i] == "#" && ai[row+j][col+i] == "#" {
				ai[row+j][col+i] = "O"
			}
		}
	}
}

// Compute the water roughness
func (ai *AssembledImage) findRoughness() (roughness int) {

	for row := range ai {
		for col := range ai[row] {
			if ai[row][col] == "#" {
				roughness++
			}
		}
	}

	return roughness
}
