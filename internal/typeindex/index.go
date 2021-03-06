package typeindex

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/g-harel/gothrough/internal/extract"
	"github.com/g-harel/gothrough/internal/stringindex"
	"github.com/g-harel/gothrough/internal/typeindex/tags"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for info items.
const (
	confidenceHigh = 120
	confidenceMed  = 80
	confidenceLow  = 20
)

// Result represents a matched value from the index.
type Result struct {
	Confidence float64
	Name       string
	Location   extract.Location
	Value      types.Type
}

// Index stores types and enables them to be searched.
type Index struct {
	textIndex        *stringindex.Index
	results          []*Result
	computedPackages *[][]string
}

// NewIndex creates a new empty index.
func NewIndex() *Index {
	return &Index{
		textIndex: stringindex.NewIndex(),
		results:   []*Result{},
	}
}

func filter(original []*Result, handler func(result *Result) bool) []*Result {
	filtered := []*Result{}
	for _, result := range original {
		if handler(result) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func and(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}

func or(bools []bool) bool {
	for _, b := range bools {
		if b {
			return true
		}
	}
	return false
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (idx *Index) Search(query string) ([]*Result, error) {
	q := tags.Parse(query)

	plain := q.GetWords() == ""

	var results []*Result
	// Use all results when no query terms.
	if plain {
		results = make([]*Result, len(idx.results))
		copy(results, idx.results)
	} else {
		matches := idx.textIndex.Search(q.GetWords())
		results = make([]*Result, len(matches))
		for i, match := range matches {
			result := *idx.results[match.ID]
			result.Confidence = match.Confidence
			results[i] = &result
		}
	}

	// Apply type filter.
	typeTags := q.GetTags("type", "t")
	if len(typeTags) > 0 {
		results = filter(results, func(result *Result) bool {
			bools := []bool{}
			for _, tag := range typeTags {
				typeString, ok := types.TypeString(result.Value)
				bools = append(bools, ok && tag == typeString)
			}
			return or(bools)
		})
	}

	// Apply inverted type filter.
	invertedTypeTags := q.GetTags("-type", "-t")
	if len(invertedTypeTags) > 0 {
		results = filter(results, func(result *Result) bool {
			bools := []bool{}
			for _, tag := range invertedTypeTags {
				typeString, ok := types.TypeString(result.Value)
				bools = append(bools, !ok || tag != typeString)
			}
			return and(bools)
		})
	}

	// Apply package filter.
	packageTags := q.GetTags("package", "p")
	if len(packageTags) > 0 {
		results = filter(results, func(result *Result) bool {
			bools := []bool{}
			for _, tag := range packageTags {
				bools = append(bools, result.Location.PackageName == tag || result.Location.PackageImportPath == tag)
			}
			return or(bools)
		})
	}

	// Apply inverted package filter.
	invertedPackageTags := q.GetTags("-package", "-p")
	if len(invertedPackageTags) > 0 {
		results = filter(results, func(result *Result) bool {
			bools := []bool{}
			for _, tag := range invertedPackageTags {
				bools = append(bools, result.Location.PackageName != tag && result.Location.PackageImportPath != tag)
			}
			return and(bools)
		})
	}

	// Sort results when plain query.
	if plain {
		sort.SliceStable(results, func(i, j int) bool {
			return types.Compare(results[i].Value, results[j].Value)
		})
	}

	// Default to 32 results or use configured number.
	count := 32
	countTags := q.GetTags("count", "c")
	if len(countTags) == 1 {
		c, err := strconv.Atoi(countTags[0])
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

// Packages returns a list of all packages in the index.
// The first item will contain the standard library packages.
// Subsequent hosted packages are grouped by the host domain.
func (idx *Index) Packages() [][]string {
	if idx.computedPackages != nil {
		return *idx.computedPackages
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

	idx.computedPackages = &packages
	return packages
}
