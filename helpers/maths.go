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

func MinInts(ints []int) int {
	min := ints[0]

	for _, i := range ints {
		if i < min {
			min = i
		}
	}

	return min
}

// Find the GCD of two integers.
func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}

	if b < 0 {
		b = -b
	}

	if a < b {
		a, b = b, a
	}

	for b != 0 {
		a, b = b, a%b
	}

	return a
}

// Find the least common multiple of a slice of integers.
func LCM(ints []int) int {
	lcm := ints[0]

	for _, i := range ints {
		lcm = lcm * i / GCD(lcm, i)
	}

	return lcm
}
