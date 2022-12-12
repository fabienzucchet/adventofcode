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
	// Output: 21956
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
	// Output: 3709435214239
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
