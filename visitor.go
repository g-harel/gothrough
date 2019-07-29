package gis

import (
	"go/ast"
	"go/token"
	"path"
	"unicode"
)

var _ ast.Visitor = visitor{}

type visitor struct {
	RelativePath     string
	InterfaceHandler func(Interface)
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if typeDeclaration, ok := n.(*ast.GenDecl); ok {
		if typeDeclaration.Tok == token.TYPE {
			for _, spec := range typeDeclaration.Specs {
				if interfaceSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := interfaceSpec.Type.(*ast.InterfaceType); ok {
						name := interfaceSpec.Name.String()
						packageName := path.Base(v.RelativePath)
						if unicode.IsUpper([]rune(name)[0]) {
							v.InterfaceHandler(Interface{
								Name:              name,
								Description:       typeDeclaration.Doc.Text(),
								PackageName:       packageName,
								PackageImportPath: v.RelativePath,
							})
							return nil
						}
					}
				}
			}
		}
	}

	return v
}
