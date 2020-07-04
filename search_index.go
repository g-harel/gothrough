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
		nameTokens := camel.Split(ifc.Name)
		methodNameTokens := []string{}
		for _, methodName := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(methodName)...)
		}

		idx.Index(i, 100, ifc.Name)
		idx.Index(i, 70, ifc.PackageName)
		idx.Index(i, 20, strings.Split(ifc.PackageImportPath, "/")...)
		idx.Index(i, 20, strings.TrimSuffix(ifc.SourceFile, ".go"))
		if len(nameTokens) > 0 {
			idx.Index(i, 60/len(nameTokens), nameTokens...)
		}
		if len(ifc.Methods) > 0 {
			idx.Index(i, 50/len(ifc.Methods), ifc.Methods...)
		}
		if len(methodNameTokens) > 0 {
			idx.Index(i, 40/len(methodNameTokens), methodNameTokens...)
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
