package pages

import (
	"net/http"

	"github.com/g-harel/gis/internal/templates"
)

func Home() http.HandlerFunc {
	return templates.NewRenderer(nil, "pages/_layout.html", "pages/home.html").Handler
}

func Results(query string, interfaces []string) http.HandlerFunc {
	context := struct {
		Query   string
		Results []string
	}{
		Query:   query,
		Results: interfaces,
	}
	return templates.NewRenderer(context, "pages/_layout.html", "pages/results.html").Handler
}
