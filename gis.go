package gis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

type Interface struct {
	Name              string
	PackageName       string
	PackageImportPath string
}

func (i *Interface) String() string {
	return fmt.Sprintf("import \"%v\"\n%v.%v\n", i.PackageImportPath, i.PackageName, i.Name)
}

func Search(dir string) ([]Interface, error) {
	fs := token.NewFileSet()

	err := filepath.Walk(dir, func(pathname string, info os.FileInfo, err error) error {
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
		if strings.Contains(pathname, "testing/") {
			return nil
		}

		f, err := parser.ParseFile(fs, pathname, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path.Dir(pathname), dir)
		relativePath = strings.TrimPrefix(relativePath, "/")
		v := visitor{RelativePath: relativePath}
		ast.Walk(v, f)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk: %v", err)
	}

	return nil, nil
}

type visitor struct {
	RelativePath string
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if typeDeclaration, ok := n.(*ast.GenDecl); ok {
		if typeDeclaration.Tok == token.TYPE {
			for _, spec := range typeDeclaration.Specs {
				if interfaceSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := interfaceSpec.Type.(*ast.InterfaceType); ok {
						name := interfaceSpec.Name.String()
						packageName := path.Base(v.RelativePath)
						if unicode.IsUpper([]rune(name)[0]) {
							i := Interface{
								Name:              name,
								PackageName:       packageName,
								PackageImportPath: v.RelativePath,
							}
							fmt.Println(i.String())
						}
					}
				}
			}
		}
	}

	return v
}
