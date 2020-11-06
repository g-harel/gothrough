package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatField(field *types.Field) *Snippet {
	snippet := NewSnippet()

	if field.Name != "" {
		snippet.fieldName(field.Name)
		snippet.space()
	}
	snippet.fieldType(field.Type)

	return snippet
}

func formatFieldList(fields []types.Field) *Snippet {
	snippet := NewSnippet()

	for i, field := range fields {
		if i > 0 {
			snippet.punctuation(",")
			snippet.space()
		}
		snippet.concat(formatField(&field))
	}

	return snippet
}
