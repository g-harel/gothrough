package extract

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/g-harel/gothrough/internal/types"
)

func newValueVisitor(handler func(Location, types.Value)) visitFunc {
	return func(filepath string, n ast.Node, fset *token.FileSet) bool {
		if n == nil {
			return true
		}

		if valueDeclaration, ok := n.(*ast.GenDecl); ok {
			if valueDeclaration.Tok == token.CONST || valueDeclaration.Tok == token.VAR {
				hasIota := false
				for _, spec := range valueDeclaration.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range valueSpec.Names {
							if name.IsExported() {
								prettyValue := ""
								if len(valueSpec.Values) > i {
									prettyValue = pretty(valueSpec.Values[i])
								}
								if strings.Contains(prettyValue, "iota") {
									hasIota = true
									prettyValue = "iota"
								}
								prettyType := ""
								if valueSpec.Type != nil {
									prettyType = pretty(valueSpec.Type)
								}
								if hasIota && prettyValue == "" && prettyType == "" {
									prettyValue = "iota"
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
