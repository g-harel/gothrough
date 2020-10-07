package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/format"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types"
)

var dest = flag.String("dest", ".index", "output filename")
var query = flag.String("query", "", "query the index")

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

	indexTime := time.Now()

	// Create index.
	// TODO index constants.
	idx := typeindex.NewIndex()
	for _, dir := range flag.Args() {
		err := extract.Types(path.Join(dir, "src"), extract.TypeHandlers{
			Interface: idx.InsertInterface,
			Function:  idx.InsertFunction,
			Value: func(location extract.Location, value types.Value) {
				// TODO
				fmt.Println(value)
			},
		})
		if err != nil {
			fatalErr(err)
		}
	}

	fmt.Printf("Indexed in %s\n", time.Since(indexTime))

	writeTime := time.Now()

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

	fmt.Printf("Written in %s\n", time.Since(writeTime))

	if *query != "" {
		searchStart := time.Now()

		results, err := idx.Search(*query)
		if err != nil {
			fatalErr(err)
		}

		fmt.Printf("Queried in %s\n", time.Since(searchStart))
		fmt.Println("========")

		if len(results) > 8 {
			results = results[:8]
		}

		for _, result := range results {
			fmt.Printf("=== %.6f ===\n", result.Confidence)
			p, err := format.String(result.Value)
			if err != nil {
				fatalErr(err)
			}
			println(p)
		}
	}
}
