package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gothrough/internal/parse"
	"github.com/g-harel/gothrough/internal/source_index"
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
	idx := source_index.NewIndex()
	rootInterfaces, err := parse.FindInterfaces(root)
	if err != nil {
		panic(err)
	}
	for _, ifc := range rootInterfaces {
		idx.Insert(*ifc)
	}
	pathInterfaces, err := parse.FindInterfaces(path)
	if err != nil {
		panic(err)
	}
	for _, ifc := range pathInterfaces {
		idx.Insert(*ifc)
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
	idx, err = source_index.NewIndexFromBytes(&buf)
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
		println(result.Interface.Pretty())
	}
}
