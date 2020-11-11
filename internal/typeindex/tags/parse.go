package tags

import "strings"

// ParsedQuery represents a parsed query.
type ParsedQuery struct {
	words string
	tags  map[string][]string
}

// Parse parses the query to extract tags.
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

// GetWords returns the parts of the query that were not parsed as tags.
func (p ParsedQuery) GetWords() string {
	return strings.TrimSpace(p.words)
}

// GetTags collects the tag values for all the given tag names.
func (p ParsedQuery) GetTags(tags ...string) []string {
	allValues := []string{}
	for _, tag := range tags {
		allValues = append(allValues, p.tags[tag]...)
	}
	return allValues
}
