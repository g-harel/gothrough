package interfaces

import (
	"fmt"
)

// TODO add docs if available
type Field struct {
	Name string
	Docs string
	Type string
}

type Method struct {
	Name         string
	Docs         string
	Arguments    []Field
	ReturnValues []Field
}

// Interface contains data about the location and shape of an interface.
type Interface struct {
	Name              string
	Docs              string
	Methods           []Method
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

func (i *Interface) DocLink() string {
	return fmt.Sprintf("https://golang.org/pkg/%v#%v", i.PackageImportPath, i.Name)
}
