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

// TODO serialize to file and create from it.
type SearchIndex struct {
	index      *index.Index
	interfaces []*interfaces.Interface
}

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
		idx.Index(i, 120, ifc.Name)
		nameTokens := camel.Split(ifc.Name)
		if len(nameTokens) > 1 {
			idx.Index(i, 160/len(nameTokens), nameTokens...)
		}

		// Index on package path and source file.
		importPathParts := strings.Split(ifc.PackageImportPath, "/")
		idx.Index(i, 120, ifc.PackageName)
		idx.Index(i, 10, strings.TrimSuffix(ifc.SourceFile, ".go"))
		if len(importPathParts) > 1 {
			idx.Index(i, 20/len(importPathParts), importPathParts...)
		}

		// Index on interface methods.
		methodNameTokens := []string{}
		for _, methodName := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(methodName)...)
		}
		if len(ifc.Methods) > 0 {
			idx.Index(i, 80/len(ifc.Methods), ifc.Methods...)
		}
		if len(methodNameTokens) > 0 {
			idx.Index(i, 80/len(methodNameTokens), methodNameTokens...)
		}
	}

	si.index = idx
	return si, nil
}

func (si *SearchIndex) Search(query string) ([]*interfaces.Interface, error) {
	searchResult := si.index.Search(query)
	res := make([]*interfaces.Interface, len(searchResult))
	for i, pos := range searchResult {
		res[i] = si.interfaces[pos]
	}

	return res, nil
}
