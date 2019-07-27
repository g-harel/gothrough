package main

import (
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")

	interfaces, err := gis.Search(root)
	if err != nil {
		panic(err)
	}

	for _, ifc := range interfaces {
		fmt.Println(ifc.String())
	}
	fmt.Println(len(interfaces))
}
