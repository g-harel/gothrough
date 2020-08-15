package interfaces

import (
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

func (ifc *Interface) PrettyTokens() []Token {
	tokens := []Token{}

	interfaceDocs := prettyDocs(ifc.Docs)
	if len(interfaceDocs) > 0 {
		for _, line := range interfaceDocs {
			tokens = append(tokens,
				Token{line, kindComment},
				tokenNewline)
		}
	}

	tokens = append(tokens,
		Token{"type", kindKeyword},
		tokenSpace,
		Token{ifc.Name, kindInterfaceName},
		tokenSpace,
		Token{"interface", kindKeyword},
		tokenSpace,
		Token{"{", kindPunctuation})

	if len(ifc.Embedded) == 0 && len(ifc.Methods) == 0 {
		tokens = append(tokens,
			Token{"}", kindPunctuation},
			tokenNewline)
		return tokens
	}

	tokens = append(tokens,
		tokenNewline)

	for _, embedded := range ifc.Embedded {
		// TODO add space if they could have docs.
		tokens = append(tokens,
			tokenIndent,
			Token{embedded, kindEmbeddedName},
			tokenNewline)
	}

	for i, method := range ifc.Methods {
		// Add newline before definition in some situations.
		prevWasEmbedded := i == 0 && len(ifc.Embedded) > 0
		prevHadDocs := i > 0 && ifc.Methods[i-1].Docs != "" // TODO update if embedded can have docs.
		selfIsNotFirstAndHasDocs := (i != 0 || prevWasEmbedded) && method.Docs != ""
		if prevHadDocs || selfIsNotFirstAndHasDocs {
			tokens = append(tokens,
				tokenNewline)
		}

		methodDocs := prettyDocs(method.Docs)
		if len(methodDocs) > 0 {
			for _, line := range methodDocs {
				tokens = append(tokens,
					tokenIndent,
					Token{line, kindComment},
					tokenNewline)
			}
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

func prettyDocs(docs string) []string {
	if docs == "" {
		return []string{}
	}

	docLines := strings.Split(strings.TrimSpace(docs), "\n")
	if len(docLines) < 1 {
		return []string{}
	}

	for i, line := range docLines {
		docLines[i] = "// " + line
	}
	return docLines
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
