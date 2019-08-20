package parse

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"path"
	"unicode"

	"github.com/g-harel/gis/internal/interfaces"
)

func Interface(relativePath string, target []*interfaces.Interface) Visitor {
	return func(n ast.Node, fset *token.FileSet) bool {
		if n == nil {
			return true
		}

		if typeDeclaration, ok := n.(*ast.GenDecl); ok {
			if typeDeclaration.Tok == token.TYPE {
				for _, spec := range typeDeclaration.Specs {
					if interfaceSpec, ok := spec.(*ast.TypeSpec); ok {
						if interfaceType, ok := interfaceSpec.Type.(*ast.InterfaceType); ok {
							name := interfaceSpec.Name.String()
							if unicode.IsUpper([]rune(name)[0]) {
								packageName := path.Base(relativePath)

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

								target = append(target, &interfaces.Interface{
									Name:              name,
									Methods:           methods,
									Body:              buf.String(),
									PackageName:       packageName,
									PackageImportPath: relativePath,
									SourceFile:        path.Base(pos.Filename),
									SourceLine:        pos.Line,
								})
								return true
							}
						}
					}
				}
			}
		}
		return false
	}
}
