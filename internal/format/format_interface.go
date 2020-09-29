package format

import (
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func formatInterface(ifc *types.Interface) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	snippet.Push(formatDocs(&ifc.Docs))

	snippet.Keyword("type")
	snippet.Space()
	snippet.InterfaceName(ifc.Name)
	snippet.Space()
	snippet.Keyword("interface")
	snippet.Space()
	snippet.Punctuation("{")

	if len(ifc.Embedded) == 0 && len(ifc.Methods) == 0 {
		snippet.Punctuation("}")
		snippet.Newline()
		return snippet
	}

	snippet.Newline()

	for i, embedded := range ifc.Embedded {
		// Add newline before definition in some situations.
		prevHadDocs := i > 0 && ifc.Embedded[i-1].Docs.Text != ""
		isNotFirstAndHasDocs := i != 0 && embedded.Docs.Text != ""
		if prevHadDocs || isNotFirstAndHasDocs {
			snippet.Newline()
		}

		embeddedDocs := formatDocs(&embedded.Docs)
		embeddedDocs.IndentLines()
		snippet.Push(embeddedDocs)

		snippet.Indent()
		if embedded.Package != "" {
			snippet.EmbeddedPackage(embedded.Package)
			snippet.Punctuation(".")
		}
		snippet.EmbeddedName(embedded.Name)
		snippet.Newline()
	}

	for i, method := range ifc.Methods {
		// Add newline before definition in some situations.
		prevWasEmbedded := i == 0 && len(ifc.Embedded) > 0
		prevEmbeddedHadDocs := len(ifc.Embedded) > 0 && ifc.Embedded[len(ifc.Embedded)-1].Docs.Text != ""
		prevMethodHadDocs := i > 0 && ifc.Methods[i-1].Docs.Text != ""
		isNotFirstAndHasDocs := (i != 0 || prevWasEmbedded) && method.Docs.Text != ""
		if prevEmbeddedHadDocs || prevMethodHadDocs || isNotFirstAndHasDocs {
			snippet.Newline()
		}

		formattedMethod := formatFunction(&method)
		formattedMethod.IndentLines()
		snippet.Push(formattedMethod)
	}

	snippet.Punctuation("}")
	snippet.Newline()

	return snippet
}
