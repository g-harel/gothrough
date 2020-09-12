package types

// Type is a placeholder for a pointer to any other type in this package.
type Type = interface{}

type Field struct {
	Name string
	Docs string
	Type string
}

type MethodSignature struct {
	Name         string
	Docs         string
	Arguments    []Field
	ReturnValues []Field
}

// TODO refactor to a more generic reference.
type EmbeddedInterface struct {
	Package string
	Name    string
	Docs    string
}

// Interface contains data about the location and shape of an interface.
type Interface struct {
	Name              string
	Docs              string
	Embedded          []EmbeddedInterface
	Methods           []MethodSignature
	PackageName       string
	PackageImportPath string
	SourceFile        string
	SourceLine        int
}
