package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types"
)

var dest = flag.String("dest", ".index", "output filename")

func usageErr() {
	fmt.Printf("Usage:\n  %s -dest=DEST_DIR SRC_DIR...\n", os.Args[0])
	os.Exit(1)
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	// Validate inputs.
	if len(flag.Args()) == 0 {
		usageErr()
	}

	// Create index.
	// TODO index constants.
	idx := typeindex.NewIndex()
	for _, dir := range flag.Args() {
		err := extract.Types(path.Join(dir, "src"), extract.TypeHandlers{
			Interface: idx.InsertInterface,
			Function: func(fn types.Function) {
				// TODO
			},
		})
		if err != nil {
			fatalErr(err)
		}
	}

	// Create output file.
	f, err := os.Create(*dest)
	if err != nil {
		fatalErr(err)
	}
	defer f.Close()

	// Write index to output file.
	err = idx.ToBytes(f)
	if err != nil {
		fatalErr(err)
	}
}
