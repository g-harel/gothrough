package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types/format"
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

	// Create index.
	indexTimer := NewTimer("index")
	widx := typeindex.NewIndex()
	for _, dir := range flag.Args() {
		err := extract.Types(path.Join(dir, "src"), extract.TypeHandlers{
			Interface: widx.InsertInterface,
			Function:  widx.InsertFunction,
			Value:     widx.InsertValue,
		})
		if err != nil {
			fatalErr(fmt.Errorf("extract types: %v", err))
		}
	}
	indexTimer.Done()

	// Create output file.
	writeTimer := NewTimer("encode/write")
	wf, err := os.Create(*dest)
	if err != nil {
		fatalErr(fmt.Errorf("create index file: %v", err))
	}
	defer wf.Close()

	// Write index to output file.
	err = widx.ToBytes(wf)
	if err != nil {
		fatalErr(fmt.Errorf("write index to file: %v", err))
	}
	writeTimer.Done()

	readTimer := NewTimer("read/decode")
	rf, err := os.Open(*dest)
	if err != nil {
		fatalErr(fmt.Errorf("open index file: %v", err))
	}

	// Read index from output file.
	ridx, err := typeindex.NewIndexFromBytes(rf)
	if err != nil {
		fatalErr(fmt.Errorf("create index form file: %v", err))
	}
	readTimer.Done()

	if *query != "" {
		searchTimer := NewTimer("search")
		results, err := ridx.Search(*query)
		if err != nil {
			fatalErr(fmt.Errorf("search index: %v", err))
		}
		searchTimer.Done()

		fmt.Println("========")

		if len(results) > 8 {
			results = results[:8]
		}

		for _, result := range results {
			fmt.Printf("=== %.6f ===\n", result.Confidence)
			p, err := format.String(result.Value)
			if err != nil {
				fatalErr(fmt.Errorf("format result: %v", err))
			}
			println(p)
		}
	}
}
