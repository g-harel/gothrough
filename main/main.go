package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	query := "ios read"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx, err := gis.NewSearchIndex(root)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed in %s\n", time.Since(indexTime))

	conversionTime := time.Now()
	var buf bytes.Buffer
	err = idx.ToBytes(&buf)
	if err != nil {
		panic(err)
	}
	idx, err = gis.NewSearchIndexFromBytes(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encoded/decoded in %s\n", time.Since(conversionTime))

	searchStart := time.Now()
	interfaces, err := idx.Search(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Searched in %s\n", time.Since(searchStart))
	fmt.Println("========")

	for _, ifc := range interfaces[:16] {
		println(ifc.String())
	}
}
