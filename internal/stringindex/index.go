package stringindex

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

// Insert stores the relationship with the given ID and strings (queries).
// It also requires a confidence to be associated with the relationship, which is used to order results.
func (idx *Index) Insert(id int, confidence int, strs ...string) {
	for _, str := range strs {
		str = strings.ToLower(str)
		for i := 1; i <= len(str); i++ {
			adjustmentFactor := div(math.Pow(float64(i), float64(len(str))), 3)
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
	// confidences[sub_query_index][result_id] -> confidence
	perSubQueryConfidences := map[int]map[int]float64{}
	maxPerSubQueryConfidence := 0.0
	for subQueryID, subQuery := range strings.Fields(query) {
		perSubQueryConfidences[subQueryID] = map[int]float64{}
		for j := 1; j <= len(subQuery); j++ {
			adjustmentFactor := div(float64(j), math.Pow(float64(len(subQuery)), substringPenaltyScale))
			for _, substr := range Substrings(subQuery, j) {
				for _, m := range idx.Matches[substr] {
					if _, ok := perSubQueryConfidences[subQueryID][m.ID]; !ok {
						perSubQueryConfidences[subQueryID][m.ID] = 0
					}
					perSubQueryConfidences[subQueryID][m.ID] += m.Confidence * adjustmentFactor
					maxPerSubQueryConfidence = math.Max(
						maxPerSubQueryConfidence,
						perSubQueryConfidences[subQueryID][m.ID],
					)
				}
			}
		}
	}

	// Combine together confidences of multiple sub-queries.
	combinedConfidences := map[int]float64{}
	for subQueryID := range perSubQueryConfidences {
		for id := range perSubQueryConfidences[subQueryID] {
			// Combine only when first seen.
			if _, ok := combinedConfidences[id]; ok {
				continue
			}

			// Collect confidence values of each sub-query, or 0 if no match.
			adjustedConfidenceValues := []float64{}
			for i := range perSubQueryConfidences {
				confidence, ok := perSubQueryConfidences[i][id]
				if ok {
					adjustedConfidence := div(confidence, maxPerSubQueryConfidence)
					adjustedConfidenceValues = append(adjustedConfidenceValues, adjustedConfidence)
				} else {
					adjustedConfidenceValues = append(adjustedConfidenceValues, 0.0)
				}
			}

			// Calculate confidence total and variance.
			adjustedConfidenceCount := float64(len(adjustedConfidenceValues))
			adjustedConfidenceTotal := 0.0
			for _, confidence := range adjustedConfidenceValues {
				adjustedConfidenceTotal += confidence
			}
			adjustedConfidenceMean := div(adjustedConfidenceTotal, adjustedConfidenceCount)
			adjustedConfidenceVariance := 0.0
			for _, confidence := range adjustedConfidenceValues {
				adjustedConfidenceVariance += div(
					math.Pow(confidence-adjustedConfidenceMean, 2),
					adjustedConfidenceCount,
				)
			}

			// Combine confidences, rewarding results that have both a large sum
			// and low standard deviation. Adjusted confidence values are <= 1.
			combinedConfidences[id] = adjustedConfidenceTotal * (1 - math.Sqrt(adjustedConfidenceVariance))
		}
	}

	// Flatten confidence map and scale values to [0, 1] range.
	results := []Match{}
	maxConfidence := 0.0
	for id := range combinedConfidences {
		maxConfidence = math.Max(maxConfidence, combinedConfidences[id])
		results = append(results, Match{
			ID:         id,
			Confidence: combinedConfidences[id],
		})
	}
	for i := range results {
		results[i].Confidence = div(results[i].Confidence, maxConfidence)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})

	return results
}

func div(a, b float64) float64 {
	if a == 0.0 {
		return 0.0
	}
	return a / b
}
