package typeindex

import (
	"strings"

	"github.com/g-harel/gothrough/internal/extract"
)

func (idx *Index) insertLocation(id int, location extract.Location) {
	// Index on package path and name.
	idx.textIndex.Insert(id, confidenceHigh, location.PackageName)
	importPathParts := strings.Split(location.PackageImportPath, "/")
	if len(importPathParts) > 1 {
		idx.textIndex.Insert(id, confidenceLow/len(importPathParts), importPathParts...)
	}
}
