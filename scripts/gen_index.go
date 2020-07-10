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
	// TODO make root and dest cli args.
	root := path.Join(os.Getenv("GOROOT"), "src")
	dest := "./index.bin"
	// TODO move query to different script + server.
	query := "ios read"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("DEST=%v\n", dest)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx, err := gis.NewSearchIndexFromSource(root)
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

	// Write to disk.
	f, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = idx.ToBytes(f)
	if err != nil {
		panic(err)
	}

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
