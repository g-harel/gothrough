package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gis/internal/interface_index"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	path := path.Join(os.Getenv("GOPATH"), "src")
	dest := "./.index"
	query := "ios reder option"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("PATH=%v\n", path)
	fmt.Printf("DEST=%v\n", dest)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	indexTime := time.Now()
	idx := interface_index.NewIndex()
	err := idx.Include(root)
	if err != nil {
		panic(err)
	}
	err = idx.Include(path)
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
	idx, err = interface_index.NewIndexFromBytes(&buf)
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

	for _, result := range results[:16] {
		println(result.Interface.Pretty())
	}
}
