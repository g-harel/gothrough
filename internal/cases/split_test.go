package cases_test

import (
	"testing"

	"github.com/g-harel/gothrough/internal/cases"
)

func TestSplit(t *testing.T) {
	tt := map[string]struct {
		Input    string
		Expected []string
	}{
		"empty": {
			Input:    "",
			Expected: []string{},
		},
		"camel simple double": {
			Input:    "TestCase",
			Expected: []string{"Test", "Case"},
		},
		"camel simple n-tuple": {
			Input:    "TestCaseThatIsLongerThanTheOtherOne",
			Expected: []string{"Test", "Case", "That", "Is", "Longer", "Than", "The", "Other", "One"},
		},
		"camel prefix acronym": {
			Input:    "HTTPTestCase",
			Expected: []string{"HTTP", "Test", "Case"},
		},
		"camel postfix acronym": {
			Input:    "ServeHTTP",
			Expected: []string{"Serve", "HTTP"},
		},
		"camel surrounded acronym": {
			Input:    "AbcABCAbc",
			Expected: []string{"Abc", "ABC", "Abc"},
		},
		"snake simple double": {
			Input:    "test_case",
			Expected: []string{"test", "case"},
		},
		"mixed camel/snake": {
			Input:    "AbcDef_xyz123",
			Expected: []string{"Abc", "Def", "xyz123"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual := cases.Split(tc.Input)
			if len(tc.Expected) != len(actual) {
				t.Fatalf("expected/actual do not match\n'%v'\n'%v'", tc.Expected, actual)
			}
			for i := range tc.Expected {
				if len(tc.Expected) != len(actual) || tc.Expected[i] != actual[i] {
					t.Fatalf("expected/actual do not match\n'%v'\n'%v'", tc.Expected, actual)
				}
			}
		})
	}
}
