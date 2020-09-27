package typeindex

import (
	"sort"
	"strings"

	"github.com/g-harel/gothrough/internal/filter"
	"github.com/g-harel/gothrough/internal/stringindex"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for info items.
const (
	confidenceHigh = 120
	confidenceMed  = 80
	confidenceLow  = 20
)

type Result struct {
	Confidence        float64
	Name              string
	PackageName       string
	PackageImportPath string
	Value             types.Type
}

type Index struct {
	textIndex         *stringindex.Index
	results           []*Result
	computed_packages *[][]string
}

func NewIndex() *Index {
	return &Index{
		textIndex: stringindex.NewIndex(),
		results:   []*Result{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (idx *Index) Search(query string) ([]*Result, error) {
	parsedQuery := filter.Parse(query)

	var results []*Result
	if parsedQuery.QueryWords == "" {
		// Use all results when no query terms.
		results = idx.results
	} else {
		matches := idx.textIndex.Search(parsedQuery.QueryWords)
		results = make([]*Result, len(matches))
		for i, match := range matches {
			result := idx.results[match.ID]
			result.Confidence = match.Confidence
			results[i] = result
		}
	}

	// Apply package filter.
	filteredResults := []*Result{}
	if len(parsedQuery.Filters["package"]) > 0 {
		for _, result := range results {
			for _, filterValue := range parsedQuery.Filters["package"] {
				if result.PackageName == filterValue ||
					result.PackageImportPath == filterValue {
					filteredResults = append(filteredResults, result)
				}
			}
		}
	} else {
		filteredResults = results
	}

	return filteredResults, nil
}

func (idx *Index) Packages() [][]string {
	if idx.computed_packages != nil {
		return *idx.computed_packages
	}

	// Collect list of unique packages, separating the standard library vs. hosted ones.
	seenPackages := map[string]bool{}
	stdPackages := []string{}
	hostedPackages := map[string][]string{}

	// Add package names.
	for _, result := range idx.results {
		packageName := result.PackageImportPath

		if seenPackages[packageName] {
			continue
		}
		seenPackages[packageName] = true

		firstNamePart := strings.Split(packageName, "/")[0]
		if !strings.Contains(firstNamePart, ".") {
			stdPackages = append(stdPackages, packageName)
			continue
		}
		if _, ok := hostedPackages[firstNamePart]; !ok {
			hostedPackages[firstNamePart] = []string{}
		}
		hostedPackages[firstNamePart] = append(hostedPackages[firstNamePart], packageName)
	}

	// Create sorted list of hosts.
	hosts := []string{}
	for host := range hostedPackages {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)

	// Created nested array of packages grouped by host and in sorted host order.
	// Standard library packages are added to the front.
	packages := [][]string{stdPackages}
	for _, host := range hosts {
		packages = append(packages, hostedPackages[host])
	}

	// Sort packages within each host's list.
	for i := range packages {
		sort.Strings(packages[i])
	}

	idx.computed_packages = &packages
	return packages
}
