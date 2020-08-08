package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// visitFunc is called for every node encountered while walking a file.
type visitFunc func(filepath string, n ast.Node, fset *token.FileSet) (done bool)

var _ ast.Visitor = visitor{}

type visitor struct {
	visitFunc visitFunc
	filepath  string
	fset      *token.FileSet
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	done := v.visitFunc(v.filepath, n, v.fset)
	if done {
		return nil
	}
	return v
}

// Visit walks the visitor on the provided file.
func visit(filepath string, v visitFunc) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %v", err)
	}

	ast.Walk(
		visitor{v, filepath, fset},
		file,
	)

	return nil
}
