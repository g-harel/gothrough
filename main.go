package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/pages"
)

func headerResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%d %s", statusCode, http.StatusText(statusCode))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	indexFilename := os.Getenv("INDEX")
	if indexFilename == "" {
		indexFilename = ".index"
	}

	fi, err := os.Stat(indexFilename)
	if err != nil {
		panic(fmt.Sprintf("stat index file: ", err))
	}
	log.Printf("Index size: %v bytes", fi.Size())

	f, err := os.Open(indexFilename)
	if err != nil {
		panic(fmt.Sprintf("open index file: ", err))
	}

	idx, err := typeindex.NewIndexFromBytes(f)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home(idx.Packages())(w, r)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			headerResponse(w, http.StatusBadRequest)
			return
		}

		query := r.Form.Get("q")
		if query == "" {
			headerResponse(w, http.StatusBadRequest)
			return
		}

		results, err := idx.Search(query)
		if err != nil {
			panic(err)
		}

		if len(results) > 16 {
			results = results[:16]
		}

		pages.Results(query, results)(w, r)
	})

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
