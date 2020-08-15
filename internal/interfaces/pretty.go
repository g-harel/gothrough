package interfaces

import (
	"fmt"
	"strings"
)

const (
	kindComment       = "comment"
	kindEmbeddedName  = "embedded_name"
	kindFieldName     = "field_name"
	kindInterfaceName = "interface_name"
	kindKeyword       = "keyword"
	kindMethodName    = "method_name"
	kindPunctuation   = "punctuation"
	kindFieldType     = "field_type"
	kindWhitespace    = "whitespace"
)

var (
	tokenIndent  = Token{"\t", kindWhitespace}
	tokenNewline = Token{"\n", kindWhitespace}
	tokenSpace   = Token{" ", kindWhitespace}
)

type Token struct {
	Text string
	Kind string
}

func (i *Interface) PrettyTokens() []Token {
	tokens := []Token{}

	prettyInterfaceDocs := prettyDocsLine(i.Name, i.Docs)
	if prettyInterfaceDocs != "" {
		tokens = append(tokens,
			Token{prettyInterfaceDocs, kindComment},
			tokenNewline)
	}

	tokens = append(tokens,
		Token{"type", kindKeyword},
		tokenSpace,
		Token{i.Name, kindInterfaceName},
		tokenSpace,
		Token{"interface", kindKeyword},
		tokenSpace,
		Token{"{", kindPunctuation})

	if len(i.Embedded) == 0 && len(i.Methods) == 0 {
		tokens = append(tokens,
			Token{"}", kindPunctuation},
			tokenNewline)
		return tokens
	}

	tokens = append(tokens,
		tokenNewline)

	for _, embedded := range i.Embedded {
		tokens = append(tokens,
			tokenIndent,
			Token{embedded, kindEmbeddedName},
			tokenNewline)
	}

	for _, method := range i.Methods {
		prettyMethodDocs := prettyDocsLine(method.Name, method.Docs)
		if prettyMethodDocs != "" {
			tokens = append(tokens,
				tokenIndent,
				Token{prettyMethodDocs, kindComment},
				tokenNewline)
		}

		tokens = append(tokens,
			tokenIndent,
			Token{method.Name, kindMethodName},
			Token{"(", kindPunctuation})

		tokens = append(tokens,
			tokenizeFieldList(method.Arguments)...)

		tokens = append(tokens,
			Token{")", kindPunctuation})

		if len(method.ReturnValues) == 1 {
			tokens = append(tokens,
				tokenSpace)
			tokens = append(tokens,
				tokenizeFieldList(method.ReturnValues)...)
		} else if len(method.ReturnValues) > 1 {
			tokens = append(tokens,
				tokenSpace,
				Token{"(", kindPunctuation})
			tokens = append(tokens,
				tokenizeFieldList(method.ReturnValues)...)
			tokens = append(tokens,
				Token{")", kindPunctuation})
		}

		tokens = append(tokens,
			tokenNewline)
	}

	tokens = append(tokens,
		Token{"}", kindPunctuation},
		tokenNewline)

	return tokens
}

func (i *Interface) Pretty() string {
	output := ""
	for _, token := range i.PrettyTokens() {
		output += token.Text
	}
	return output
}

// TODO break lines that are longer than 80.
func prettyDocsLine(name, docs string) string {
	if docs == "" {
		return ""
	}

	docLines := strings.Split(docs, "\n")
	if len(docLines) < 1 {
		return ""
	}

	// Return first line if it ends with a period.
	// This would better capture doc lines that include periods than the next block.
	if strings.HasSuffix(docLines[0], ".") {
		return fmt.Sprintf("// %v", docLines[0])
	}

	// Return first sentence if it starts with identifier.
	if strings.HasPrefix(docLines[0], name) {
		allDocs := strings.Join(docLines, " ")
		firstSentence := strings.Split(allDocs, ".")[0]
		return fmt.Sprintf("// %v.", firstSentence)
	}

	return ""
}

func tokenizeFieldList(fields []Field) []Token {
	tokens := []Token{}
	for i, field := range fields {
		if i > 0 {
			tokens = append(tokens,
				Token{",", kindPunctuation},
				tokenSpace)
		}
		if field.Name != "" {
			tokens = append(tokens,
				Token{field.Name, kindFieldName},
				tokenSpace)
		}
		tokens = append(tokens,
			Token{field.Type, kindFieldType})
	}
	return tokens
}
