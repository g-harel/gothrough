package main

import (
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gis"
)

func main() {
	root := path.Join(os.Getenv("GOROOT"), "src")
	query := "read write"

	fmt.Printf("ROOT=%v\n", root)
	fmt.Printf("QUERY=%v\n", query)
	fmt.Println("========")

	interfaces, err := gis.Search(root, query)
	if err != nil {
		panic(err)
	}

	for _, ifc := range interfaces {
		fmt.Printf("%v %v (%v)\n", ifc.Name, ifc.Methods, ifc.Address())
	}
}
