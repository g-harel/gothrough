package types

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
	Embedded          []string
	Methods           []Method
	PackageName       string
	PackageImportPath string
	SourceFile        string
	SourceLine        int
}
