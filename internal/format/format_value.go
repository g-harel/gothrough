package format

import (
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func formatValue(value *types.Value) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	snippet.Push(formatDocs(&value.Docs))

	if value.Const {
		snippet.Keyword("const")
	} else {
		snippet.Keyword("var")
	}
	snippet.Space()

	snippet.DeclName(value.Name)
	snippet.Space()

	if value.Value != "" {
		snippet.Punctuation("=")
		snippet.Space()
		snippet.Literal(value.Value)
	} else {
		snippet.FieldType(value.Type)
	}

	snippet.Newline()

	return snippet
}
