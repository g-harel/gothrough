package types

func TypeString(t Type) (string, bool) {
	switch v := t.(type) {
	case *Value:
		if v.Const {
			return "const", true
		}
		return "var", true
	case *Interface:
		return "interface", true
	case *Function:
		return "function", true
	default:
		return "", false
	}
}

func Compare(a, b Type) bool {
	aVal := extractTypeSortValue(a)
	bVal := extractTypeSortValue(b)
	if aVal != bVal {
		return aVal < bVal
	}

	aName, aOk := extractName(a)
	bName, bOk := extractName(b)
	if aOk && !bOk {
		return true
	}
	if !aOk && bOk {
		return false
	}
	return aName < bName
}

func extractName(t Type) (string, bool) {
	switch v := t.(type) {
	case *Value:
		return v.Name, true
	case *Interface:
		return v.Name, true
	case *Function:
		return v.Name, true
	default:
		return "", false
	}
}

func extractTypeSortValue(t Type) int {
	switch v := t.(type) {
	case *Value:
		if v.Const {
			return 0
		}
		return 1
	case *Interface:
		return 2
	case *Function:
		return 3
	default:
		return 999
	}
}
