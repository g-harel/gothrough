package types

// Interface contains data about the location and shape of an interface.
type Interface struct {
	Name              string
	Docs              string
	Embedded          []EmbeddedInterface
	Methods           []MethodSignature
	PackageName       string
	PackageImportPath string
	SourceFile        string
	SourceLine        int
}

var _ Prettier = &Interface{}

func (ifc *Interface) Pretty() string {
	return flattenTokens(ifc.PrettyTokens())
}

func (ifc *Interface) PrettyTokens() []Token {
	// TODO Create tokenlist type with helper methods.
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

	for i, embedded := range ifc.Embedded {
		// Add newline before definition in some situations.
		prevHadDocs := i > 0 && ifc.Embedded[i-1].Docs != ""
		isNotFirstAndHasDocs := i != 0 && embedded.Docs != ""
		if prevHadDocs || isNotFirstAndHasDocs {
			tokens = append(tokens,
				tokenNewline)
		}

		embeddedDocs := prettyDocs(embedded.Docs)
		if len(embeddedDocs) > 0 {
			for _, line := range embeddedDocs {
				tokens = append(tokens,
					tokenIndent,
					Token{line, kindComment},
					tokenNewline)
			}
		}

		tokens = append(tokens,
			tokenIndent)
		if embedded.Package != "" {
			tokens = append(tokens,
				Token{embedded.Package, kindEmbeddedPackage},
				Token{".", kindPunctuation})
		}
		tokens = append(tokens,
			Token{embedded.Name, kindEmbeddedName},
			tokenNewline)
	}

	for i, method := range ifc.Methods {
		// Add newline before definition in some situations.
		prevWasEmbedded := i == 0 && len(ifc.Embedded) > 0
		prevEmbeddedHadDocs := len(ifc.Embedded) > 0 && ifc.Embedded[len(ifc.Embedded)-1].Docs != ""
		prevMethodHadDocs := i > 0 && ifc.Methods[i-1].Docs != ""
		isNotFirstAndHasDocs := (i != 0 || prevWasEmbedded) && method.Docs != ""
		if prevEmbeddedHadDocs || prevMethodHadDocs || isNotFirstAndHasDocs {
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
