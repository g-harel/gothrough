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
	Name    string
	Docs    string
	Methods []Method
	// TODO read + index + print
	Embedded          []string
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
		// Create pretty method arguments.
		arguments := []string{}
		for _, argument := range method.Arguments {
			if argument.Name != "" {
				arguments = append(arguments, fmt.Sprintf("%v %v", argument.Name, argument.Type))
			} else {
				arguments = append(arguments, argument.Type)
			}
		}
		prettyArguments := strings.Join(arguments, ", ")

		// Create pretty method return values.
		returnValues := []string{}
		for _, returnValue := range method.ReturnValues {
			if returnValue.Name != "" {
				returnValues = append(returnValues, fmt.Sprintf("%v %v", returnValue.Name, returnValue.Type))
			} else {
				returnValues = append(returnValues, returnValue.Type)
			}
		}
		prettyReturnValues := ""
		if len(returnValues) > 1 {
			prettyReturnValues = fmt.Sprintf(" (%v)", strings.Join(returnValues, ", "))
		} else if len(returnValues) == 1 {
			prettyReturnValues = fmt.Sprintf(" %v", returnValues[0])
		}

		methodDocs := ""
		if method.Docs != "" {
			// TODO docs printing helper that only shows if first line is a sentence.
			methodDocs = "// " + strings.Split(method.Docs, "\n")[0] + "\n"
		}

		prettyMethod := fmt.Sprintf("%v%v(%v)%v", methodDocs, method.Name, prettyArguments, prettyReturnValues)

		methods = append(methods, prettyMethod)
	}

	interfaceDocs := ""
	if i.Docs != "" {
		interfaceDocs = "// " + strings.Split(i.Docs, "\n")[0] + "\n"
	}

	interfaceBody := ""
	if len(i.Methods) > 0 {
		allMethods := strings.Join(methods, "\n")
		indentedMethods := prefixLines(allMethods, "\t")
		interfaceBody = fmt.Sprintf("\n%v\n", indentedMethods)
	}

	return fmt.Sprintf("%vtype %v interface {%v}", interfaceDocs, i.Name, interfaceBody)
}

func prefixLines(s string, p string) string {
	out := []string{}
	for _, line := range strings.Split(s, "\n") {
		out = append(out, p+line)
	}
	return strings.Join(out, "\n")
}
