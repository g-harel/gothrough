package gis

import (
	"unicode"
)

func CamelSplit(str string) []string {
	tokens := []string{}
	var lastRune rune

	for _, r := range str {
		if len(tokens) == 0 {
			tokens = append(tokens, string(r))
			continue
		}

		// Append to previous token.
		if !unicode.IsUpper(r) || unicode.IsUpper(lastRune) {
			tokens[len(tokens)-1] = tokens[len(tokens)-1] + string(r)
			continue
		}

		tokens = append(tokens, string(r))
	}

	return tokens
}
