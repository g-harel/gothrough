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
func extract(dir string, in []string) ([]Interface, error) {
	fs := token.NewFileSet()
	out := []Interface{}
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
				InterfaceHandler: func(i Interface) {
					out = append(out, i)
				},
			},
			f,
		)
	}
	return out, nil
}

// Search finds interfaces in the given directory.
func Search(dir string) ([]Interface, error) {
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

	return extractOut, nil
}
