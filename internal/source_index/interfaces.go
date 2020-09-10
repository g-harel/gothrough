package source_index

import (
	"strings"

	"github.com/g-harel/gothrough/internal/camel"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for interface info items.
const (
	interfaceNameVal           = 120
	totalInterfaceNameTokenVal = 120
	packageNameVal             = 120
	sourceFileVal              = 10
	totalImportPathPartVal     = 20
	totalEmbeddedNameVal       = 80
	totalEmbeddedNameTokenVal  = 80
	totalMethodNameVal         = 80
	totalMethodNameTokenVal    = 80
)

func (si *Index) InsertInterface(ifc types.Interface) {
	si.interfaces = append(si.interfaces, &ifc)
	// TODO id generator that would work with multiple types.
	id := len(si.interfaces) - 1

	// Index on interface name.
	si.textIndex.Insert(id, interfaceNameVal, ifc.Name)
	nameTokens := camel.Split(ifc.Name)
	if len(nameTokens) > 1 {
		si.textIndex.Insert(id, totalInterfaceNameTokenVal/len(nameTokens), nameTokens...)
	}

	// Index on package path and source file.
	importPathParts := strings.Split(ifc.PackageImportPath, "/")
	si.textIndex.Insert(id, packageNameVal, ifc.PackageName)
	si.textIndex.Insert(id, sourceFileVal, strings.TrimSuffix(ifc.SourceFile, ".go"))
	if len(importPathParts) > 1 {
		si.textIndex.Insert(id, totalImportPathPartVal/len(importPathParts), importPathParts...)
	}

	// Index on embedded interfaces.
	if len(ifc.Embedded) > 0 {
		for _, embedded := range ifc.Embedded {
			si.textIndex.Insert(id, totalEmbeddedNameVal/len(ifc.Embedded), embedded.Name)
		}
		embeddedNameTokens := []string{}
		for _, embedded := range ifc.Embedded {
			if embedded.Package != "" {
				embeddedNameTokens = append(embeddedNameTokens, embedded.Package)
			}
			embeddedNameTokens = append(embeddedNameTokens, camel.Split(embedded.Name)...)
		}
		if len(embeddedNameTokens) > 1 {
			si.textIndex.Insert(id, totalEmbeddedNameTokenVal/len(embeddedNameTokens), embeddedNameTokens...)
		}
	}

	// Index on interface methods.
	if len(ifc.Methods) > 0 {
		for _, method := range ifc.Methods {
			si.textIndex.Insert(id, totalMethodNameVal/len(ifc.Methods), method.Name)
		}
		methodNameTokens := []string{}
		for _, method := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(method.Name)...)
		}
		if len(methodNameTokens) > 1 {
			si.textIndex.Insert(id, totalMethodNameTokenVal/len(methodNameTokens), methodNameTokens...)
		}
	}
}
