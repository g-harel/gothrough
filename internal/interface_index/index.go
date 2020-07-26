package interface_index

import (
	"github.com/g-harel/gis/internal/interfaces"
	"github.com/g-harel/gis/internal/string_index"
)

type Index struct {
	index      *string_index.Index
	interfaces []*interfaces.Interface
}

type Result struct {
	Interface  *interfaces.Interface
	Confidence float64
}

func NewIndex() *Index {
	return &Index{
		index:      string_index.NewIndex(),
		interfaces: []*interfaces.Interface{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *Index) Search(query string) ([]*Result, error) {
	searchResult := si.index.Search(query)
	if len(searchResult) == 0 {
		return []*Result{}, nil
	}

	maxConfidence := searchResult[0].Confidence
	results := make([]*Result, len(searchResult))
	for i, result := range searchResult {
		results[i] = &Result{
			Interface:  si.interfaces[result.ID],
			Confidence: result.Confidence / maxConfidence,
		}
	}

	return results, nil
}
