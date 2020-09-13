package format

import (
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func formatField(field *types.Field) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	if field.Name != "" {
		snippet.FieldName(field.Name)
		snippet.Space()
	}
	snippet.FieldType(field.Type)

	return snippet
}

func formatFieldList(fields []types.Field) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	for i, field := range fields {
		if i > 0 {
			snippet.Punctuation(",")
			snippet.Space()
		}
		snippet.Push(formatField(&field))
	}

	return snippet
}
