package pages

import (
	"log"
	"net/http"

	"github.com/g-harel/gothrough/internal/templates"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types/format"
)

// PrettyResult wraps the index results into a template-friendly struct.
type PrettyResult struct {
	Name              string
	Confidence        int
	PackageName       string
	PackageImportPath string
	PrettyTokens      []format.Token
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

// Home returns a baked handler for the homepage.
func Home(packages [][]string) http.HandlerFunc {
	context := struct {
		Packages [][]string
	}{
		Packages: packages,
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/home.html").Handler
}

// Results returns a baked handler for the results page.
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
