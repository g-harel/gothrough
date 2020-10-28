package typeindex

import (
	"github.com/g-harel/gothrough/internal/cases"
	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/types"
)

func (idx *Index) InsertFunction(location extract.Location, fnc types.Function) {
	idx.results = append(idx.results, &Result{
		Name:     fnc.Name,
		Location: location,
		Value:    &fnc,
	})
	id := len(idx.results) - 1

	idx.insertLocation(id, location)

	// Index on name.
	idx.textIndex.Insert(id, confidenceHigh, fnc.Name)
	nameTokens := cases.Split(fnc.Name)
	if len(nameTokens) > 1 {
		idx.textIndex.Insert(id, confidenceHigh/len(nameTokens), nameTokens...)
	}

	// Index on arguments.
	if len(fnc.Arguments) > 0 {
		for _, argument := range fnc.Arguments {
			if argument.Name == "" {
				continue
			}
			idx.textIndex.Insert(id, confidenceMed/len(fnc.Arguments), argument.Name)
			argumentNameTokens := cases.Split(argument.Name)
			if len(argumentNameTokens) > 1 {
				idx.textIndex.Insert(id, confidenceMed/len(argumentNameTokens), argumentNameTokens...)
			}
		}
	}

	// Index on return values.
	if len(fnc.ReturnValues) > 0 {
		for _, returnValue := range fnc.ReturnValues {
			if returnValue.Name == "" {
				continue
			}
			idx.textIndex.Insert(id, confidenceMed/len(fnc.ReturnValues), returnValue.Name)
			returnValueNameTokens := cases.Split(returnValue.Name)
			if len(returnValueNameTokens) > 1 {
				idx.textIndex.Insert(id, confidenceMed/len(returnValueNameTokens), returnValueNameTokens...)
			}
		}
	}
}
