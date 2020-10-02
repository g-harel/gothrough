package format

import (
	"fmt"
	"reflect"

	"github.com/g-harel/gothrough/internal/tokens"
	"github.com/g-harel/gothrough/internal/types"
)

func String(value types.Type) (string, error) {
	snippet, err := Format(value)
	if err != nil {
		return "", err
	}

	output := ""
	for _, token := range snippet.Dump() {
		output += token.Text
	}
	return output, nil
}

func Format(value interface{}) (*tokens.Snippet, error) {
	if v, ok := value.(*types.Docs); ok {
		return formatDocs(v), nil
	}
	if v, ok := value.(*types.Field); ok {
		return formatField(v), nil
	}
	if v, ok := value.(*types.Function); ok {
		return formatFunction(v), nil
	}
	if v, ok := value.(*types.Interface); ok {
		return formatInterface(v), nil
	}
	return nil, fmt.Errorf("No matching tokenizer: %v", reflect.TypeOf(value))
}
