package index_test

import (
	"reflect"
	"testing"

	"github.com/g-harel/gis/internal/index"
)

func testContentEqual(t *testing.T, actual, expected []int) {
	copyActual := make([]int, len(actual))
	copy(copyActual, actual)

	copyExpected := make([]int, len(expected))
	copy(copyExpected, expected)

	if !reflect.DeepEqual(copyActual, copyExpected) {
		t.Fatalf("actual and expected do not match\n%v\n%v", actual, expected)
	}
}

func TestIndex(t *testing.T) {
	t.Run("should return indexed value", func(t *testing.T) {
		id := 32
		query := "return_indexed_value"

		idx := index.NewIndex()
		idx.Index(id, 0, query)
		res := idx.Search(query)

		testContentEqual(t, res, []int{id})
	})

	t.Run("should return multiple matching values", func(t *testing.T) {
		ids := []int{12, 78}
		query := "multiple_indexed_values"

		idx := index.NewIndex()
		idx.Index(ids[0], 0, query)
		idx.Index(ids[1], 0, query)
		res := idx.Search(query)

		testContentEqual(t, res, ids)
	})

	t.Run("should only return matching values", func(t *testing.T) {
		ids := []int{54, 76}
		query := "only_matching"

		idx := index.NewIndex()
		idx.Index(ids[0], 0, query)
		idx.Index(ids[1], 0, "_")
		res := idx.Search(query)

		testContentEqual(t, res, []int{ids[0]})
	})
}
