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

		prettyMethodDocs := prettyDocsLine(method.Name, method.Docs)
		prettyMethod := fmt.Sprintf("%v%v(%v)%v", prettyMethodDocs, method.Name, prettyArguments, prettyReturnValues)

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

	prettyInterfaceDocs := prettyDocsLine(i.Name, i.Docs)
	return fmt.Sprintf("%vtype %v interface {%v}", prettyInterfaceDocs, i.Name, interfaceBody)
}

func prettyDocsLine(name, docs string) string {
	if docs == "" {
		return ""
	}

	docLines := strings.Split(docs, "\n")
	if len(docLines) < 1 {
		return ""
	}

	// Return first line if it ends with a period.
	// This would better capture doc lines that include periods than the next block.
	if strings.HasSuffix(docLines[0], ".") {
		return fmt.Sprintf("// %v\n", docLines[0])
	}

	// Return first sentence if it starts with identifier.
	if strings.HasPrefix(docLines[0], name) {
		allDocs := strings.Join(docLines, " ")
		firstSentence := strings.Split(allDocs, ".")[0]
		return fmt.Sprintf("// %v.\n", firstSentence)
	}

	return ""
}

func prefixLines(s string, p string) string {
	out := []string{}
	for _, line := range strings.Split(s, "\n") {
		out = append(out, p+line)
	}
	return strings.Join(out, "\n")
}
