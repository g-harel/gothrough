package gis

import (
	"fmt"
)

// Interface represents the location of a discovered interface.
type Interface struct {
	Name              string
	PackageName       string
	PackageImportPath string
}

// String returns a string representation of the interface.
func (i *Interface) String() string {
	return fmt.Sprintf("import \"%v\"\n%v.%v\n", i.PackageImportPath, i.PackageName, i.Name)
}

// Equals compares Interfaces to determine if they are equal.
func (i *Interface) Equals(alt *Interface) bool {
	return i.Name == alt.Name && i.PackageImportPath == i.PackageImportPath
}
