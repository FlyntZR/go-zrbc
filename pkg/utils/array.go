package utils

// ArrayDiff returns the values in array1 that are not present in any of the other arrays.
// This is similar to PHP's array_diff function.
func ArrayDiff[T comparable](array1 []T, arrays ...[]T) []T {
	if len(array1) == 0 {
		return []T{}
	}

	// Create a map to store values from other arrays
	exists := make(map[T]bool)

	// Add all values from other arrays to the map
	for _, arr := range arrays {
		for _, val := range arr {
			exists[val] = true
		}
	}

	// Create result slice for values in array1 that don't exist in other arrays
	result := make([]T, 0)

	// Check each value in array1
	for _, val := range array1 {
		if !exists[val] {
			result = append(result, val)
		}
	}

	return result
}

// ArrayDiffString is a convenience function for string slices
func ArrayDiffString(array1 []string, arrays ...[]string) []string {
	return ArrayDiff(array1, arrays...)
}

// ArrayDiffInt is a convenience function for int slices
func ArrayDiffInt(array1 []int, arrays ...[]int) []int {
	return ArrayDiff(array1, arrays...)
}

// ArrayDiffInt64 is a convenience function for int64 slices
func ArrayDiffInt64(array1 []int64, arrays ...[]int64) []int64 {
	return ArrayDiff(array1, arrays...)
}
