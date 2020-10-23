package extract

import (
	"go/ast"
	"go/token"

	"github.com/g-harel/gothrough/internal/types"
)

// TODO check behavior for iota
func newValueVisitor(handler func(Location, types.Value)) visitFunc {
	return func(filepath string, n ast.Node, fset *token.FileSet) bool {
		if n == nil {
			return true
		}

		if valueDeclaration, ok := n.(*ast.GenDecl); ok {
			if valueDeclaration.Tok == token.CONST || valueDeclaration.Tok == token.VAR {
				for _, spec := range valueDeclaration.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range valueSpec.Names {
							if name.IsExported() {
								prettyValue := ""
								if len(valueSpec.Values) > i {
									prettyValue = pretty(valueSpec.Values[i])
								}
								prettyType := ""
								if valueSpec.Type != nil {
									prettyType = pretty(valueSpec.Type)
								}
								handler(
									getLocation(filepath),
									types.Value{
										Name: name.String(),
										Docs: types.Docs{
											Text: valueDeclaration.Doc.Text() + valueSpec.Doc.Text(),
										},
										Const: valueDeclaration.Tok == token.CONST,
										Value: prettyValue,
										Type:  prettyType,
									},
								)
							}
						}
					}
				}
			}
		}
		return false
	}
}
