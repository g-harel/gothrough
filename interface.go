package gis

import (
	"fmt"
)

// Interface represents the location of a discovered interface.
type Interface struct {
	Name              string
	Methods           []string
	Body              string
	PackageName       string
	PackageImportPath string
	SourceFile        string
	SourceLine        int
}

func (i *Interface) Address() string {
	return fmt.Sprintf("%v/%v:%v", i.PackageImportPath, i.SourceFile, i.SourceLine)
}

// String returns a string representation of the interface.
func (i *Interface) String() string {
	return fmt.Sprintf("%v %v (%v)", i.Name, i.Methods, i.Address())
	// return fmt.Sprintf("%v\nimport \"%v\"\n%v.%v\n%v\n", i.Address(), i.PackageImportPath, i.PackageName, i.Name, i.Body)
}

// Equals compares Interfaces to determine if they are equal.
func (i *Interface) Equals(alt *Interface) bool {
	return i.Name == alt.Name && i.PackageImportPath == i.PackageImportPath
}
