package utils

// InArray to get if value is in array.
func InArray(arr []string, v string) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}
