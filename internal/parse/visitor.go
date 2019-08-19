package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Visitor func(n ast.Node, fset *token.FileSet) (done bool)

var _ ast.Visitor = visitor{}

type visitor struct {
	visitFunc Visitor
	fset      *token.FileSet
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	done := v.visitFunc(n, v.fset)
	if done {
		return nil
	} else {
		return v
	}
}

func Visit(filepath string, v Visitor) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("parse file: %v", err)
	}

	ast.Walk(
		visitor{v, fset},
		file,
	)

	return nil
}
