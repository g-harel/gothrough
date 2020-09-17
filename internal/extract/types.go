package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/gothrough/internal/types"
)

type TypeHandlers struct {
	Interface func(types.Interface)
}

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
		visit(pathname, newInterfaceVisitor(handlers.Interface))
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory: %v", err)
	}

	return nil
}
