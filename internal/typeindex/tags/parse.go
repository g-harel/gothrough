package tags

import "strings"

type ParsedQuery struct {
	words string
	tags  map[string][]string
}

func Parse(query string) ParsedQuery {
	parts := strings.Fields(query)

	parsed := ParsedQuery{
		words: "",
		tags:  map[string][]string{},
	}
	for _, part := range parts {
		if !strings.Contains(part, ":") {
			parsed.words += " " + part
			continue
		}

		tagQuery := strings.SplitN(part, ":", 2)
		prefix := tagQuery[0]
		query := tagQuery[1]

		// Only add the prefix/query combination if it is not a duplicate.
		isNew := true
		for _, existingQuery := range parsed.tags[prefix] {
			if query == existingQuery {
				isNew = false
				break
			}
		}
		if isNew {
			parsed.tags[prefix] = append(parsed.tags[prefix], query)
		}
	}

	return parsed
}

func (p ParsedQuery) GetWords() string {
	return strings.TrimSpace(p.words)
}

func (p ParsedQuery) GetTags(tags ...string) []string {
	allValues := []string{}
	for _, tag := range tags {
		allValues = append(allValues, p.tags[tag]...)
	}
	return allValues
}
