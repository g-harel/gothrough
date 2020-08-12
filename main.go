package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/gis/internal/interface_index"
	"github.com/g-harel/gis/pages"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	indexFilename := os.Getenv("INDEX")
	if indexFilename == "" {
		indexFilename = ".index"
	}
	f, err := os.Open(indexFilename)
	if err != nil {
		panic("missing index file")
	}

	idx, err := interface_index.NewIndexFromBytes(f)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", pages.Home())

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		headerResponse := func(statusCode int) {
			w.WriteHeader(statusCode)
			fmt.Fprintf(w, "%d %s", statusCode, http.StatusText(statusCode))
		}

		err := r.ParseForm()
		if err != nil {
			headerResponse(http.StatusBadRequest)
			return
		}

		query := r.Form.Get("query")
		if query == "" {
			headerResponse(http.StatusBadRequest)
			return
		}

		results, err := idx.Search(query)
		if err != nil {
			panic(err)
		}

		prettyResults := []string{}
		if len(results) > 16 {
			results = results[:16]
		}
		for _, result := range results {
			prettyResult := fmt.Sprintf("\n// === %.6f === %v\n", result.Confidence, result.Interface.DocLink())
			prettyResult += result.Interface.Pretty()
			prettyResults = append(prettyResults, prettyResult)
		}

		pages.Results(query, prettyResults)(w, r)
	})

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
