package index

// Substrings produces the list of all contiguous substrings of the given size.
func Substrings(str string, size int) []string {
	if size < 1 || len(str) < 1 {
		return []string{}
	}
	if size >= len(str) {
		return []string{str}
	}

	substrings := []string{}
	for i := 0; i+size <= len(str); i++ {
		substrings = append(substrings, str[i:i+size])
	}

	return substrings
}
