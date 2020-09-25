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
	textIndex *stringindex.Index
	results   []*Result

	computed_packages         *[][]string
	computed_package_contents *map[string][]int
}

func NewIndex() *Index {
	return &Index{
		textIndex: stringindex.NewIndex(),
		results:   []*Result{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (idx *Index) Search(query string) ([]*Result, error) {
	filters := filter.Parse(query)
	// TODO error if unknown filters

	matches := idx.textIndex.Search(filters.Query)
	if len(matches) == 0 {
		return []*Result{}, nil
	}

	results := make([]*Result, len(matches))
	for i, match := range matches {
		result := idx.results[match.ID]
		result.Confidence = match.Confidence
		results[i] = result
	}

	// Use all results when no query terms.
	if filters.Query == "" {
		results = idx.results
	}

	// TODO use filters

	return results, nil
}

func (idx *Index) Package(importPath string) []*Result {
	if idx.computed_package_contents == nil {
		idx.computed_package_contents = &map[string][]int{}
	}

	resultIDs := []int{}
	if contents, ok := (*idx.computed_package_contents)[importPath]; ok {
		resultIDs = contents
	} else {
		resultIDs = []int{}
		for i, result := range idx.results {
			if result.PackageImportPath == importPath {
				resultIDs = append(resultIDs, i)
			}
		}
		(*idx.computed_package_contents)[importPath] = resultIDs
	}

	results := make([]*Result, len(resultIDs))
	for i, id := range resultIDs {
		results[i] = idx.results[id]
	}

	// TODO sort results
	return results
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
