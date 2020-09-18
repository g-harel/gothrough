package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"github.com/g-harel/gothrough/internal/types"
)

func pretty(node interface{}) string {
	buf := bytes.NewBufferString("")
	renderer := token.NewFileSet()
	format.Node(buf, renderer, node)
	return buf.String()
}

func collectFields(fieldList *ast.FieldList) []types.Field {
	result := []types.Field{}
	if fieldList != nil {
		for _, field := range fieldList.List {
			if len(field.Names) == 0 {
				result = append(result, types.Field{
					Name: "",
					Docs: types.Docs{Text: field.Doc.Text()},
					Type: pretty(field.Type),
				})
				continue
			}
			for _, fieldNames := range field.Names {
				result = append(result, types.Field{
					Name: fieldNames.Name,
					Docs: types.Docs{Text: field.Doc.Text()},
					Type: pretty(field.Type),
				})
			}
		}
	}
	return result
}
