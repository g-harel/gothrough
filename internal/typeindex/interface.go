package typeindex

import (
	"github.com/g-harel/gothrough/internal/cases"
	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/types"
)

func (idx *Index) InsertInterface(location extract.Location, ifc types.Interface) {
	idx.results = append(idx.results, &Result{
		Name:     ifc.Name,
		Location: location,
		Value:    &ifc,
	})
	id := len(idx.results) - 1

	idx.insertLocation(id, location)

	// Index on name.
	idx.textIndex.Insert(id, confidenceHigh, ifc.Name)
	nameTokens := cases.Split(ifc.Name)
	if len(nameTokens) > 1 {
		idx.textIndex.Insert(id, confidenceHigh/len(nameTokens), nameTokens...)
	}

	// Index on embedded interfaces.
	if len(ifc.Embedded) > 0 {
		for _, embedded := range ifc.Embedded {
			idx.textIndex.Insert(id, confidenceMed/len(ifc.Embedded), embedded.Name)
		}
		embeddedNameTokens := []string{}
		for _, embedded := range ifc.Embedded {
			if embedded.Package != "" {
				embeddedNameTokens = append(embeddedNameTokens, embedded.Package)
			}
			embeddedNameTokens = append(embeddedNameTokens, cases.Split(embedded.Name)...)
		}
		if len(embeddedNameTokens) > 1 {
			idx.textIndex.Insert(id, confidenceMed/len(embeddedNameTokens), embeddedNameTokens...)
		}
	}

	// Index on methods.
	if len(ifc.Methods) > 0 {
		for _, method := range ifc.Methods {
			idx.textIndex.Insert(id, confidenceMed/len(ifc.Methods), method.Name)
		}
		methodNameTokens := []string{}
		for _, method := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, cases.Split(method.Name)...)
		}
		if len(methodNameTokens) > 1 {
			idx.textIndex.Insert(id, confidenceMed/len(methodNameTokens), methodNameTokens...)
		}
	}
}
