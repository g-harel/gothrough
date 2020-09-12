package types

const (
	kindComment         = "comment"
	kindEmbeddedName    = "embedded_name"
	kindEmbeddedPackage = "embedded_package"
	kindFieldName       = "field_name"
	kindInterfaceName   = "interface_name"
	kindKeyword         = "keyword"
	kindMethodName      = "method_name"
	kindPunctuation     = "punctuation"
	kindFieldType       = "field_type"
	kindWhitespace      = "whitespace"
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

// TODO make this a method that works on all types instead of attaching it.
type Prettier interface {
	PrettyTokens() []Token
	Pretty() string
}
