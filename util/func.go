package util

import "strconv"

// Pow10 ...
func Pow10(a uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < a; i++ {
		result *= 10
	}
	return result
}

// ParseUint64 ...
func ParseUint64(a string) uint64 {
	result, err := strconv.ParseUint(a, 10, 64)
	if err != nil {
		return 0
	}
	return result
}
