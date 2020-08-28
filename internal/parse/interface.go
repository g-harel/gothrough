package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/g-harel/gothrough/internal/types"
)

// TODO make more utils and move somewhere else.
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
					Docs: field.Doc.Text(),
					Type: pretty(field.Type),
				})
				continue
			}
			for _, fieldNames := range field.Names {
				result = append(result, types.Field{
					Name: fieldNames.Name,
					Docs: field.Doc.Text(),
					Type: pretty(field.Type),
				})
			}
		}
	}
	return result
}

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
								methods := []types.Method{}
								embedded := []string{}
								for _, method := range interfaceType.Methods.List {
									if identType, ok := method.Type.(*ast.Ident); ok {
										// TODO docs?
										embedded = append(embedded, identType.Name)
										continue
									}
									if selectorExprType, ok := method.Type.(*ast.SelectorExpr); ok {
										embedded = append(embedded, fmt.Sprintf("%v.%v", selectorExprType.X, selectorExprType.Sel))
										continue
									}

									arguments := []types.Field{}
									returnValues := []types.Field{}
									if methodType, ok := method.Type.(*ast.FuncType); ok {
										arguments = collectFields(methodType.Params)
										returnValues = collectFields(methodType.Results)
									}
									for _, methodName := range method.Names {
										if methodName.IsExported() {
											methods = append(methods, types.Method{
												Name:         methodName.String(),
												Docs:         method.Doc.Text(),
												Arguments:    arguments,
												ReturnValues: returnValues,
											})
										}
									}
								}

								handler(types.Interface{
									Name:              name,
									Docs:              typeDeclaration.Doc.Text(),
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

// FindInterfaces adds interfaces in the provided src directory to the index.
func FindInterfaces(srcDir string) ([]*types.Interface, error) {
	found := []*types.Interface{}

	// Collect all interfaces in the provided directory.
	err := filepath.Walk(srcDir, func(pathname string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(pathname, ".go") {
			return nil
		}
		if strings.HasSuffix(pathname, "_test.go") {
			return nil
		}
		if strings.Contains(pathname, "internal/") {
			return nil
		}
		if strings.Contains(pathname, "vendor/") {
			return nil
		}
		if strings.Contains(pathname, "testdata/") {
			return nil
		}
		visit(pathname, newInterfaceVisitor(func(ifc types.Interface) {
			found = append(found, &ifc)
		}))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %v", err)
	}

	return found, nil
}
