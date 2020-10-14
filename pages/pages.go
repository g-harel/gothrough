package pages

import (
	"log"
	"net/http"

	"github.com/g-harel/gothrough/internal/format"
	"github.com/g-harel/gothrough/internal/templates"
	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/typeindex"
)

type PrettyResult struct {
	Name              string
	Confidence        int
	PackageName       string
	PackageImportPath string
	PrettyTokens      []tokens.Token
}

func formatAll(results []*typeindex.Result) []PrettyResult {
	formatted := []PrettyResult{}
	for _, result := range results {
		snippet, err := format.Format(result.Value)
		if err != nil {
			log.Printf("could not format: %v", err)
			continue
		}
		formatted = append(formatted, PrettyResult{
			Name:              result.Name,
			Confidence:        int(result.Confidence * 1000),
			PackageName:       result.Location.PackageName,
			PackageImportPath: result.Location.PackageImportPath,
			PrettyTokens:      snippet.Dump(),
		})
	}
	return formatted
}

func Home(packages [][]string) http.HandlerFunc {
	context := struct {
		Packages [][]string
	}{
		Packages: packages,
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/home.html").Handler
}

func Results(query string, results []*typeindex.Result) http.HandlerFunc {
	context := struct {
		Query   string
		Results []PrettyResult
	}{
		Query:   query,
		Results: formatAll(results),
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
