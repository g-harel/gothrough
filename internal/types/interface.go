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

type Embedded struct {
	Package string
	Name    string
	Docs    string
}

// Interface contains data about the location and shape of an interface.
type Interface struct {
	Name              string
	Docs              string
	Embedded          []Embedded
	Methods           []Method
	PackageName       string
	PackageImportPath string
	SourceFile        string
	SourceLine        int
}
