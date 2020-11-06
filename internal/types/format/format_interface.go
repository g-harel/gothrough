package format

import (
	"github.com/g-harel/gothrough/internal/types"
)

func formatInterface(ifc *types.Interface, decl bool) *Snippet {
	snippet := NewSnippet()

	snippet.concat(formatDocs(&ifc.Docs))

	if decl {
		snippet.keyword("type")
		snippet.space()
	}

	snippet.interfaceName(ifc.Name)
	snippet.space()
	snippet.keyword("interface")
	snippet.space()
	snippet.punctuation("{")

	if len(ifc.Embedded) == 0 && len(ifc.Methods) == 0 {
		snippet.punctuation("}")
		snippet.newline()
		return snippet
	}

	snippet.newline()

	for i, embedded := range ifc.Embedded {
		// Add newline before definition in some situations.
		prevHadDocs := i > 0 && ifc.Embedded[i-1].Docs.Text != ""
		isNotFirstAndHasDocs := i != 0 && embedded.Docs.Text != ""
		if prevHadDocs || isNotFirstAndHasDocs {
			snippet.newline()
		}

		embeddedDocs := formatDocs(&embedded.Docs)
		embeddedDocs.indentSnippet()
		snippet.concat(embeddedDocs)

		snippet.indent()
		if embedded.Package != "" {
			snippet.embeddedPackage(embedded.Package)
			snippet.punctuation(".")
		}
		snippet.embeddedName(embedded.Name)
		snippet.newline()
	}

	for i, method := range ifc.Methods {
		// Add newline before definition in some situations.
		prevWasEmbedded := i == 0 && len(ifc.Embedded) > 0
		prevEmbeddedHadDocs := prevWasEmbedded && len(ifc.Embedded) > 0 && ifc.Embedded[len(ifc.Embedded)-1].Docs.Text != ""
		prevMethodHadDocs := i > 0 && ifc.Methods[i-1].Docs.Text != ""
		isNotFirstAndHasDocs := (i != 0 || prevWasEmbedded) && method.Docs.Text != ""
		if prevEmbeddedHadDocs || prevMethodHadDocs || isNotFirstAndHasDocs {
			snippet.newline()
		}

		formattedMethod := formatFunction(&method, false)
		formattedMethod.indentSnippet()
		snippet.concat(formattedMethod)
	}

	snippet.punctuation("}")
	snippet.newline()

	return snippet
}
