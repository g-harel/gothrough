package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/gis/internal/interface_index"
)

// http://localhost:3000/?query=io%20reader
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.String())

		query, ok := r.URL.Query()["query"]
		if !ok || len(query) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		results, err := idx.Search(query[0])
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "%s\n========\n", query[0])
		if len(results) > 16 {
			results = results[:16]
		}
		for _, result := range results {
			fmt.Fprintf(w, "%4.3f %s\n", result.Confidence, result.Interface.String())
		}
	})

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
