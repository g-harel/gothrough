package string_index_test

import (
	"testing"

	"github.com/g-harel/gis/internal/string_index"
)

func TestSubstring(t *testing.T) {
	tt := map[string]struct {
		InputString string
		InputSize   int
		Expected    []string
	}{
		"empty": {
			InputString: "",
			InputSize:   2,
			Expected:    []string{},
		},
		"0 size": {
			InputString: "abc",
			InputSize:   0,
			Expected:    []string{},
		},
		"1 size": {
			InputString: "abc",
			InputSize:   1,
			Expected:    []string{"a", "b", "c"},
		},
		"2 size": {
			InputString: "abcde",
			InputSize:   2,
			Expected:    []string{"ab", "bc", "cd", "de"},
		},
		"4 size": {
			InputString: "abcde",
			InputSize:   4,
			Expected:    []string{"abcd", "bcde"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual := string_index.Substrings(tc.InputString, tc.InputSize)
			if len(tc.Expected) != len(actual) {
				t.Fatalf("expected/actual do not match\n%v\n%v", tc.Expected, actual)
			}
			for i := range tc.Expected {
				if tc.Expected[i] != actual[i] {
					t.Fatalf("expected/actual do not match\n%v\n%v", tc.Expected, actual)
				}
			}
		})
	}
}
