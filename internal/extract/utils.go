package extract

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"path"
	"strings"

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

func getLocation(filepath string) Location {
	pathname, filename := path.Split(filepath)

	// Attempt to detect source dir by looking for the closest "src" directory.
	pathParts := strings.Split(path.Clean(pathname), "/src")
	srcDir := path.Join(pathParts[:len(pathParts)-1]...) + "/src/"

	relativePath := strings.TrimPrefix(path.Dir(filepath), srcDir)

	return Location{
		PackageName:       path.Base(relativePath),
		PackageImportPath: relativePath,
		SourceFile:        filename,
	}
}
