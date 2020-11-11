package format

import (
	"fmt"
	"reflect"

	"github.com/g-harel/gothrough/internal/types"
)

// String prints a value to a string.
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

// Format produces a snippet representation of the value.
func Format(value interface{}) (*Snippet, error) {
	if v, ok := value.(*types.Docs); ok {
		return formatDocs(v), nil
	}
	if v, ok := value.(*types.Field); ok {
		return formatField(v), nil
	}
	if v, ok := value.(*types.Function); ok {
		return formatFunction(v, true), nil
	}
	if v, ok := value.(*types.Interface); ok {
		return formatInterface(v, true), nil
	}
	if v, ok := value.(*types.Value); ok {
		return formatValue(v), nil
	}
	return nil, fmt.Errorf("No matching tokenizer: %v", reflect.TypeOf(value))
}
