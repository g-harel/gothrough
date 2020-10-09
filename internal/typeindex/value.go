package typeindex

import (
	"github.com/g-harel/gothrough/internal/camel"
	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/types"
)

func (idx *Index) InsertValue(location extract.Location, val types.Value) {
	idx.results = append(idx.results, &Result{
		Name:     val.Name,
		Location: location,
		Value:    &val,
	})
	id := len(idx.results) - 1

	idx.insertLocation(id, location)

	// Index on name.
	idx.textIndex.Insert(id, confidenceHigh, val.Name)
	nameTokens := camel.Split(val.Name)
	if len(nameTokens) > 1 {
		idx.textIndex.Insert(id, confidenceHigh/len(nameTokens), nameTokens...)
	}
}
