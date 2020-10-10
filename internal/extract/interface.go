package extract

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/g-harel/gothrough/internal/types"
)

func newInterfaceVisitor(handler func(Location, types.Interface)) visitFunc {
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
							if interfaceSpec.Name.IsExported() {
								// Collect methods.
								methods := []types.Function{}
								embedded := []types.Reference{}
								for _, member := range interfaceType.Methods.List {
									if identType, ok := member.Type.(*ast.Ident); ok {
										embedded = append(embedded, types.Reference{
											Docs: types.Docs{Text: member.Doc.Text()},
											Name: identType.Name,
										})
										continue
									}
									if selectorExprType, ok := member.Type.(*ast.SelectorExpr); ok {
										embedded = append(embedded, types.Reference{
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
											methods = append(methods, types.Function{
												Name:         methodName.String(),
												Docs:         types.Docs{Text: member.Doc.Text()},
												Arguments:    arguments,
												ReturnValues: returnValues,
											})
										}
									}
								}

								handler(
									getLocation(filepath),
									types.Interface{
										Name:     name,
										Docs:     types.Docs{Text: typeDeclaration.Doc.Text()},
										Embedded: embedded,
										Methods:  methods,
									},
								)

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
