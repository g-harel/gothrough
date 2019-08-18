package main

import (
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	query := "io reader"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	idx, err := gis.NewSearchIndex(root)
	if err != nil {
		panic(err)
	}

	interfaces, err := idx.Search(query)
	if err != nil {
		panic(err)
	}

	for _, ifc := range interfaces {
		println(ifc.String())
	}
}
