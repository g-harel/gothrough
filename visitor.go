package gis

import (
	"go/ast"
	"go/token"

	"github.com/g-harel/gis/internal/interfaces"
)

var _ ast.Visitor = visitor{}

type visitor struct {
	FileSet          *token.FileSet
	RelativePath     string
	InterfaceHandler func(*interfaces.Interface)
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	i, ok := interfaces.FromDecl(n, v.RelativePath, v.FileSet)
	if ok {
		v.InterfaceHandler(i)
		return nil
	}

	return v
}
