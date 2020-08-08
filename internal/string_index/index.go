package string_index

import (
	"math"
	"sort"
	"strings"
)

const (
	substringPenaltyScale = 3
)

type Match struct {
	ID         int
	Confidence float64
}

// Index is a simple search index.
// It does not store values, only the value's ID provided by the consumer of this package.
// The IDs of all matching values are returned as the result of a search.
type Index struct {
	Matches map[string][]Match
}

// NewIndex creates a new index.
func NewIndex() *Index {
	return &Index{
		Matches: map[string][]Match{},
	}
}

// Index stores the relationship with the given ID and strings (queries).
// It also requires a confidence to be associated with the relationship, which is used to order results.
func (idx *Index) Index(id int, confidence int, strs ...string) {
	for _, str := range strs {
		str = strings.ToLower(str)
		for i := 1; i <= len(str); i++ {
			adjustmentFactor := math.Pow(float64(i)/float64(len(str)), 3)
			adjustedConfidence := float64(confidence) * adjustmentFactor
			for _, substr := range Substrings(str, i) {
				if len(idx.Matches[substr]) == 0 {
					idx.Matches[substr] = []Match{}
				}
				idx.Matches[substr] = append(idx.Matches[substr], Match{id, adjustedConfidence})
			}
		}
	}
}

// Search searches for indexed values matching the query.
// Results are ordered by descending order of confidence.
func (idx *Index) Search(query string) []Match {
	query = strings.ToLower(query)

	// Sum confidences from matches.
	confidences := map[int]float64{}
	for _, subQuery := range strings.Fields(query) {
		for i := 1; i <= len(subQuery); i++ {
			adjustmentFactor := float64(i) / math.Pow(float64(len(subQuery)), substringPenaltyScale)
			for _, substr := range Substrings(subQuery, i) {
				for _, m := range idx.Matches[substr] {
					if _, ok := confidences[m.ID]; !ok {
						confidences[m.ID] = 0
					}
					confidences[m.ID] += m.Confidence * adjustmentFactor
				}
			}
		}
	}

	// Collect matched IDs (no duplicates).
	results := []Match{}
	for id := range confidences {
		results = append(results, Match{
			ID:         id,
			Confidence: confidences[id],
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})

	return results
}
