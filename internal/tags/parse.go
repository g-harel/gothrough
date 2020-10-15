package tags

import "strings"

type ParsedQuery struct {
	Words string
	Tags  map[string][]string
}

func Parse(query string) ParsedQuery {
	parts := strings.Fields(query)

	parsed := ParsedQuery{
		Words: "",
		Tags:  map[string][]string{},
	}
	for _, part := range parts {
		if !strings.Contains(part, ":") {
			parsed.Words += " " + part
			continue
		}

		tagQuery := strings.SplitN(part, ":", 2)
		prefix := tagQuery[0]
		query := tagQuery[1]

		// Only add the prefix/query combination if it is not a duplicate.
		isNew := true
		for _, existingQuery := range parsed.Tags[prefix] {
			if query == existingQuery {
				isNew = false
				break
			}
		}
		if isNew {
			parsed.Tags[prefix] = append(parsed.Tags[prefix], query)
		}
	}

	return parsed
}
