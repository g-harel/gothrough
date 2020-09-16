package source_index

import (
	"strings"

	"github.com/g-harel/gothrough/internal/camel"
	"github.com/g-harel/gothrough/internal/types"
)

func (si *Index) InsertInterface(ifc types.Interface) {
	si.results = append(si.results, &Result{
		Name:              ifc.Name,
		PackageName:       ifc.PackageName,
		PackageImportPath: ifc.PackageImportPath,
		Value:             &ifc,
	})
	id := len(si.results) - 1

	// Index on interface name.
	si.textIndex.Insert(id, confidenceHigh, ifc.Name)
	nameTokens := camel.Split(ifc.Name)
	if len(nameTokens) > 1 {
		si.textIndex.Insert(id, confidenceHigh/len(nameTokens), nameTokens...)
	}

	// Index on package path and source file.
	importPathParts := strings.Split(ifc.PackageImportPath, "/")
	si.textIndex.Insert(id, confidenceHigh, ifc.PackageName)
	si.textIndex.Insert(id, confidenceLow, strings.TrimSuffix(ifc.SourceFile, ".go"))
	if len(importPathParts) > 1 {
		si.textIndex.Insert(id, confidenceLow/len(importPathParts), importPathParts...)
	}

	// Index on embedded interfaces.
	if len(ifc.Embedded) > 0 {
		for _, embedded := range ifc.Embedded {
			si.textIndex.Insert(id, confidenceMed/len(ifc.Embedded), embedded.Name)
		}
		embeddedNameTokens := []string{}
		for _, embedded := range ifc.Embedded {
			if embedded.Package != "" {
				embeddedNameTokens = append(embeddedNameTokens, embedded.Package)
			}
			embeddedNameTokens = append(embeddedNameTokens, camel.Split(embedded.Name)...)
		}
		if len(embeddedNameTokens) > 1 {
			si.textIndex.Insert(id, confidenceMed/len(embeddedNameTokens), embeddedNameTokens...)
		}
	}

	// Index on interface methods.
	if len(ifc.Methods) > 0 {
		for _, method := range ifc.Methods {
			si.textIndex.Insert(id, confidenceMed/len(ifc.Methods), method.Name)
		}
		methodNameTokens := []string{}
		for _, method := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(method.Name)...)
		}
		if len(methodNameTokens) > 1 {
			si.textIndex.Insert(id, confidenceMed/len(methodNameTokens), methodNameTokens...)
		}
	}
}
