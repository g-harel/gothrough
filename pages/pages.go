package pages

import (
	"log"
	"net/http"

	"github.com/g-harel/gothrough/internal/format"
	"github.com/g-harel/gothrough/internal/source_index"
	"github.com/g-harel/gothrough/internal/templates"
	"github.com/g-harel/gothrough/internal/tokens"
)

func Home(packages [][]string) http.HandlerFunc {
	context := struct {
		Packages [][]string
	}{
		Packages: packages,
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/home.html").Handler
}

type ResultsResult struct {
	Name              string
	PackageName       string
	PackageImportPath string
	PrettyTokens      []tokens.Token
}

func Results(query string, results []*source_index.Result) http.HandlerFunc {
	context := struct {
		Query   string
		Results []ResultsResult
	}{
		Query:   query,
		Results: []ResultsResult{},
	}
	for _, result := range results {
		snippet, err := format.Format(result.Value)
		if err != nil {
			log.Printf("could not format: %v", err)
			continue
		}
		context.Results = append(context.Results, ResultsResult{
			Name:              result.Name,
			PackageName:       result.PackageName,
			PackageImportPath: result.PackageImportPath,
			PrettyTokens:      snippet.Dump(),
		})
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
