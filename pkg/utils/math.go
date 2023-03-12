package utils

import "golang.org/x/exp/constraints"

// Max is a generic function that returns the larger of two numbers
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
