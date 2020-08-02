package interfaces

import (
	"fmt"
	"strings"
)

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

// TODO docs
func (i *Interface) Pretty() string {
	methods := []string{}
	for _, method := range i.Methods {
		arguments := []string{}
		for _, argument := range method.Arguments {
			if argument.Name != "" {
				arguments = append(arguments, fmt.Sprintf("%v %v", argument.Name, argument.Type))
			} else {
				arguments = append(arguments, argument.Type)
			}
		}

		returnValues := []string{}
		for _, returnValue := range method.ReturnValues {
			if returnValue.Name != "" {
				returnValues = append(returnValues, fmt.Sprintf("%v %v", returnValue.Name, returnValue.Type))
			} else {
				returnValues = append(returnValues, returnValue.Type)
			}
		}

		result := fmt.Sprintf("%v(%v)", method.Name, strings.Join(arguments, ", "))
		if len(returnValues) > 1 {
			result += fmt.Sprintf(" (%v)", strings.Join(returnValues, ", "))
		} else if len(returnValues) == 1 {
			result += " " + returnValues[0]
		}

		methods = append(methods, result)
	}
	// TODO empty interfaces.
	return fmt.Sprintf("type %v interface {\n\t%v\n}\n", i.Name, strings.Join(methods, "\n\t"))
}
