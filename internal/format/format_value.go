package format

import (
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func formatValue(value *types.Value) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	snippet.Push(formatDocs(&value.Docs))

	snippet.Keyword("const")
	snippet.Space()

	// TODO more specific kind
	snippet.FunctionName(value.Name)
	snippet.Space()
	snippet.Punctuation("=")
	snippet.Space()
	snippet.Literal(value.Value)

	snippet.Newline()

	return snippet
}
