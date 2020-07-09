package gis

import (
	"encoding/gob"
	"io"

	"github.com/g-harel/gis/internal/index"
	"github.com/g-harel/gis/internal/interfaces"
)

type encodableSearchIndex struct {
	Index      *index.Index
	Interfaces []*interfaces.Interface
}

func (si *SearchIndex) ToBytes(w io.Writer) error {
	esi := encodableSearchIndex{
		Index:      si.index,
		Interfaces: si.interfaces,
	}

	enc := gob.NewEncoder(w)
	err := enc.Encode(esi)
	if err != nil {
		return err
	}

	return nil
}

func NewSearchIndexFromBytes(r io.Reader) (*SearchIndex, error) {
	var esi encodableSearchIndex

	dec := gob.NewDecoder(r)
	err := dec.Decode(&esi)
	if err != nil {
		return nil, err
	}

	si := &SearchIndex{
		index:      esi.Index,
		interfaces: esi.Interfaces,
	}

	return si, nil
}
