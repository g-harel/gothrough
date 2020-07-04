package index

import (
	"sort"
	"strings"
)

type mappedValue struct {
	id         int
	confidence float32
}

// Index is a simple search index.
// The indexed values are never stored, only the ID provided by the consumer of this package.
// The IDs of all matching values are returned as the result of a search.
type Index struct {
	mappings map[string][]mappedValue
}

// NewIndex creates a new index.
func NewIndex() *Index {
	return &Index{
		mappings: map[string][]mappedValue{},
	}
}

// Index stores the relationship with the given ID and strings (queries).
// It also requires a confidence to be associated with the relationship, which is used to order results.
func (idx *Index) Index(id int, confidence int, strs ...string) {
	for _, str := range strs {
		str = strings.ToLower(str)
		for i := 1; i <= len(str); i++ {
			adjustedConfidence := float32(confidence) * (float32(i) / float32(len(str)))
			// TODO test
			for _, substr := range Substrings(str, i) {
				if len(idx.mappings[substr]) == 0 {
					idx.mappings[substr] = []mappedValue{}
				}
				idx.mappings[substr] = append(idx.mappings[substr], mappedValue{id, adjustedConfidence})
			}
		}
	}
}

// Search searches for indexed values matching the query.
// Results are ordered by descending order of confidence.
func (idx *Index) Search(query string) []int {
	query = strings.ToLower(query)

	// Sum confidences from matching mappings.
	confidences := map[int]float32{}
	for _, subQuery := range strings.Fields(query) {
		// TODO query substrings too (adjust for sub query length).
		for _, m := range idx.mappings[subQuery] {
			if _, ok := confidences[m.id]; !ok {
				confidences[m.id] = 0
			}
			confidences[m.id] += m.confidence
		}
	}

	// Collect matched IDs (no duplicates).
	ids := make([]int, len(confidences))
	i := 0
	for id := range confidences {
		ids[i] = id
		i++
	}

	sort.Slice(ids, func(i, j int) bool {
		return confidences[ids[i]] > confidences[ids[j]]
	})

	return ids
}
