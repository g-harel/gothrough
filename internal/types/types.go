package types

// Type is a placeholder for a pointer to any other type in this package.
type Type interface{}

// Docs represents documentation text attached to other types.
type Docs struct {
	Text string
}

// Field represents a function argument or return value.
// Docs and name are optional.
type Field struct {
	Name string
	Docs Docs
	Type string
}

// Function represents a function type.
// Docs, arguments and return values are optional.
type Function struct {
	Name         string
	Docs         Docs
	Arguments    []Field
	ReturnValues []Field
}

// Reference represents a reference to another value.
// Docs and package are optional.
type Reference struct {
	Package string
	Name    string
	Docs    Docs
}

// Interface represents an interface type.
// Docs, embedded and methods are optional.
type Interface struct {
	Name     string
	Docs     Docs
	Embedded []Reference
	Methods  []Function
}

// Value represents a const or var declaration.
// Docs and const are optional. Exactly one of type or value should be specified.
type Value struct {
	Name  string
	Docs  Docs
	Const bool
	Type  string
	Value string
}
