package gis

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"path"
	"unicode"
)

var _ ast.Visitor = visitor{}

type visitor struct {
	FileSet          *token.FileSet
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
							buf := bytes.NewBufferString("")
							fset := token.NewFileSet()
							printer.Fprint(buf, fset, n)
							pos := v.FileSet.Position(typeDeclaration.Pos())
							v.InterfaceHandler(Interface{
								Name:              name,
								Body:              buf.String(),
								PackageName:       packageName,
								PackageImportPath: v.RelativePath,
								SourceFile:        path.Base(pos.Filename),
								SourceLine:        pos.Line,
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
