package filter

import "strings"

type Filter struct {
	Query          string
	PackageFilters []string
	Extra          map[string][]string
}

func Parse(query string) Filter {
	parts := strings.Fields(query)

	filter := Filter{
		Query:          "",
		PackageFilters: []string{},
		Extra:          map[string][]string{},
	}
	for _, part := range parts {
		if !strings.Contains(part, ":") {
			filter.Query += " " + part
			continue
		}

		filterQuery := strings.SplitN(part, ":", 2)
		prefix := filterQuery[0]
		query := filterQuery[1]
		if prefix == "package" {
			filter.PackageFilters = append(filter.PackageFilters, query)
		} else {
			if _, ok := filter.Extra[prefix]; !ok {
				filter.Extra[prefix] = []string{}
			}
			filter.Extra[prefix] = append(filter.Extra[prefix], query)
		}
	}

	return filter
}
