package types

// Type is a placeholder for a pointer to any other type in this package.
type Type interface{}

type Docs struct {
	Text string
}

type Field struct {
	Name string
	Docs Docs
	Type string
}

type Function struct {
	Name         string
	Docs         Docs
	Arguments    []Field
	ReturnValues []Field
}

type Reference struct {
	Package string
	Name    string
	Docs    Docs
}

type Interface struct {
	Name     string
	Docs     Docs
	Embedded []Reference
	Methods  []Function
}

type Value struct {
	Name  string
	Docs  Docs
	Type  string
	Value string
}
