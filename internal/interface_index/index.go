package interface_index

import (
	"sort"
	"strings"

	"github.com/g-harel/gothrough/internal/camel"
	"github.com/g-harel/gothrough/internal/string_index"
	"github.com/g-harel/gothrough/internal/types"
)

// Confidence values for interface info items.
const (
	perfectMatchVal            = 120
	interfaceNameVal           = 120
	totalInterfaceNameTokenVal = 120
	packageNameVal             = 120
	sourceFileVal              = 10
	totalImportPathPartVal     = 20
	totalEmbeddedNameVal       = 80
	totalEmbeddedNameTokenVal  = 80
	totalMethodNameVal         = 80
	totalMethodNameTokenVal    = 80
)

type Index struct {
	index             *string_index.Index
	interfaces        []*types.Interface
	computed_packages *[][]string
}

type Result struct {
	Interface  *types.Interface
	Confidence float64
}

func NewIndex() *Index {
	return &Index{
		index:      string_index.NewIndex(),
		interfaces: []*types.Interface{},
	}
}

// Search returns a interfaces that match the query in deacreasing order of confidence.
func (si *Index) Search(query string) ([]*Result, error) {
	searchResult := si.index.Search(query)
	if len(searchResult) == 0 {
		return []*Result{}, nil
	}

	maxConfidence := searchResult[0].Confidence
	results := make([]*Result, len(searchResult))
	for i, result := range searchResult {
		results[i] = &Result{
			Interface:  si.interfaces[result.ID],
			Confidence: result.Confidence / maxConfidence,
		}
	}

	return results, nil
}

func (si *Index) Insert(ifc types.Interface) {
	si.interfaces = append(si.interfaces, &ifc)
	id := len(si.interfaces) - 1

	// Index on interface name.
	si.index.Insert(id, interfaceNameVal, ifc.Name)
	nameTokens := camel.Split(ifc.Name)
	if len(nameTokens) > 1 {
		si.index.Insert(id, totalInterfaceNameTokenVal/len(nameTokens), nameTokens...)
	}

	// Index on package path and source file.
	importPathParts := strings.Split(ifc.PackageImportPath, "/")
	si.index.Insert(id, packageNameVal, ifc.PackageName)
	si.index.Insert(id, sourceFileVal, strings.TrimSuffix(ifc.SourceFile, ".go"))
	if len(importPathParts) > 1 {
		si.index.Insert(id, totalImportPathPartVal/len(importPathParts), importPathParts...)
	}

	// Index on embedded interfaces.
	if len(ifc.Embedded) > 0 {
		si.index.Insert(id, totalEmbeddedNameVal/len(ifc.Embedded), ifc.Embedded...)
		embeddedNameTokens := []string{}
		for _, embedded := range ifc.Embedded {
			parts := strings.Split(embedded, ".")
			if len(parts) > 1 {
				packageNameParts := strings.Split(parts[0], "_")
				embeddedNameTokens = append(embeddedNameTokens, packageNameParts...)
				embeddedNameTokens = append(embeddedNameTokens, camel.Split(parts[1])...)
			} else {
				embeddedNameTokens = append(embeddedNameTokens, camel.Split(embedded)...)
			}
		}
		if len(embeddedNameTokens) > 1 {
			si.index.Insert(id, totalEmbeddedNameTokenVal/len(embeddedNameTokens), embeddedNameTokens...)
		}
	}

	// Index on interface methods.
	if len(ifc.Methods) > 0 {
		for _, method := range ifc.Methods {
			si.index.Insert(id, totalMethodNameVal/len(ifc.Methods), method.Name)
		}
		methodNameTokens := []string{}
		for _, method := range ifc.Methods {
			methodNameTokens = append(methodNameTokens, camel.Split(method.Name)...)
		}
		if len(methodNameTokens) > 1 {
			si.index.Insert(id, totalMethodNameTokenVal/len(methodNameTokens), methodNameTokens...)
		}
	}
}

func (si *Index) Packages() [][]string {
	if si.computed_packages != nil {
		return *si.computed_packages
	}

	// Collect list of unique packages, separating the standard library vs. hosted ones.
	seenPackages := map[string]bool{}
	stdPackages := []string{}
	hostedPackages := map[string][]string{}
	for _, ifc := range si.interfaces {
		if seenPackages[ifc.PackageImportPath] {
			continue
		}
		seenPackages[ifc.PackageImportPath] = true

		firstSection := strings.Split(ifc.PackageImportPath, "/")[0]
		if !strings.Contains(firstSection, ".") {
			stdPackages = append(stdPackages, ifc.PackageImportPath)
			continue
		}
		if _, ok := hostedPackages[firstSection]; !ok {
			hostedPackages[firstSection] = []string{}
		}
		hostedPackages[firstSection] = append(hostedPackages[firstSection], ifc.PackageImportPath)
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
