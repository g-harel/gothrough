package gis

import (
	"unicode"
)

// CamelSplit breaks the camel cased input string into its tokens.
func CamelSplit(str string) []string {
	tokens := []string{}
	var lastRune rune

	for i, r := range str {
		// First rune will always start a new token.
		if len(tokens) == 0 {
			tokens = append(tokens, string(r))
			lastRune = r
			continue
		}

		// Start a new token when current rune is uppercase and next one is
		// lowercase. This condition will trigger to end a token correctly when
		// it is comprised of uppercase letters (ex. HTTPTest => HTTP, Test).
		if unicode.IsUpper(r) && len(str) > i+1 && !unicode.IsUpper(rune(str[i+1])) {
			tokens = append(tokens, string(r))
			lastRune = r
			continue
		}

		// When current rune is lowercase or previous one is uppercase, append
		// to last token. This logic captures the case where multiple uppercase
		// letters are part of the same token. (ex ServeHTTP => Serve, HTTP).
		if !unicode.IsUpper(r) || unicode.IsUpper(lastRune) {
			tokens[len(tokens)-1] = tokens[len(tokens)-1] + string(r)
			lastRune = r
			continue
		}

		// Create new token when no special conditions are met.
		tokens = append(tokens, string(r))
		lastRune = r
	}

	return tokens
}
