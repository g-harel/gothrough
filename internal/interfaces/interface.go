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
	Name string
	Docs string
	// TODO index
	Embedded          []string
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

		prettyMethod := fmt.Sprintf("%v%v(%v)%v", prettyDocsLine(method.Docs), method.Name, prettyArguments, prettyReturnValues)

		methods = append(methods, prettyMethod)
	}

	interfaceBody := ""
	if len(i.Methods) > 0 || len(i.Embedded) > 0 {
		allEmbedded := strings.Join(i.Embedded, "\n")
		allMethods := strings.Join(methods, "\n")
		allBody := strings.TrimSpace(allEmbedded + "\n" + allMethods)
		indentedBody := prefixLines(allBody, "\t")
		interfaceBody = fmt.Sprintf("\n%v\n", indentedBody)
	}

	return fmt.Sprintf("%vtype %v interface {%v}", prettyDocsLine(i.Docs), i.Name, interfaceBody)
}

func prettyDocsLine(docs string) string {
	if docs == "" {
		return ""
	}

	docLines := strings.Split(docs, "\n")
	if len(docLines) < 1 {
		return ""
	}

	// TODO use entire first sentence even when broken up if first word matches name.
	docLine := docLines[0]
	if !strings.HasSuffix(docLine, ".") {
		return ""
	}

	return fmt.Sprintf("// %v\n", docLine)
}

func prefixLines(s string, p string) string {
	out := []string{}
	for _, line := range strings.Split(s, "\n") {
		out = append(out, p+line)
	}
	return strings.Join(out, "\n")
}
