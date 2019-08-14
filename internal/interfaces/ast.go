package interfaces

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"path"
	"unicode"
)

func FromNode(n ast.Node, relPath string, fset *token.FileSet) (*Interface, bool) {
	if n == nil {
		return nil, false
	}

	if typeDeclaration, ok := n.(*ast.GenDecl); ok {
		if typeDeclaration.Tok == token.TYPE {
			for _, spec := range typeDeclaration.Specs {
				if interfaceSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := interfaceSpec.Type.(*ast.InterfaceType); ok {
						name := interfaceSpec.Name.String()
						if unicode.IsUpper([]rune(name)[0]) {
							packageName := path.Base(relPath)

							// Render declaration as it appeared originally.
							buf := bytes.NewBufferString("")
							renderer := token.NewFileSet()
							printer.Fprint(buf, renderer, n)

							// Collect declaration position (for filename and line number).
							pos := fset.Position(typeDeclaration.Pos())

							// Collect method names.
							methods := []string{}
							for _, method := range interfaceType.Methods.List {
								for _, methodName := range method.Names {
									if methodName.IsExported() {
										methods = append(methods, methodName.String())
									}
								}
							}

							return &Interface{
								Name:              name,
								Methods:           methods,
								Body:              buf.String(),
								PackageName:       packageName,
								PackageImportPath: relPath,
								SourceFile:        path.Base(pos.Filename),
								SourceLine:        pos.Line,
							}, true
						}
					}
				}
			}
		}
	}

	return nil, false
}
