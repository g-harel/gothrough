package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	query := "io reader"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx, err := gis.NewSearchIndex(root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Indexed in %s\n", time.Since(indexTime))

	searchStart := time.Now()
	interfaces, err := idx.Search(query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Searched in %s\n", time.Since(searchStart))
	fmt.Println("========")

	for _, ifc := range interfaces {
		println(ifc.Pretty())
	}
}
