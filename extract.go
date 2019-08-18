package gis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/g-harel/gis/internal/interfaces"
)

var _ ast.Visitor = visitor{}

type visitor struct {
	FileSet          *token.FileSet
	RelativePath     string
	InterfaceHandler func(*interfaces.Interface)
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	i, ok := interfaces.FromNode(n, v.RelativePath, v.FileSet)
	if ok {
		v.InterfaceHandler(i)
		return nil
	}

	return v
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
