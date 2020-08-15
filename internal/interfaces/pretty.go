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

func (i *Interface) PrettyTokens() []Token {
	tokens := []Token{}

	interfaceDocs := prettyDocs(i.Name, i.Docs, 80)
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
		methodDocs := prettyDocs(method.Name, method.Docs, 76)
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

// TODO break lines that are longer than max length.
func prettyDocs(name, docs string, maxLength int) []string {
	if docs == "" {
		return []string{}
	}

	if maxLength <= 3 {
		return []string{}
	}

	docLines := strings.Split(docs, "\n")
	if len(docLines) < 1 {
		return []string{}
	}

	// Find contents of the docs that should be kept.
	finalDocs := ""
	if strings.HasSuffix(docLines[0], ".") {
		// Use first line if it ends with a period.
		// This would better capture doc lines that include periods vs. the next block.
		finalDocs = docLines[0]
	} else if strings.HasPrefix(docLines[0], name) {
		// Use entire first sentence if docs start with target name.
		allDocs := strings.Join(docLines, " ")
		finalDocs = strings.Split(allDocs, ".")[0] + "."
	} else {
		return []string{}
	}

	// Break docs into lines along whitespace.
	actualMaxLength := maxLength - 3
	words := strings.Fields(finalDocs)
	lines := []string{""}
	for _, word := range words {
		lastLineIndex := len(lines) - 1
		if len(lines[lastLineIndex])+len(word)+1 < actualMaxLength {
			// If word fits, add it to the last line.
			lines[lastLineIndex] += " " + word
		} else {
			// If doesn't fit, add it to a new line.
			lines = append(lines, " "+word)
		}
	}

	for i, line := range lines {
		lines[i] = "//" + line
	}

	return lines
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
