package types

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
