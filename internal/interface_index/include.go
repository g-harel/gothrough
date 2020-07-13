package interface_index

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/g-harel/gis/internal/camel"
	"github.com/g-harel/gis/internal/interfaces"
	"github.com/g-harel/gis/internal/parse"
)

// Confidence values for interface info items.
const (
	interfaceNameVal           = 120
	totalInterfaceNameTokenVal = 160
	packageNameVal             = 120
	sourceFileVal              = 10
	totalImportPathPartVal     = 20
	totalMethodNameVal         = 80
	totalMethodNameTokenVal    = 80
)

// Include adds interfaces in the provided src directory to the index.
func (si *Index) Include(srcDir string) error {
	// Store current count to skip already-indexed interfaces.
	startLength := len(si.interfaces)

	// Collect all interfaces in the provided directory.
	err := filepath.Walk(srcDir, func(pathname string, info os.FileInfo, err error) error {
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
		parse.Visit(pathname, parse.NewInterfaceVisitor(func(ifc interfaces.Interface) {
			si.interfaces = append(si.interfaces, &ifc)
		}))
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory: %v", err)
	}

	// Add the interfaces to the internal index with a confidence value.
	for i, ifc := range si.interfaces[startLength:] {
		id := i + startLength

		// Index on interface name.
		si.index.Index(id, interfaceNameVal, ifc.Name)
		nameTokens := camel.Split(ifc.Name)
		if len(nameTokens) > 1 {
			si.index.Index(id, totalInterfaceNameTokenVal/len(nameTokens), nameTokens...)
		}

		// Index on package path and source file.
		importPathParts := strings.Split(ifc.PackageImportPath, "/")
		si.index.Index(id, packageNameVal, ifc.PackageName)
		si.index.Index(id, sourceFileVal, strings.TrimSuffix(ifc.SourceFile, ".go"))
		if len(importPathParts) > 1 {
			si.index.Index(i, totalImportPathPartVal/len(importPathParts), importPathParts...)
		}

		// Index on interface methods.
		methodNameTokens := []string{}
		for _, methodName := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(methodName)...)
		}
		if len(ifc.Methods) > 0 {
			si.index.Index(id, totalMethodNameVal/len(ifc.Methods), ifc.Methods...)
		}
		if len(methodNameTokens) > 0 {
			si.index.Index(id, totalInterfaceNameTokenVal/len(methodNameTokens), methodNameTokens...)
		}
	}

	return nil
}
