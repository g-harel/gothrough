package pages

import (
	"net/http"

	"github.com/g-harel/gis/internal/templates"
)

func Home() http.HandlerFunc {
	return templates.NewRenderer(nil, "pages/home.html").Handler
}