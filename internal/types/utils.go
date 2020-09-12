package types

import "strings"

func flattenTokens(tokens []Token) string {
	output := ""
	for _, token := range tokens {
		output += token.Text
	}
	return output
}

func prettyDocs(docs string) []string {
	if docs == "" {
		return []string{}
	}

	docLines := strings.Split(strings.TrimSpace(docs), "\n")
	if len(docLines) < 1 {
		return []string{}
	}

	for i, line := range docLines {
		docLines[i] = "// " + line
	}
	return docLines
}

func tokenizeFieldList(fields []Field) []Token {
	tokens := []Token{}
	for i, field := range fields {
		if i > 0 {
			tokens = append(tokens,
				Token{",", kindPunctuation},
				tokenSpace)
		}
		if field.Name != "" {
			tokens = append(tokens,
				Token{field.Name, kindFieldName},
				tokenSpace)
		}
		tokens = append(tokens,
			Token{field.Type, kindFieldType})
	}
	return tokens
}
