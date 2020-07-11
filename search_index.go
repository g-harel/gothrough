package gis

import (
	"github.com/g-harel/gis/internal/index"
	"github.com/g-harel/gis/internal/interfaces"
)

type SearchIndex struct {
	index      *index.Index
	interfaces []*interfaces.Interface
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		index:      index.NewIndex(),
		interfaces: []*interfaces.Interface{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *SearchIndex) Search(query string) ([]*interfaces.Interface, error) {
	searchResult := si.index.Search(query)
	res := make([]*interfaces.Interface, len(searchResult))
	for i, pos := range searchResult {
		res[i] = si.interfaces[pos]
	}

	return res, nil
}
