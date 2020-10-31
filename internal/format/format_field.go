package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatField(field *types.Field) *Snippet {
	snippet := NewSnippet()

	if field.Name != "" {
		snippet.FieldName(field.Name)
		snippet.Space()
	}
	snippet.FieldType(field.Type)

	return snippet
}

func formatFieldList(fields []types.Field) *Snippet {
	snippet := NewSnippet()

	for i, field := range fields {
		if i > 0 {
			snippet.Punctuation(",")
			snippet.Space()
		}
		snippet.Push(formatField(&field))
	}

	return snippet
}
