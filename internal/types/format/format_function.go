package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatFunction(function *types.Function, decl bool) *Snippet {
	snippet := NewSnippet()

	snippet.concat(formatDocs(&function.Docs))

	if decl {
		snippet.keyword("func")
		snippet.space()
	}

	snippet.declName(function.Name)
	snippet.punctuation("(")
	snippet.concat(formatFieldList(function.Arguments))
	snippet.punctuation(")")

	hasNamedReturnValue := false
	for _, returnValue := range function.ReturnValues {
		if returnValue.Name != "" {
			hasNamedReturnValue = true
			break
		}
	}

	if len(function.ReturnValues) == 1 && !hasNamedReturnValue {
		snippet.space()
		snippet.concat(formatFieldList(function.ReturnValues))
	} else if len(function.ReturnValues) > 1 || hasNamedReturnValue {
		snippet.space()
		snippet.punctuation("(")
		snippet.concat(formatFieldList(function.ReturnValues))
		snippet.punctuation(")")
	}

	snippet.newline()

	return snippet
}
