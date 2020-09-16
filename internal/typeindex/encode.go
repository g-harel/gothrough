package typeindex

import (
	"encoding/gob"
	"io"

	"github.com/g-harel/gothrough/internal/stringindex"
	"github.com/g-harel/gothrough/internal/types"
)

// Register types that could be hidden behind "types.Type".
func init() {
	gob.Register(&types.Interface{})
}

type encodableSearchIndex struct {
	TextIndex *stringindex.Index
	Results   []*Result
}

func (si *Index) ToBytes(w io.Writer) error {
	esi := encodableSearchIndex{
		TextIndex: si.textIndex,
		Results:   si.results,
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
		results:   esi.Results,
	}

	return si, nil
}
