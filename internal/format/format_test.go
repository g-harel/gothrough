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
		"simple function": {
			Input: &types.Function{
				Name: "Test",
			},
			Expected: "func Test()\n",
		},
		"simple value": {
			Input: &types.Value{
				Name: "Test",
				Type: "string",
			},
			Expected: "var Test string\n",
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
