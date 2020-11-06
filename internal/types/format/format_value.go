package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatValue(value *types.Value) *Snippet {
	snippet := NewSnippet()

	snippet.concat(formatDocs(&value.Docs))

	if value.Const {
		snippet.keyword("const")
	} else {
		snippet.keyword("var")
	}
	snippet.space()

	snippet.declName(value.Name)
	snippet.space()

	if value.Value != "" {
		snippet.punctuation("=")
		snippet.space()
		snippet.literal(value.Value)
	} else {
		snippet.fieldType(value.Type)
	}

	snippet.newline()

	return snippet
}
