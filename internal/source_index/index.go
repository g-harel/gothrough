// TODO rename to type_index
package source_index

import (
	"sort"
	"strings"

	"github.com/g-harel/gothrough/internal/string_index"
	"github.com/g-harel/gothrough/internal/types"
)

type Index struct {
	textIndex         *string_index.Index
	interfaces        []*types.Interface
	computed_packages *[][]string
}

type Result struct {
	Confidence        float64
	Name              string
	PackageName       string
	PackageImportPath string
	Value             types.Type
}

func NewIndex() *Index {
	return &Index{
		textIndex:  string_index.NewIndex(),
		interfaces: []*types.Interface{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *Index) Search(query string) ([]*Result, error) {
	searchResult := si.textIndex.Search(query)
	if len(searchResult) == 0 {
		return []*Result{}, nil
	}

	results := make([]*Result, len(searchResult))
	for i, result := range searchResult {
		ifc := si.interfaces[result.ID]
		results[i] = &Result{
			Confidence:        result.Confidence,
			Name:              ifc.Name,
			PackageName:       ifc.PackageName,
			PackageImportPath: ifc.PackageImportPath,
			Value:             ifc,
		}
	}

	return results, nil
}

func (si *Index) Packages() [][]string {
	if si.computed_packages != nil {
		return *si.computed_packages
	}

	// Collect list of unique packages, separating the standard library vs. hosted ones.
	seenPackages := map[string]bool{}
	stdPackages := []string{}
	hostedPackages := map[string][]string{}
	addPackage := func(packageName string) {
		if seenPackages[packageName] {
			return
		}
		seenPackages[packageName] = true

		firstNamePart := strings.Split(packageName, "/")[0]
		if !strings.Contains(firstNamePart, ".") {
			stdPackages = append(stdPackages, packageName)
			return
		}
		if _, ok := hostedPackages[firstNamePart]; !ok {
			hostedPackages[firstNamePart] = []string{}
		}
		hostedPackages[firstNamePart] = append(hostedPackages[firstNamePart], packageName)
	}

	// Add interface package names.
	for _, ifc := range si.interfaces {
		addPackage(ifc.PackageImportPath)
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

	si.computed_packages = &packages
	return packages
}
