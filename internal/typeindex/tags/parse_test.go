package tags_test

import (
	"testing"

	"github.com/g-harel/gothrough/internal/typeindex/tags"
)

type Expected struct {
	Words string
	Tags  map[string][]string
}

func TestParse(t *testing.T) {
	tt := map[string]struct {
		Input    string
		Expected Expected
	}{
		"empty": {
			Input:    "",
			Expected: Expected{},
		},
		"trim words": {
			Input: " abc ",
			Expected: Expected{
				Words: "abc",
			},
		},
		"single tags": {
			Input: "test:true",
			Expected: Expected{
				Tags: map[string][]string{"test": {"true"}},
			},
		},
		"multiple tags": {
			Input: "abc:123 xyz:456",
			Expected: Expected{
				Tags: map[string][]string{"abc": {"123"}, "xyz": {"456"}},
			},
		},
		"duplicate tags": {
			Input: "a:1 a:2",
			Expected: Expected{
				Tags: map[string][]string{"a": {"1", "2"}},
			},
		},
		"tags and words": {
			Input: "test:false abc test:true",
			Expected: Expected{
				Words: "abc",
				Tags:  map[string][]string{"test": {"false", "true"}},
			},
		},
		"multiple colons": {
			Input: "test:test:test",
			Expected: Expected{
				Tags: map[string][]string{"test": {"test:test"}},
			},
		},
		"duplicate values": {
			Input: "test:a test:a",
			Expected: Expected{
				Tags: map[string][]string{"test": {"a"}},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual := tags.Parse(tc.Input)
			if tc.Expected.Words != actual.GetWords() {
				t.Fatalf("expected/actual words do not match\n'%v'\n'%v'", tc.Expected.Words, actual.GetWords())
			}
			for tag, values := range tc.Expected.Tags {
				actualValues := actual.GetTags(tag)
				for i := range values {
					if values[i] != actualValues[i] {
						t.Fatalf("expected/actual tags do not match\n'%v'\n'%v'", values, actualValues)
					}
				}
			}
		})
	}
}
