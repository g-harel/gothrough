package string_index_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/g-harel/gis/internal/string_index"
)

// TODO util to index matches clearer.

func assertEqual(t *testing.T, msg string, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual and expected %v do not match\n%v\n%v", msg, actual, expected)
	}
}

func TestIndex(t *testing.T) {
	t.Run("should return indexed value", func(t *testing.T) {
		id := 32
		query := "return_indexed_value"

		idx := string_index.NewIndex()
		idx.Index(id, 0, query)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 1)
		assertEqual(t, "id", actual[0].ID, id)
	})

	t.Run("should return multiple matching values", func(t *testing.T) {
		matchA := string_index.Match{12, 0}
		matchB := string_index.Match{78, 0}
		query := "multiple_indexed_values"

		idx := string_index.NewIndex()
		idx.Index(matchA.ID, int(matchA.Confidence), query)
		idx.Index(matchB.ID, int(matchB.Confidence), query)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, matchA.ID)
		assertEqual(t, "second id", actual[1].ID, matchB.ID)
	})

	t.Run("should only return matching values", func(t *testing.T) {
		match := string_index.Match{54, 0}
		not_match := string_index.Match{76, 0}
		query := "only_matching"

		idx := string_index.NewIndex()
		idx.Index(match.ID, int(match.Confidence), query)
		idx.Index(not_match.ID, int(not_match.Confidence), "%")
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 1)
		assertEqual(t, "id", actual[0].ID, match.ID)
	})

	t.Run("should return partially matching values", func(t *testing.T) {
		matchA := string_index.Match{98, 0}
		matchB := string_index.Match{81, 0}
		query := "abc xy"

		idx := string_index.NewIndex()
		idx.Index(matchA.ID, int(matchA.Confidence), "ab")
		idx.Index(matchB.ID, int(matchB.Confidence), "xyz")
		// TODO
		idx.Index(123456, 0, "%")
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, matchA.ID)
		assertEqual(t, "second id", actual[1].ID, matchB.ID)
	})

	t.Run("should return matched values in order of confidence", func(t *testing.T) {
		matchA := string_index.Match{21, 50}
		matchB := string_index.Match{82, 100}
		query := "matching_order_confidence"

		idx := string_index.NewIndex()
		idx.Index(matchA.ID, int(matchA.Confidence), query)
		idx.Index(matchB.ID, int(matchB.Confidence), query)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, matchB.ID)
		assertEqual(t, "second id", actual[1].ID, matchA.ID)
	})

	t.Run("should accumulate confidence from multiple index calls", func(t *testing.T) {
		matchA := string_index.Match{81, 100}
		matchB := string_index.Match{43, 60}
		query := "confidence_order_sum"

		idx := string_index.NewIndex()
		idx.Index(matchA.ID, int(matchA.Confidence), query)
		idx.Index(matchB.ID, int(matchB.Confidence), query)
		idx.Index(matchB.ID, int(matchB.Confidence), query)
		actual := idx.Search(query)

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, matchB.ID)
		assertEqual(t, "second id", actual[1].ID, matchA.ID)
	})

	t.Run("should accumulate confidence from multiple query parts", func(t *testing.T) {
		matchA := string_index.Match{43, 60}
		matchB := string_index.Match{91, 100}
		queries := []string{"query", "parts"}

		idx := string_index.NewIndex()
		idx.Index(matchA.ID, int(matchA.Confidence), queries...)
		idx.Index(matchB.ID, int(matchB.Confidence), queries[0])
		actual := idx.Search(strings.Join(queries, " "))

		assertEqual(t, "length", len(actual), 2)
		assertEqual(t, "first id", actual[0].ID, matchA.ID)
		assertEqual(t, "second id", actual[1].ID, matchB.ID)
	})
}
