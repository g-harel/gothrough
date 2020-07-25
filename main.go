package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/gis/internal/interface_index"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		interfaces, err := idx.Search("io reader")
		if err != nil {
			panic(err)
		}

		for _, ifc := range interfaces[:16] {
			fmt.Fprintf(w, "%s\n", ifc.String())
		}
	})

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
