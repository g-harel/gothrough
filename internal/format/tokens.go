package format

const (
	textIndent  = "\t"
	textNewline = "\n"
	textSpace   = " "
)

const (
	kindComment         = "comment"
	kindDeclName        = "decl_name"
	kindEmbeddedName    = "embedded_name"
	kindEmbeddedPackage = "embedded_package"
	kindFieldName       = "field_name"
	kindInterfaceName   = "interface_name"
	kindLiteral         = "literal"
	kindKeyword         = "keyword"
	kindPunctuation     = "punctuation"
	kindFieldType       = "field_type"
	kindWhitespace      = "whitespace"
)

type Token struct {
	Text string
	Kind string
}

type Snippet struct {
	tokens []Token
}

func NewSnippet() *Snippet {
	return &Snippet{[]Token{}}
}

func (snippet *Snippet) Dump() []Token {
	return snippet.tokens
}

func (snippet *Snippet) push(token Token) {
	snippet.tokens = append(snippet.tokens, token)
}

func (snippet *Snippet) Push(s *Snippet) {
	snippet.tokens = append(snippet.tokens, s.Dump()...)
}

func (snippet *Snippet) Indent() {
	snippet.push(Token{textIndent, kindWhitespace})
}

func (snippet *Snippet) IndentLines() {
	indentToken := Token{textIndent, kindWhitespace}

	if len(snippet.tokens) == 0 {
		return
	}
	snippet.tokens = append([]Token{indentToken}, snippet.tokens...)

	for i := 0; i < len(snippet.tokens); i++ {
		if i == len(snippet.tokens)-1 {
			continue
		}
		if snippet.tokens[i].Text == textNewline && snippet.tokens[i].Kind == kindWhitespace {
			snippet.tokens = append(snippet.tokens[:i+1], snippet.tokens[i:]...)
			snippet.tokens[i+1] = indentToken
		}
	}
}

func (snippet *Snippet) Newline() {
	snippet.push(Token{textNewline, kindWhitespace})
}

func (snippet *Snippet) Space() {
	snippet.push(Token{textSpace, kindWhitespace})
}

func (snippet *Snippet) Comment(text string) {
	snippet.push(Token{text, kindComment})
}

func (snippet *Snippet) EmbeddedName(text string) {
	snippet.push(Token{text, kindEmbeddedName})
}

func (snippet *Snippet) EmbeddedPackage(text string) {
	snippet.push(Token{text, kindEmbeddedPackage})
}

func (snippet *Snippet) FieldName(text string) {
	snippet.push(Token{text, kindFieldName})
}

func (snippet *Snippet) InterfaceName(text string) {
	snippet.push(Token{text, kindInterfaceName})
}

func (snippet *Snippet) Literal(text string) {
	snippet.push(Token{text, kindLiteral})
}

func (snippet *Snippet) Keyword(text string) {
	snippet.push(Token{text, kindKeyword})
}

func (snippet *Snippet) DeclName(text string) {
	snippet.push(Token{text, kindDeclName})
}

func (snippet *Snippet) Punctuation(text string) {
	snippet.push(Token{text, kindPunctuation})
}

func (snippet *Snippet) FieldType(text string) {
	snippet.push(Token{text, kindFieldType})
}

func (snippet *Snippet) Whitespace(text string) {
	snippet.push(Token{text, kindWhitespace})
}
