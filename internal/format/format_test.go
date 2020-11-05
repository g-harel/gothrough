package format_test

import (
	"testing"

	"github.com/g-harel/gothrough/internal/format"
	"github.com/g-harel/gothrough/internal/types"
)

func TestString(t *testing.T) {
	tt := map[string]struct {
		Input    types.Type
		Expected string
	}{
		"simple interface": {
			Input: &types.Interface{
				Name: "Test",
			},
			Expected: "type Test interface {}\n",
		},
		"interface with function and embedded": {
			Input: &types.Interface{
				Name:     "Extender",
				Embedded: []types.Reference{{Package: "io", Name: "Reader"}},
				Methods:  []types.Function{{Name: "Print"}},
			},
			Expected: `type Extender interface {
	io.Reader
	Print()
}
`,
		},
		"interface with docs": {
			Input: &types.Interface{
				Name: "Potato",
				Docs: types.Docs{Text: "Potato does the thing."},
				Embedded: []types.Reference{
					{
						Name: "Tester",
					},
					{
						Name: "Water",
					},
					{
						Package: "earth",
						Name:    "Grower",
						Docs: types.Docs{
							Text: "Grows the thing with the thing.",
						},
					},
				},
				Methods: []types.Function{
					{
						Name: "Pick",
						Docs: types.Docs{
							Text: "Takes the thing from the thing.",
						},
					},
					{
						Name: "Squish",
						ReturnValues: []types.Field{
							{Type: "int"},
						},
					},
					{
						Name: "Eat",
						ReturnValues: []types.Field{
							{Type: "string"},
							{Type: "error"},
						},
					},
				},
			},
			Expected: `// Potato does the thing.
type Potato interface {
	Tester
	Water

	// Grows the thing with the thing.
	earth.Grower

	// Takes the thing from the thing.
	Pick()

	Squish() int
	Eat() (string, error)
}
`,
		},
		"simple function": {
			Input: &types.Function{
				Name: "Test",
			},
			Expected: "func Test()\n",
		},
		"function with single return value": {
			Input: &types.Function{
				Name:         "GetCount",
				ReturnValues: []types.Field{{Type: "int"}},
			},
			Expected: "func GetCount() int\n",
		},
		"function with multiple return values": {
			Input: &types.Function{
				Name:         "Calc",
				ReturnValues: []types.Field{{Type: "int"}, {Type: "error"}},
			},
			Expected: "func Calc() (int, error)\n",
		},
		"function with named return values": {
			Input: &types.Function{
				Name:         "Test",
				ReturnValues: []types.Field{{Name: "expected", Type: "string"}},
			},
			Expected: "func Test() (expected string)\n",
		},
		"function with argument": {
			Input: &types.Function{
				Name:         "Square",
				Arguments:    []types.Field{{Name: "num", Type: "int"}},
				ReturnValues: []types.Field{{Type: "int"}},
			},
			Expected: "func Square(num int) int\n",
		},
		"simple value": {
			Input: &types.Value{
				Name: "Test",
				Type: "string",
			},
			Expected: "var Test string\n",
		},
		"const value": {
			Input: &types.Value{
				Name:  "Count",
				Type:  "int",
				Const: true,
			},
			Expected: "const Count int\n",
		},
		"value with literal": {
			Input: &types.Value{
				Name:  "TEST_NAME",
				Value: "\"abc\"",
			},
			Expected: "var TEST_NAME = \"abc\"\n",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual, err := format.String(tc.Input)
			if err != nil {
				t.Fatalf("format error: %v", err)
				return
			}
			if actual != tc.Expected {
				t.Fatalf("expected/actual do not match\n'%v'\n'%v'", tc.Expected, actual)
			}
		})
	}
}
