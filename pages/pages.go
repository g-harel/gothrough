package pages

import (
	"net/http"

	"github.com/g-harel/gothrough/internal/templates"
	"github.com/g-harel/gothrough/internal/types"
)

func Home() http.HandlerFunc {
	return templates.NewRenderer(nil, "pages/_layout.html", "pages/home.html").Handler
}

type ResultsResult struct {
	Name              string
	PackageName       string
	PackageImportPath string
	PrettyTokens      []types.Token
}

func Results(query string, interfaces []types.Interface) http.HandlerFunc {
	context := struct {
		Query   string
		Results []ResultsResult
	}{
		Query:   query,
		Results: []ResultsResult{},
	}
	for _, ifc := range interfaces {
		context.Results = append(context.Results, ResultsResult{
			Name:              ifc.Name,
			PackageName:       ifc.PackageName,
			PackageImportPath: ifc.PackageImportPath,
			PrettyTokens:      ifc.PrettyTokens(),
		})
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
