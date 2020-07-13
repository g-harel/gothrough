package interface_index

import (
	"github.com/g-harel/gis/internal/interfaces"
	"github.com/g-harel/gis/internal/string_index"
)

type Index struct {
	index      *string_index.Index
	interfaces []*interfaces.Interface
}

func NewIndex() *Index {
	return &Index{
		index:      string_index.NewIndex(),
		interfaces: []*interfaces.Interface{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *Index) Search(query string) ([]*interfaces.Interface, error) {
	searchResult := si.index.Search(query)
	results := make([]*interfaces.Interface, len(searchResult))
	for i, result := range searchResult {
		results[i] = si.interfaces[result.ID]
	}

	return results, nil
}
