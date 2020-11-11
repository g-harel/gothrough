package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/gothrough/internal/types"
)

// Location is a generic representation of where a type was found.
type Location struct {
	PackageName       string
	PackageImportPath string
	SourceFile        string
}

// TypeHandlers is a collection of functions that accept all type types.
type TypeHandlers struct {
	Value     func(Location, types.Value)
	Function  func(Location, types.Function)
	Interface func(Location, types.Interface)
}

// Types walks the given directory to extract all the types and pass them on to the handlers.
func Types(srcDir string, handlers TypeHandlers) error {
	// Collect all types in the provided directory.
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
		visit(pathname, newValueVisitor(handlers.Value))
		visit(pathname, newFunctionVisitor(handlers.Function))
		visit(pathname, newInterfaceVisitor(handlers.Interface))
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory: %v", err)
	}

	return nil
}
