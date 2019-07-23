package gis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

var root = path.Join(os.Getenv("GOROOT"), "src")

// var root = "/home/g-harel/Documents/dev/gis/main"
var pattern = regexp.MustCompile("^type [A-Z][A-Za-z]* interface ?{")

func List() {
	fmt.Println("looking in", root)

	fs := token.NewFileSet()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.Contains(path, "internal/") {
			return nil
		}
		if strings.Contains(path, "vendor/") {
			return nil
		}
		if strings.Contains(path, "testdata/") {
			return nil
		}
		if strings.Contains(path, "testing/") {
			return nil
		}

		f, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		if 0 == 1 {
			spew.Dump(path, f)
		}

		v := visitor{Path: path}
		ast.Walk(v, f)

		return nil
	})
	if err != nil {
		panic(err)
	}
}

type visitor struct {
	Path string
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if decl, ok := n.(*ast.GenDecl); ok {
		if decl.Tok == token.TYPE {
			for _, spec := range decl.Specs {
				if typ, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typ.Type.(*ast.InterfaceType); ok {
						name := typ.Name.String()
						imprt := strings.TrimPrefix(v.Path, root)
						imprt = strings.TrimPrefix(imprt, "/")
						imprt = strings.TrimSuffix(imprt, ".go")
						if unicode.IsUpper([]rune(name)[0]) {
							fmt.Println(imprt, name)
						}
					}
				}
			}
		}
	}

	return v
}
