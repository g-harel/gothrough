package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatFunction(function *types.Function, decl bool) *Snippet {
	snippet := NewSnippet()

	snippet.Push(formatDocs(&function.Docs))

	if decl {
		snippet.Keyword("func")
		snippet.Space()
	}

	snippet.DeclName(function.Name)
	snippet.Punctuation("(")
	snippet.Push(formatFieldList(function.Arguments))
	snippet.Punctuation(")")

	hasNamedReturnValue := false
	for _, returnValue := range function.ReturnValues {
		if returnValue.Name != "" {
			hasNamedReturnValue = true
			break
		}
	}

	if len(function.ReturnValues) == 1 && !hasNamedReturnValue {
		snippet.Space()
		snippet.Push(formatFieldList(function.ReturnValues))
	} else if len(function.ReturnValues) > 1 || hasNamedReturnValue {
		snippet.Space()
		snippet.Punctuation("(")
		snippet.Push(formatFieldList(function.ReturnValues))
		snippet.Punctuation(")")
	}

	snippet.Newline()

	return snippet
}
