package helpers

// Find the maximum of a slice of integers.
func MaxInts(ints []int) int {
	max := ints[0]

	for _, i := range ints {
		if i > max {
			max = i
		}
	}

	return max
}
