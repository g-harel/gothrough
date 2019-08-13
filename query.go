package gis

import (
	"sort"
	"strings"

	"github.com/g-harel/gis/internal/camel"
	"github.com/g-harel/gis/internal/interfaces"
)

type mappedValue struct {
	addr       string
	index      int
	confidence float32
}

type Querier struct {
	values   []*interfaces.Interface
	mappings map[string][]mappedValue
}

func NewQuerier() *Querier {
	return &Querier{
		values:   []*interfaces.Interface{},
		mappings: map[string][]mappedValue{},
	}
}

func (q *Querier) createMappings(m mappedValue, queries ...string) {
	for _, query := range queries {
		query = strings.ToLower(query)
		if len(q.mappings[query]) == 0 {
			q.mappings[query] = []mappedValue{m}
		}
		q.mappings[query] = append(q.mappings[query], m)
	}
}

func (q *Querier) Write(i *interfaces.Interface) {
	index := len(q.values)
	q.values = append(q.values, i)

	mapping := mappedValue{
		addr:  i.Address(),
		index: index,
	}

	// interfaces.Interface name.
	mapping.confidence = 10
	q.createMappings(mapping, i.Name)

	// interfaces.Interface name tokens.
	tokens := camel.Split(i.Name)
	if len(tokens) > 1 {
		mapping.confidence = 6 / float32(len(tokens))
		q.createMappings(mapping, tokens...)
	}

	// Package name.
	mapping.confidence = 7
	q.createMappings(mapping, i.PackageName)

	// Source file name.
	mapping.confidence = 2
	q.createMappings(mapping, strings.TrimSuffix(i.SourceFile, ".go"))

	// Method names.
	mapping.confidence = 5 / float32(len(i.Methods))
	q.createMappings(mapping, i.Methods...)

	// Method name tokens.
	for _, methodName := range i.Methods {
		tokens := camel.Split(methodName)
		if len(tokens) > 1 {
			mapping.confidence = 4 / float32(len(i.Methods)) / float32(len(tokens))
			q.createMappings(mapping, tokens...)
		}
	}

	// Import path parts.
	mapping.confidence = 3
	q.createMappings(mapping, strings.Split(i.PackageImportPath, "/")...)
}

func (q *Querier) Query(query string) []*interfaces.Interface {
	query = strings.ToLower(query)

	indexes := map[string]int{}
	confidences := map[string]float32{}

	// Sum confidences for all matches.
	for _, subQuery := range strings.Fields(query) {
		for _, m := range q.mappings[subQuery] {
			indexes[m.addr] = m.index
			if _, ok := confidences[m.addr]; !ok {
				confidences[m.addr] = 0
			}
			confidences[m.addr] += m.confidence
		}
	}

	// Combine results into list without duplicates.
	res := make([]*interfaces.Interface, len(indexes))
	i := 0
	for _, index := range indexes {
		res[i] = q.values[index]
		i++
	}

	// Sort return value by decreasing confidence.
	sort.Slice(res, func(i, j int) bool {
		return confidences[res[i].Address()] > confidences[res[j].Address()]
	})

	return res
}
