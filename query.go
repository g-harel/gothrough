package gis

import (
	"sort"
	"strings"
)

type mappedValue struct {
	addr       string
	index      int
	confidence float32
}

type Querier struct {
	values   []*Interface
	mappings map[string][]mappedValue
}

func NewQuerier() *Querier {
	return &Querier{
		values:   []*Interface{},
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

func (q *Querier) Write(i Interface) {
	// TODO break identifiers by camelCase

	index := len(q.values)
	q.values = append(q.values, &i)

	mapping := mappedValue{
		addr:  i.Address(),
		index: index,
	}

	mapping.confidence = 10
	q.createMappings(mapping, i.Name)

	tokens := CamelSplit(i.Name)
	mapping.confidence = 6 + 3/float32(len(tokens))
	q.createMappings(mapping, tokens...)

	mapping.confidence = 7
	q.createMappings(mapping, i.PackageName)

	mapping.confidence = 5
	q.createMappings(mapping, strings.TrimSuffix(i.SourceFile, ".go"))

	// Methods in larger interfaces are given lower confidence.
	mapping.confidence = 5 + 2/float32(len(i.Methods))
	q.createMappings(mapping, i.Methods...)

	mapping.confidence = 3
	q.createMappings(mapping, strings.Split(i.PackageImportPath, "/")...)
}

func (q *Querier) Query(query string) []*Interface {
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
	res := make([]*Interface, len(indexes))
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
