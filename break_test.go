package gis_test

import (
	"testing"

	"github.com/g-harel/gis"
)

func TestCamelSplit(t *testing.T) {
	tt := map[string]struct {
		Input    string
		Expected []string
	}{
		"empty": {
			Input:    "",
			Expected: []string{},
		},
		"simple double": {
			Input:    "TestCase",
			Expected: []string{"Test", "Case"},
		},
		"simple n-tuple": {
			Input:    "TestCaseThatIsLongerThanTheOtherOne",
			Expected: []string{"Test", "Case", "That", "Is", "Longer", "Than", "The", "Other", "One"},
		},
		"prefix acronym": {
			Input:    "HTTPTestCase",
			Expected: []string{"HTTP", "Test", "Case"},
		},
		"postfix acronym": {
			Input:    "ServeHTTP",
			Expected: []string{"Serve", "HTTP"},
		},
		"surrounded acronym": {
			Input:    "AbcABCAbc",
			Expected: []string{"Abc", "ABC", "Abc"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual := gis.CamelSplit(tc.Input)
			for i := range tc.Expected {
				if len(tc.Expected) != len(actual) || tc.Expected[i] != actual[i] {
					t.Fatalf("expected/actual do not match\n%v\n%v", tc.Expected, actual)
				}
			}
		})
	}
}
