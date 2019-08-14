package gis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/g-harel/gis/internal/camel"
	"github.com/g-harel/gis/internal/index"
	"github.com/g-harel/gis/internal/interfaces"
)

// Find writes all files in the given directory to the output.
func find(dir string) ([]string, error) {
	out := []string{}
	err := filepath.Walk(dir, func(pathname string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			out = append(out, pathname)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %v", err)
	}
	return out, err
}

// Filter only forwards paths to the output if they match certain conditions.
func filter(in []string) ([]string, error) {
	out := []string{}
	for _, pathname := range in {
		if !strings.HasSuffix(pathname, ".go") {
			continue
		}
		if strings.HasSuffix(pathname, "_test.go") {
			continue
		}
		if strings.Contains(pathname, "internal/") {
			continue
		}
		if strings.Contains(pathname, "vendor/") {
			continue
		}
		if strings.Contains(pathname, "testdata/") {
			continue
		}
		if strings.Contains(pathname, "testing/") {
			continue
		}
		out = append(out, pathname)
	}
	return out, nil
}

// Extract parses the file and walks the AST to extract interfaces.
func extract(dir string, in []string) ([]*interfaces.Interface, error) {
	fs := token.NewFileSet()
	out := []*interfaces.Interface{}
	for _, pathname := range in {
		f, err := parser.ParseFile(fs, pathname, nil, parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("parse file: %v", err)
		}

		relativePath := strings.TrimPrefix(path.Dir(pathname), dir)
		relativePath = strings.TrimPrefix(relativePath, "/")
		ast.Walk(
			visitor{
				FileSet:      fs,
				RelativePath: relativePath,
				InterfaceHandler: func(i *interfaces.Interface) {
					out = append(out, i)
				},
			},
			f,
		)
	}
	return out, nil
}

// Search finds interfaces in the given directory.
func Search(dir, query string) ([]*interfaces.Interface, error) {
	findOut, err := find(dir)
	if err != nil {
		return nil, fmt.Errorf("find files: %v", err)
	}

	filterOut, err := filter(findOut)
	if err != nil {
		return nil, fmt.Errorf("filter files: %v", err)
	}

	extractOut, err := extract(dir, filterOut)
	if err != nil {
		return nil, fmt.Errorf("extract interfaces: %v", err)
	}

	idx := index.NewIndex()
	for i, ifc := range extractOut {
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

	searchResult := idx.Search(query)
	res := make([]*interfaces.Interface, len(searchResult))
	for i, pos := range searchResult {
		res[i] = extractOut[pos]
	}

	return res, nil
}
