package integration_test

import (
	"testing"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/typeindex"
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
}
