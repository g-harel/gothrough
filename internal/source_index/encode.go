package source_index

import (
	"encoding/gob"
	"io"

	"github.com/g-harel/gothrough/internal/string_index"
	"github.com/g-harel/gothrough/internal/types"
)

type encodableSearchIndex struct {
	TextIndex  *string_index.Index
	Interfaces []*types.Interface
}

func (si *Index) ToBytes(w io.Writer) error {
	esi := encodableSearchIndex{
		TextIndex:  si.textIndex,
		Interfaces: si.interfaces,
	}

	enc := gob.NewEncoder(w)
	err := enc.Encode(esi)
	if err != nil {
		return err
	}

	return nil
}

func NewIndexFromBytes(r io.Reader) (*Index, error) {
	var esi encodableSearchIndex

	dec := gob.NewDecoder(r)
	err := dec.Decode(&esi)
	if err != nil {
		return nil, err
	}

	si := &Index{
		textIndex: esi.TextIndex,
		interfaces: esi.Interfaces,
	}

	return si, nil
}
