package index_test

import (
	"testing"

	"github.com/g-harel/gis/internal/index"
)

func TestIndex(t *testing.T) {
	t.Run("should return indexed value", func(t *testing.T) {
		id := 32
		query := "return_indexed_value"

		idx := index.NewIndex()
		idx.Index(id, 0, query)
		res := idx.Search(query)

		if len(res) < 1 || res[0] != id {
			t.Fatalf("indexed value was not found")
		}
	})
}
