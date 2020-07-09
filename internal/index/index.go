package index

import (
	"math"
	"sort"
	"strings"
)

// Penalty applied to confidence value when indexing substrings.
// Refer to implementation for usage details.
const substringPenalty = 4.0

type mappedValue struct {
	ID         int
	Confidence float64
}

// Index is a simple search index.
// It does not store values, only the value's ID provided by the consumer of this package.
// The IDs of all matching values are returned as the result of a search.
type Index struct {
	Mappings map[string][]mappedValue
}

// NewIndex creates a new index.
func NewIndex() *Index {
	return &Index{
		Mappings: map[string][]mappedValue{},
	}
}

// Index stores the relationship with the given ID and strings (queries).
// It also requires a confidence to be associated with the relationship, which is used to order results.
func (idx *Index) Index(id int, confidence int, strs ...string) {
	for _, str := range strs {
		str = strings.ToLower(str)
		for i := 1; i <= len(str); i++ {
			substringFraction := float64(i) / float64(len(str))
			adjustedConfidence := float64(confidence) * math.Pow(substringFraction, substringPenalty)
			for _, substr := range Substrings(str, i) {
				if len(idx.Mappings[substr]) == 0 {
					idx.Mappings[substr] = []mappedValue{}
				}
				idx.Mappings[substr] = append(idx.Mappings[substr], mappedValue{id, adjustedConfidence})
			}
		}
	}
}

// Search searches for indexed values matching the query.
// Results are ordered by descending order of confidence.
func (idx *Index) Search(query string) []int {
	query = strings.ToLower(query)

	// Sum confidences from matching mappings.
	confidences := map[int]float64{}
	for _, subQuery := range strings.Fields(query) {
		for i := 1; i <= len(subQuery); i++ {
			for _, substr := range Substrings(subQuery, i) {
				for _, m := range idx.Mappings[substr] {
					if _, ok := confidences[m.ID]; !ok {
						confidences[m.ID] = 0
					}
					confidences[m.ID] += m.Confidence
				}
			}
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
