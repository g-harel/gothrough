package gis

import (
	"strings"
)

type mappedValue struct {
	id         string
	index      int
	confidence int
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

func (q *Querier) createMapping(query string, m mappedValue) {
	query = strings.ToLower(query)

	if len(q.mappings[query]) == 0 {
		q.mappings[query] = []mappedValue{m}
	}

	for i, currentMapping := range q.mappings[query] {
		if currentMapping.confidence < m.confidence {
			q.mappings[query] = make([]mappedValue, len(q.mappings[query])+1)
			copy(q.mappings[query][:i], q.mappings[query][:i])
			copy(q.mappings[query][i+1:], q.mappings[query][i:])
			q.mappings[query][i] = m
			break
		}
	}

	q.mappings[query] = append(q.mappings[query], m)
}

func (q *Querier) Write(i Interface) {
	index := len(q.values)
	q.values = append(q.values, &i)

	mapping := mappedValue{
		id:    i.Address(),
		index: index,
	}

	mapping.confidence = 10
	q.createMapping(i.Name, mapping)

	mapping.confidence = 7
	q.createMapping(i.PackageName, mapping)

	mapping.confidence = 5
	q.createMapping(strings.TrimSuffix(i.SourceFile, ".go"), mapping)

	for _, method := range i.Methods {
		mapping.confidence = 5
		q.createMapping(method, mapping)
	}

	for _, part := range strings.Split(i.PackageImportPath, "/") {
		if part == i.PackageName {
			continue
		}
		mapping.confidence = 3
		q.createMapping(part, mapping)
	}
}

func (q *Querier) Query(query string) []*Interface {
	query = strings.ToLower(query)

	// TODO multi-toke search (combine confidences)
	res := []*Interface{}
	for _, m := range q.mappings[query] {
		res = append(res, q.values[m.index])
	}
	return res
}
