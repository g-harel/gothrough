package pages

import (
	"net/http"

	"github.com/g-harel/gothrough/internal/source_index"
	"github.com/g-harel/gothrough/internal/templates"
	"github.com/g-harel/gothrough/internal/types"
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
	PrettyTokens      []types.Token
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
		context.Results = append(context.Results, ResultsResult{
			Name:              result.Name,
			PackageName:       result.PackageName,
			PackageImportPath: result.PackageImportPath,
			PrettyTokens:      result.Value.PrettyTokens(),
		})
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
