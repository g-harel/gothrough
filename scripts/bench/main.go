package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/format"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	path := path.Join(os.Getenv("GOPATH"), "src")
	dest := "./.index"
	query := "ios reder option"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("PATH=%v\n", path)
	fmt.Printf("DEST=%v\n", dest)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx := typeindex.NewIndex()

	err := extract.Types(root, extract.TypeHandlers{
		Interface: idx.InsertInterface,
		Function: func(fn types.Function) {
			// TODO
		},
	})
	if err != nil {
		panic(err)
	}
	err = extract.Types(path, extract.TypeHandlers{
		Interface: idx.InsertInterface,
		Function: func(fn types.Function) {
			// TODO
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed in %s\n", time.Since(indexTime))

	encodeTime := time.Now()
	var buf bytes.Buffer
	err = idx.ToBytes(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encoded in %s\n", time.Since(encodeTime))

	decodeTime := time.Now()
	idx, err = typeindex.NewIndexFromBytes(&buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded in %s\n", time.Since(decodeTime))

	writeTime := time.Now()
	f, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = idx.ToBytes(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Written in %s\n", time.Since(writeTime))

	searchStart := time.Now()
	results, err := idx.Search(query)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Queried in %s\n", time.Since(searchStart))
	fmt.Println("========")

	for _, result := range results[:8] {
		fmt.Printf("\n// === %.6f ===\n", result.Confidence)
		p, err := format.String(result.Value)
		if err != nil {
			panic(err)
		}
		println(p)
	}
}
