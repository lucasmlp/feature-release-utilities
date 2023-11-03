package utils

// indexOf returns the index of a string in a slice, or a large number if not found.
func IndexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return 1<<31 - 1 // Use max int to denote not found; this will sort it to the end.
}
