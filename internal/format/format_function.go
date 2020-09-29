package format

import (
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

// TODO primary one should be linked in top-lever format and add func keyword.
func formatFunction(function *types.Function) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	snippet.Push(formatDocs(&function.Docs))

	snippet.FunctionName(function.Name)
	snippet.Punctuation("(")
	snippet.Push(formatFieldList(function.Arguments))
	snippet.Punctuation(")")

	if len(function.ReturnValues) == 1 {
		snippet.Space()
		snippet.Push(formatFieldList(function.ReturnValues))
	} else if len(function.ReturnValues) > 1 {
		snippet.Space()
		snippet.Punctuation("(")
		snippet.Push(formatFieldList(function.ReturnValues))
		snippet.Punctuation(")")
	}

	snippet.Newline()

	return snippet
}
