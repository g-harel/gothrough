package pages

import (
	"net/http"

	"github.com/g-harel/gis/internal/interfaces"
	"github.com/g-harel/gis/internal/templates"
)

func Home() http.HandlerFunc {
	return templates.NewRenderer(nil, "pages/_layout.html", "pages/home.html").Handler
}

type ResultsResult struct {
	Name        string
	PackageName string
	Pretty      string
}

func Results(query string, interfaces []interfaces.Interface) http.HandlerFunc {
	context := struct {
		Query   string
		Results []ResultsResult
	}{
		Query:   query,
		Results: []ResultsResult{},
	}
	for _, ifc := range interfaces {
		context.Results = append(context.Results, ResultsResult{
			Name:        ifc.Name,
			PackageName: ifc.PackageName,
			Pretty:      ifc.Pretty(),
		})
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
