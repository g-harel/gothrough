package extract

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"
	"unicode"

	"github.com/g-harel/gothrough/internal/types"
)

// newInterfaceVisitor creates a visitor that collects interfaces into the target array.
func newInterfaceVisitor(handler func(types.Interface)) visitFunc {
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
							// Only parse public interfaces.
							if unicode.IsUpper([]rune(name)[0]) {
								pathname, filename := path.Split(filepath)

								// Attempt to detect source dir by looking for the closest "src" directory.
								pathParts := strings.Split(path.Clean(pathname), "/src")
								srcDir := path.Join(pathParts[:len(pathParts)-1]...) + "/src/"

								relativePath := strings.TrimPrefix(path.Dir(filepath), srcDir)

								// Collect methods.
								methods := []types.MethodSignature{}
								embedded := []types.EmbeddedInterface{}
								for _, member := range interfaceType.Methods.List {
									if identType, ok := member.Type.(*ast.Ident); ok {
										embedded = append(embedded, types.EmbeddedInterface{
											Docs: types.Docs{Text: member.Doc.Text()},
											Name: identType.Name,
										})
										continue
									}
									if selectorExprType, ok := member.Type.(*ast.SelectorExpr); ok {
										embedded = append(embedded, types.EmbeddedInterface{
											Package: fmt.Sprintf("%v", selectorExprType.X),
											Name:    selectorExprType.Sel.String(),
											Docs:    types.Docs{Text: member.Doc.Text()},
										})
										continue
									}

									arguments := []types.Field{}
									returnValues := []types.Field{}
									if methodType, ok := member.Type.(*ast.FuncType); ok {
										arguments = collectFields(methodType.Params)
										returnValues = collectFields(methodType.Results)
									}
									for _, methodName := range member.Names {
										if methodName.IsExported() {
											methods = append(methods, types.MethodSignature{
												Name:         methodName.String(),
												Docs:         types.Docs{Text: member.Doc.Text()},
												Arguments:    arguments,
												ReturnValues: returnValues,
											})
										}
									}
								}

								handler(types.Interface{
									Name:              name,
									Docs:              types.Docs{Text: typeDeclaration.Doc.Text()},
									Embedded:          embedded,
									Methods:           methods,
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
