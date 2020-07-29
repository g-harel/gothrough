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

func pretty(node interface{}) string {
	buf := bytes.NewBufferString("")
	renderer := token.NewFileSet()
	format.Node(buf, renderer, node)
	return buf.String()
}

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

								// Collect methods.
								methods := []interfaces.Method{}
								for _, method := range interfaceType.Methods.List {
									arguments := []interfaces.Argument{}
									returnValues := []interfaces.ReturnValue{}
									if methodType, ok := method.Type.(*ast.FuncType); ok {
										if methodType.Params != nil {
											for _, param := range methodType.Params.List {
												for _, paramNames := range param.Names {
													arguments = append(arguments, interfaces.Argument{
														Name: paramNames.Name,
														Type: pretty(param.Type),
													})
												}
											}
										}
										if methodType.Results != nil {
											// TODO check that unnamed return values are captured
											for _, returnValue := range methodType.Results.List {
												for _, returnValueNames := range returnValue.Names {
													returnValues = append(returnValues, interfaces.ReturnValue{
														Name: returnValueNames.Name,
														Type: pretty(returnValue.Type),
													})
												}
											}
										}
									}
									for _, methodName := range method.Names {
										if methodName.IsExported() {
											methods = append(methods, interfaces.Method{
												Name:         methodName.String(),
												Docs:         method.Doc.Text(),
												Arguments:    arguments,
												ReturnValues: returnValues,
											})
										}
									}
								}

								handler(interfaces.Interface{
									Name:              name,
									Docs:              typeDeclaration.Doc.Text(),
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
