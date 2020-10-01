package typeindex

import (
	"strings"

	"github.com/g-harel/gothrough/internal/extract"
)

func (idx *Index) insertLocation(id int, location extract.Location) {
	// Index on package path and source file.
	importPathParts := strings.Split(location.PackageImportPath, "/")
	idx.textIndex.Insert(id, confidenceHigh, location.PackageName)
	idx.textIndex.Insert(id, confidenceLow, strings.TrimSuffix(location.SourceFile, ".go"))
	if len(importPathParts) > 1 {
		idx.textIndex.Insert(id, confidenceLow/len(importPathParts), importPathParts...)
	}
}
