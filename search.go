package gis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/gis/internal/camel"
	"github.com/g-harel/gis/internal/index"
	"github.com/g-harel/gis/internal/interfaces"
	"github.com/g-harel/gis/internal/parse"
)

// Confidence values for interface info items.
const (
	interfaceNameVal           = 120
	totalInterfaceNameTokenVal = 160
	packageNameVal             = 120
	sourceFileVal              = 10
	totalImportPathPartVal     = 20
	totalMethodNameVal         = 80
	totalMethodNameTokenVal    = 80
)

type SearchIndex struct {
	index      *index.Index
	interfaces []*interfaces.Interface
}

// NewSearchIndex creates a searchable index of interfaces in the provided root directory.
func NewSearchIndex(rootDir string) (*SearchIndex, error) {
	// Collect all interfaces in the provided directory.
	si := &SearchIndex{interfaces: []*interfaces.Interface{}}
	err := filepath.Walk(rootDir, func(pathname string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(pathname, ".go") {
			return nil
		}
		if strings.HasSuffix(pathname, "_test.go") {
			return nil
		}
		if strings.Contains(pathname, "internal/") {
			return nil
		}
		if strings.Contains(pathname, "vendor/") {
			return nil
		}
		if strings.Contains(pathname, "testdata/") {
			return nil
		}
		if strings.Contains(pathname, "testing/") {
			return nil
		}
		parse.Visit(pathname, parse.NewInterfaceVisitor(func(ifc interfaces.Interface) {
			si.interfaces = append(si.interfaces, &ifc)
		}))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %v", err)
	}

	// Add the interfaces to the index with a confidence value.
	idx := index.NewIndex()
	for i, ifc := range si.interfaces {
		// Index on interface name.
		idx.Index(i, interfaceNameVal, ifc.Name)
		nameTokens := camel.Split(ifc.Name)
		if len(nameTokens) > 1 {
			idx.Index(i, totalInterfaceNameTokenVal/len(nameTokens), nameTokens...)
		}

		// Index on package path and source file.
		importPathParts := strings.Split(ifc.PackageImportPath, "/")
		idx.Index(i, packageNameVal, ifc.PackageName)
		idx.Index(i, sourceFileVal, strings.TrimSuffix(ifc.SourceFile, ".go"))
		if len(importPathParts) > 1 {
			idx.Index(i, totalImportPathPartVal/len(importPathParts), importPathParts...)
		}

		// Index on interface methods.
		methodNameTokens := []string{}
		for _, methodName := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(methodName)...)
		}
		if len(ifc.Methods) > 0 {
			idx.Index(i, totalMethodNameVal/len(ifc.Methods), ifc.Methods...)
		}
		if len(methodNameTokens) > 0 {
			idx.Index(i, totalInterfaceNameTokenVal/len(methodNameTokens), methodNameTokens...)
		}
	}

	si.index = idx
	return si, nil
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
