package example

import "io"

type test interface {
	Plain()
}

type Test interface {
	io.Reader

	Plain()
	String(string) string
	Int() (int, int)
	Name(arg bool) (count int, err error)
}
