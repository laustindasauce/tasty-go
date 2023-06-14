package utils

// ContainsInt returns whether or not the int exists in the slice
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
