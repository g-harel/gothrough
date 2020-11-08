package integration_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/typeindex"
	"github.com/g-harel/gothrough/internal/types/format"
)

func TestIndex(t *testing.T) {
	idx := typeindex.NewIndex()
	err := extract.Types("src", extract.TypeHandlers{
		Interface: idx.InsertInterface,
		Function:  idx.InsertFunction,
		Value:     idx.InsertValue,
	})
	if err != nil {
		t.Fatalf("extract types: %v", err)
	}

	t.Run("should find results", func(t *testing.T) {
		results, err := idx.Search("")
		if err != nil {
			t.Fatalf("search index: %v", err)
		}

		if len(results) < 1 {
			t.Fatalf("expected results but got none.")
		}
	})

	t.Run("should be able to filter results", func(t *testing.T) {
		results, err := idx.Search("type:const")
		if err != nil {
			t.Fatalf("search index: %v", err)
		}

		for _, result := range results {
			pretty, err := format.String(result.Value)
			if err != nil {
				t.Fatalf("format result: %v", err)
			}
			if !strings.Contains(pretty, "const") {
				t.Fatalf("a result was not a const:\n%s", pretty)
			}
		}
	})

	t.Run("should be able to filter out results", func(t *testing.T) {
		results, err := idx.Search("-type:interface")
		if err != nil {
			t.Fatalf("search index: %v", err)
		}

		for _, result := range results {
			pretty, err := format.String(result.Value)
			if err != nil {
				t.Fatalf("format result: %v", err)
			}
			if strings.Contains(pretty, "interface") {
				t.Fatalf("a result was an interface:\n%s", pretty)
			}
		}
	})

	t.Run("should rank by match quality", func(t *testing.T) {
		results, err := idx.Search("test")
		if err != nil {
			t.Fatalf("search index: %v", err)
		}

		if results[0].Name != "Test" {
			pretty, err := format.String(results[0].Value)
			if err != nil {
				t.Fatalf("format result: %v", err)
			}
			for _, result := range results {
				pretty, _ := format.String(result.Value)
				fmt.Printf("%.6f\n", result.Confidence)
				println(pretty)
			}
			t.Fatalf("did not get expected result order:\n%s", pretty)
		}
	})
}
