package pretty

import (
	"strings"

	"github.com/g-harel/gothrough/internal/types"
)

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

func tokenizeFieldList(fields []types.Field) []Token {
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
