package filter

import "strings"

type ParsedQuery struct {
	QueryWords string
	// TODO rename to tag
	Filters map[string][]string
}

func Parse(query string) ParsedQuery {
	parts := strings.Fields(query)

	parsed := ParsedQuery{
		QueryWords: "",
		Filters:    map[string][]string{},
	}
	for _, part := range parts {
		if !strings.Contains(part, ":") {
			parsed.QueryWords += " " + part
			continue
		}

		filterQuery := strings.SplitN(part, ":", 2)
		prefix := filterQuery[0]
		query := filterQuery[1]

		// Only add the prefix/query combination if it is not a duplicate.
		isNew := true
		for _, existingQuery := range parsed.Filters[prefix] {
			if query == existingQuery {
				isNew = false
				break
			}
		}
		if isNew {
			parsed.Filters[prefix] = append(parsed.Filters[prefix], query)
		}
	}

	return parsed
}
