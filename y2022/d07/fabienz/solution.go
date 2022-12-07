package fabienz

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"gitlab.com/padok-team/adventofcode/helpers"
)

// PartOne solves the first problem of day 7 of Advent of Code 2022.
func PartOne(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input to create the filesystem.
	fs, err := ParseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Compute the size of the filesystem.
	fs.ComputeSize()

	// Find all directories with a size of at most 100000.
	dirs := fs.FindDirectoriesWithSizeAtMost(100000)

	_, err = fmt.Fprintf(answer, "%d", TotalSize(dirs))
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 7 of Advent of Code 2022.
func PartTwo(input io.Reader, answer io.Writer) error {
	// Read the input. Feel free to change it depending on the input.
	lines, err := helpers.LinesFromReader(input)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	// Parse the input to create the filesystem.
	fs, err := ParseLines(lines)
	if err != nil {
		return fmt.Errorf("could not parse input: %w", err)
	}

	// Compute the size of the filesystem.
	fs.ComputeSize()

	// Find all directories with a size of at least 30000000 - (70000000 - fs.size) aka space required - space already free.
	dirs := fs.FindDirectoriesWithSizeAtLeast(30000000 - (70000000 - fs.size))

	// Find the best candidate to remove.
	bestCandidate := FindSmallestDirectory(dirs)
	if (bestCandidate == nil) || (bestCandidate.size == 0) {
		return fmt.Errorf("could not find a candidate to remove")
	}

	_, err = fmt.Fprintf(answer, "%d", bestCandidate.size)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// We will represent the system with a tree.
type File struct {
	name        string
	children    []*File
	parent      *File
	isDirectory bool
	size        int
}

// Various regexps to parse the input.
var (
	// cd regexp.
	cdRegexp = regexp.MustCompile(`^\$ cd ([a-zA-Z0-9\/.]+)$`)
	// ls regexp.
	lsRegexp = regexp.MustCompile(`^\$ ls$`)
	// directory regexp.
	dirRegexp = regexp.MustCompile(`^dir ([a-zA-Z0-9]+)$`)
	// file regexp.
	fileRegexp = regexp.MustCompile(`^([0-9]+) ([a-zA-Z0-9.]+)$`)
)

// Parse the input and return a filesystem.
func ParseLines(lines []string) (fs *File, err error) {
	// Create the filesystem.
	fs = &File{
		name:        "/",
		children:    make([]*File, 0),
		parent:      nil,
		isDirectory: true,
		size:        0,
	}

	// Keep track of the current directory.
	currentDir := fs

	for _, line := range lines {
		// If the line is a cd, change the current directory.
		if cdRegexp.MatchString(line) {
			// Get the directory name.
			dirName := cdRegexp.FindStringSubmatch(line)[1]

			// If the directory name is "..", go to the parent directory.
			if dirName == ".." {
				currentDir = currentDir.parent
				continue
			}

			// Otherwise, find the directory in the current directory.
			for _, child := range currentDir.children {
				if child.name == dirName {
					currentDir = child
					break
				}
			}

			continue
		}

		// If the line is a ls, then continue.
		if lsRegexp.MatchString(line) {
			continue
		}

		// If the line is a directory, then create it.
		if dirRegexp.MatchString(line) {
			// Get the directory name.
			dirName := dirRegexp.FindStringSubmatch(line)[1]

			// Create the directory.
			dir := &File{
				name:        dirName,
				children:    make([]*File, 0),
				parent:      currentDir,
				isDirectory: true,
				size:        0,
			}

			// Add the directory to the current directory.
			currentDir.children = append(currentDir.children, dir)

			continue
		}

		// If the line is a file, then create it.
		if fileRegexp.MatchString(line) {
			// Get the file size and name.
			fileSize, err := strconv.Atoi(fileRegexp.FindStringSubmatch(line)[1])
			if err != nil {
				return fs, fmt.Errorf("could not parse file size: %w", err)
			}
			fileName := fileRegexp.FindStringSubmatch(line)[2]

			// Create the file.
			file := &File{
				name:        fileName,
				children:    make([]*File, 0),
				parent:      currentDir,
				isDirectory: false,
				size:        fileSize,
			}

			// Add the file to the current directory.
			currentDir.children = append(currentDir.children, file)

			continue
		}

		// If we reach this point, then the line is not valid.
		return fs, fmt.Errorf("invalid line: %s", line)
	}

	return fs, nil
}

// Recursively compute the size of the directories.
func (f *File) ComputeSize() {
	// If the file is a directory, then compute the size of its children.
	if f.isDirectory {
		for _, child := range f.children {
			child.ComputeSize()
			f.size += child.size
		}
	}
}

// Recursively find all directories with a total size of at most maxSize.
func (f *File) FindDirectoriesWithSizeAtMost(maxSize int) (dirs []*File) {
	// If the file is a directory, then check if its size is at most maxSize.
	if f.isDirectory {
		if f.size <= maxSize {
			dirs = append(dirs, f)
		}

		for _, child := range f.children {
			dirs = append(dirs, child.FindDirectoriesWithSizeAtMost(maxSize)...)
		}
	}

	return dirs
}

// Recursively find all directories with a total size of at least minSize.
func (f *File) FindDirectoriesWithSizeAtLeast(minSize int) (dirs []*File) {
	// If the file is a directory, then check if its size is at least minSize.
	if f.isDirectory {
		if f.size >= minSize {
			dirs = append(dirs, f)
		}

		for _, child := range f.children {
			dirs = append(dirs, child.FindDirectoriesWithSizeAtLeast(minSize)...)
		}
	}

	return dirs
}

// Compute the total size of a slice of files.
func TotalSize(files []*File) (size int) {
	for _, file := range files {
		size += file.size
	}

	return size
}

// Find the directory with the smallest size in a slice of files.
func FindSmallestDirectory(files []*File) (dir *File) {
	// If the slice is empty, return nil.
	if len(files) == 0 {
		return nil
	}

	// Otherwise, find the smallest directory.
	dir = files[0]
	for _, file := range files {
		if file.size < dir.size {
			dir = file
		}
	}

	return dir
}

// Print a file.
func (f *File) Print() {
	helpers.Println("Name:", f.name, "Size:", f.size, "Is directory:", f.isDirectory)
}
