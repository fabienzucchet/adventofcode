package fabienz

import (
	"log"
	"os"
	"testing"

	"github.com/fabienzucchet/adventofcode/helpers"
)

func ExamplePartOne() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer file.Close()

	if err := PartOne(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}
	// Output: 2562
}

func ExamplePartTwo() {
	file, err := os.Open("testdata/input.txt")
	if err != nil {
		log.Fatalf("could open input file: %v", err)
	}
	defer file.Close()

	if err := PartTwo(file, os.Stdout); err != nil {
		log.Fatalf("could not solve: %v", err)
	}
	// Output: [1 1 1 1 0 1 1 1 1 0 1 0 0 0 0 1 1 1 0 0 1 0 0 0 1 0 0 0 1 0 1 0 0 0 0 1 0 0 0 0 1 0 0 1 0 1 0 0 0 1 0 0 1 0 0 1 1 1 0 0 1 0 0 0 0 1 1 1 0 0 0 1 0 1 0 0 1 0 0 0 1 0 0 0 0 1 0 0 0 0 1 0 0 1 0 0 0 1 0 0 1 0 0 0 0 1 0 0 0 0 1 0 0 0 0 1 0 0 1 0 0 0 1 0 0 1 1 1 1 0 1 0 0 0 0 1 1 1 1 0 1 1 1 0 0 0 0 1 0 0]
}

func Benchmark(b *testing.B) {
	testCases := map[string]struct {
		solution  helpers.Solution
		inputFile string
	}{
		"PartOne": {
			solution:  helpers.SolutionFunc(PartOne),
			inputFile: "testdata/input.txt",
		},

		"PartTwo": {
			solution:  helpers.SolutionFunc(PartTwo),
			inputFile: "testdata/input.txt",
		},
	}

	for name, test := range testCases {
		b.Run(name, func(b *testing.B) {
			helpers.BenchmarkSolution(b, test.solution, test.inputFile)
		})
	}
}
