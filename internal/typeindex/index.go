package typeindex

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/stringindex"
	"github.com/g-harel/gothrough/internal/tags"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for info items.
const (
	confidenceHigh = 120
	confidenceMed  = 80
	confidenceLow  = 20
)

type Result struct {
	Confidence float64
	Name       string
	Location   extract.Location
	Value      types.Type
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
	q := tags.Parse(query)

	var results []*Result
	// Use all results when no query terms.
	if q.Words == "" {
		results = idx.results
	} else {
		matches := idx.textIndex.Search(q.Words)
		results = make([]*Result, len(matches))
		for i, match := range matches {
			result := idx.results[match.ID]
			result.Confidence = match.Confidence
			results[i] = result
		}
	}

	// Apply type filter.
	if len(q.Tags["type"]) > 0 {
		temp := results[:]
		results = []*Result{}
		for _, result := range temp {
			for _, tag := range q.Tags["type"] {
				typeString, ok := types.TypeString(result.Value)
				if ok && tag == typeString {
					results = append(results, result)
				}
			}
		}
	}

	// Apply package filter.
	if len(q.Tags["package"]) > 0 {
		temp := results[:]
		results = []*Result{}
		for _, result := range temp {
			for _, tag := range q.Tags["package"] {
				if result.Location.PackageName == tag || result.Location.PackageImportPath == tag {
					results = append(results, result)
				}
			}
		}
	}

	// Sort results when there is no query.
	// This happens when the query has no words.
	if q.Words == "" {
		sort.SliceStable(results, func(i, j int) bool {
			return types.Compare(results[i].Value, results[j].Value)
		})
	}

	// Default to 32 results or use configured number.
	count := 32
	if len(q.Tags["count"]) == 1 {
		c, err := strconv.Atoi(q.Tags["count"][0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse count tag: %v", err)
		} else {
			count = c
		}
	}
	if len(results) > count {
		results = results[:count]
	}

	return results, nil
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
		packageName := result.Location.PackageImportPath

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
