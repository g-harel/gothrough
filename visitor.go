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
					if interfaceType, ok := interfaceSpec.Type.(*ast.InterfaceType); ok {
						name := interfaceSpec.Name.String()
						if unicode.IsUpper([]rune(name)[0]) {
							packageName := path.Base(v.RelativePath)

							// Render declaration as it appeared originally.
							buf := bytes.NewBufferString("")
							fset := token.NewFileSet()
							printer.Fprint(buf, fset, n)

							// Collect declaration position (for filename and line number).
							pos := v.FileSet.Position(typeDeclaration.Pos())

							// Collect method names.
							methods := []string{}
							for _, method := range interfaceType.Methods.List {
								for _, methodName := range method.Names {
									if methodName.IsExported() {
										methods = append(methods, methodName.String())
									}
								}
							}

							v.InterfaceHandler(Interface{
								Name:              name,
								Methods:           methods,
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
