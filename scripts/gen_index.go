package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gothrough/internal/interface_index"
	"github.com/g-harel/gothrough/internal/parse"
)

var dest = flag.String("dest", ".index", "output filename")

func usageErr() {
	fmt.Printf("Usage:\n  %s -dest=DEST_DIR SRC_DIR...\n", os.Args[0])
	os.Exit(1)
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "fatal: %s", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	// Validate inputs.
	if len(flag.Args()) == 0 {
		usageErr()
	}

	// Create index.
	idx := interface_index.NewIndex()
	for _, dir := range flag.Args() {
		interfaces, err := parse.FindInterfaces(path.Join(dir, "src"))
		if err != nil {
			fatalErr(err)
		}

		for _, ifc := range interfaces {
			idx.Insert(*ifc)
		}
	}

	// Create output file.
	f, err := os.Create(*dest)
	if err != nil {
		fatalErr(err)
	}
	defer f.Close()

	// Write index to ouput file.
	err = idx.ToBytes(f)
	if err != nil {
		fatalErr(err)
	}
}
