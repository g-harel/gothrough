package index_test

import (
	"reflect"
	"strings"
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
		idx.Index(ids[0], 2, query)
		idx.Index(ids[1], 1, query)
		res := idx.Search(query)

		testContentEqual(t, res, ids)
	})

	t.Run("should only return matching values", func(t *testing.T) {
		ids := []int{54, 76}
		query := "only_matching"

		idx := index.NewIndex()
		idx.Index(ids[0], 2, query)
		idx.Index(ids[1], 1, "%")
		res := idx.Search(query)

		testContentEqual(t, res, []int{ids[0]})
	})

	t.Run("should return matched values in order of confidence", func(t *testing.T) {
		ids := []int{21, 82}
		query := "matching_order_confidence"

		idx := index.NewIndex()
		idx.Index(ids[0], 100, query)
		idx.Index(ids[1], 50, query)
		res := idx.Search(query)

		testContentEqual(t, res, ids)
	})

	t.Run("should accumulate confidence from multiple index calls", func(t *testing.T) {
		ids := []int{81, 43}
		query := "confidence_order_sum"

		idx := index.NewIndex()
		idx.Index(ids[0], 100, query)
		idx.Index(ids[1], 60, query)
		idx.Index(ids[1], 60, query)
		res := idx.Search(query)

		testContentEqual(t, res, []int{ids[1], ids[0]})
	})

	t.Run("should accumulate confidence from multiple query parts", func(t *testing.T) {
		ids := []int{43, 91}
		queries := []string{"query", "parts"}

		idx := index.NewIndex()
		idx.Index(ids[0], 60, queries...)
		idx.Index(ids[1], 100, queries[0])
		res := idx.Search(strings.Join(queries, " "))

		testContentEqual(t, res, ids)
	})
}
