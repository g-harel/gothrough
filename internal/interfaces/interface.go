package interfaces

import (
	"fmt"
)

// Interface contains data about the location and shape of an interface.
type Interface struct {
	Name              string
	Methods           []string
	Printed              string
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
}

func (i *Interface) Pretty() string {
	return fmt.Sprintf("%v\npackage \"%v\"\n// %v\n%v\n", i.Address(), i.PackageImportPath, i.DocLink(), i.Printed)
}

func (i *Interface) DocLink() string {
	return fmt.Sprintf("https://golang.org/pkg/%v#%v", i.PackageImportPath, i.Name)
}
