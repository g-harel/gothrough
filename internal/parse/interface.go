package parse

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"path"
	"strings"
	"unicode"

	"github.com/g-harel/gis/internal/interfaces"
)

// NewInterfaceVisitor creates a visitor that collects interfaces into the target array.
func NewInterfaceVisitor(handler func(interfaces.Interface)) Visitor {
	return func(filepath string, n ast.Node, fset *token.FileSet) bool {
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
								pathname, filename := path.Split(filepath)

								// Attempt to detect source dir by looking for the closest "src" directory.
								pathParts := strings.Split(path.Clean(pathname), "/src")
								srcDir := path.Join(pathParts[:len(pathParts)-1]...) + "/src/"

								relativePath := strings.TrimPrefix(path.Dir(filepath), srcDir)

								// Render declaration as it appeared originally.
								buf := bytes.NewBufferString("")
								renderer := token.NewFileSet()
								format.Node(buf, renderer, n)

								// Collect method names.
								methods := []string{}
								for _, method := range interfaceType.Methods.List {
									for _, methodName := range method.Names {
										if methodName.IsExported() {
											methods = append(methods, methodName.String())
										}
									}
								}

								handler(interfaces.Interface{
									Name:              name,
									Methods:           methods,
									Printed:           buf.String(),
									PackageName:       path.Base(relativePath),
									PackageImportPath: relativePath,
									SourceFile:        filename,
									SourceLine:        fset.Position(typeDeclaration.Pos()).Line,
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
