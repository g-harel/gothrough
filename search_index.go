package gis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/gis/internal/camel"
	"github.com/g-harel/gis/internal/index"
	"github.com/g-harel/gis/internal/interfaces"
)

type SearchIndex struct {
	index      *index.Index
	interfaces []*interfaces.Interface
}

func NewSearchIndex(dir string) (*SearchIndex, error) {
	indexedFiles, err := find(dir)
	if err != nil {
		return nil, fmt.Errorf("find files: %v", err)
	}

	extractedInterfaces, err := extract(dir, indexedFiles)
	if err != nil {
		return nil, fmt.Errorf("extract interfaces: %v", err)
	}

	idx := index.NewIndex()
	for i, ifc := range extractedInterfaces {
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

	return &SearchIndex{
		index:      idx,
		interfaces: extractedInterfaces,
	}, nil
}

func (s *SearchIndex) Search(query string) ([]*interfaces.Interface, error) {
	searchResult := s.index.Search(query)
	res := make([]*interfaces.Interface, len(searchResult))
	for i, pos := range searchResult {
		res[i] = s.interfaces[pos]
	}

	return res, nil
}

// Find writes all files in the given directory to the output.
func find(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(pathname string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasSuffix(pathname, ".go") &&
				!strings.HasSuffix(pathname, "_test.go") &&
				!strings.Contains(pathname, "internal/") &&
				!strings.Contains(pathname, "vendor/") &&
				!strings.Contains(pathname, "testdata/") &&
				!strings.Contains(pathname, "testing/") {
				files = append(files, pathname)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %v", err)
	}
	return files, err
}
