// TODO rename to format
package pretty

import (
	"fmt"
	"reflect"

	"github.com/g-harel/gothrough/internal/types"
)

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

func Pretty(value interface{}) (string, error) {
	tokens, err := PrettyTokens(value)
	if err != nil {
		return "", err
	}

	output := ""
	for _, token := range tokens {
		output += token.Text
	}
	return output, nil
}

func PrettyTokens(value interface{}) ([]Token, error) {
	if ifc, ok := value.(*types.Interface); ok {
		return prettyTokensInterface(ifc), nil
	}
	return []Token{}, fmt.Errorf("No matching tokenizer: %v", reflect.TypeOf(value))
}
