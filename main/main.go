package main

import (
	"os"
	"path"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")

	_, err := gis.Search(root)
	if err != nil {
		panic(err)
	}
}
