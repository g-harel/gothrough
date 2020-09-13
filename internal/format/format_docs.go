package format

import (
	"strings"

	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func formatDocs(docs *types.Docs) *tokens.Snippet {
	snippet := tokens.NewSnippet()

	if docs.Text == "" {
		return snippet
	}

	docLines := strings.Split(strings.TrimSpace(docs.Text), "\n")
	if len(docLines) < 1 {
		return snippet
	}

	for _, line := range docLines {
		snippet.Comment("// " + line)
		snippet.Newline()
	}

	return snippet
}
