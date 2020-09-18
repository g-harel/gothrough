package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/pages"
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

	idx, err := typeindex.NewIndexFromBytes(f)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home(idx.Packages())(w, r)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		headerResponse := func(statusCode int) {
			w.WriteHeader(statusCode)
			fmt.Fprintf(w, "%d %s", statusCode, http.StatusText(statusCode))
		}

		err := r.ParseForm()
		if err != nil {
			headerResponse(http.StatusBadRequest)
			return
		}

		query := r.Form.Get("q")
		if query == "" {
			headerResponse(http.StatusBadRequest)
			return
		}

		results, err := idx.Search(query)
		if err != nil {
			panic(err)
		}

		pages.Results(query, results)(w, r)
	})

	// TODO package viewer

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
