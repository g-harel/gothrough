package string_index_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/g-harel/gis/internal/string_index"
)

type indexItem struct {
	id           int
	confidence   int
	matchStrings []string
}

func indexFrom(items ...indexItem) *string_index.Index {
	idx := string_index.NewIndex()
	for _, item := range items {
		idx.Index(item.id, item.confidence, item.matchStrings...)
	}
	return idx
}

func assertEqual(t *testing.T, msg string, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual and expected %v do not match\n%v\n%v", msg, actual, expected)
	}
}

func TestIndex(t *testing.T) {
	t.Run("should return indexed value", func(t *testing.T) {
		query := "return_indexed_value"
		item := indexItem{32, 0, []string{query}}

		idx := indexFrom(item)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 1)
		assertEqual(t, "id", actual[0].ID, item.id)
	})

	t.Run("should return multiple matching values", func(t *testing.T) {
		query := "multiple_indexed_values"
		item0 := indexItem{12, 0, []string{query}}
		item1 := indexItem{78, 0, []string{query}}

		idx := indexFrom(item0, item1)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, item0.id)
		assertEqual(t, "second id", actual[1].ID, item1.id)
	})

	t.Run("should only return matching values", func(t *testing.T) {
		query := "only_matching"
		item0 := indexItem{54, 0, []string{query}}
		item1 := indexItem{76, 0, []string{"%"}}

		idx := indexFrom(item0, item1)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 1)
		assertEqual(t, "id", actual[0].ID, item0.id)
	})

	t.Run("should return partially matching values", func(t *testing.T) {
		query := "abc xy"
		item0 := indexItem{98, 0, []string{"ab"}}
		item1 := indexItem{81, 0, []string{"xyz"}}
		item2 := indexItem{123456, 0, []string{"*"}}

		idx := indexFrom(item0, item1, item2)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, item0.id)
		assertEqual(t, "second id", actual[1].ID, item1.id)
	})

	t.Run("should return matched values in order of confidence", func(t *testing.T) {
		query := "matching_order_confidence"
		item0 := indexItem{21, 50, []string{query}}
		item1 := indexItem{82, 100, []string{query}}

		idx := indexFrom(item0, item1)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, item1.id)
		assertEqual(t, "second id", actual[1].ID, item0.id)
	})

	t.Run("should accumulate confidence from multiple index calls", func(t *testing.T) {
		query := "confidence_order_sum"
		item0 := indexItem{81, 100, []string{query}}
		item1 := indexItem{43, 60, []string{query}}
		item2 := item1

		idx := indexFrom(item0, item1, item2)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, item1.id)
		assertEqual(t, "second id", actual[1].ID, item0.id)
	})

	t.Run("should accumulate confidence from multiple query parts", func(t *testing.T) {
		queries := []string{"query", "parts"}
		item0 := indexItem{43, 60, queries}
		item1 := indexItem{91, 100, []string{queries[0]}}

		idx := indexFrom(item0, item1)
		actual := idx.Search(strings.Join(queries, " "))

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, item0.id)
		assertEqual(t, "second id", actual[1].ID, item1.id)
	})
}
