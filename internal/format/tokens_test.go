package format_test

import (
	"testing"

	"github.com/g-harel/gothrough/internal/format"
)

func TestSnippet(t *testing.T) {
	tt := map[string]struct {
		Input    func() *format.Snippet
		Expected string
	}{
		"empty": {
			Input:    func() *format.Snippet { return format.NewSnippet() },
			Expected: "",
		},
		"single token": {
			Input: func() *format.Snippet {
				s := format.NewSnippet()
				s.Keyword("test")
				return s
			},
			Expected: "test",
		},
		"multiple tokens": {
			Input: func() *format.Snippet {
				s := format.NewSnippet()
				s.Keyword("var")
				s.Space()
				s.Punctuation("=")
				s.Space()
				s.Literal("4")
				return s
			},
			Expected: "var = 4",
		},
		"indent lines": {
			Input: func() *format.Snippet {
				s := format.NewSnippet()
				s.DeclName("Test")
				s.Newline()
				s.EmbeddedName("abc")
				s.Newline()
				s.InterfaceName("123")
				s.IndentLines()
				return s
			},
			Expected: `	Test
	abc
	123`,
		},
		"push snippet": {
			Input: func() *format.Snippet {
				inside := format.NewSnippet()
				inside.Indent()
				inside.Literal("'test'")

				s := format.NewSnippet()
				s.Punctuation("{")
				s.Newline()
				s.Push(inside)
				s.Newline()
				s.Punctuation("}")
				return s
			},
			Expected: `{
	'test'
}`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actual := ""
			for _, token := range tc.Input().Dump() {
				actual += token.Text
			}
			if actual != tc.Expected {
				t.Fatalf("expected/actual do not match\n'%v'\n'%v'", tc.Expected, actual)
			}
		})
	}
}
