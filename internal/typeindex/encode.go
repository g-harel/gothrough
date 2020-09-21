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

type encodableTypeIndex struct {
	TextIndex *stringindex.Index
	Results   []*Result
}

func (idx *Index) ToBytes(w io.Writer) error {
	edx := encodableTypeIndex{
		TextIndex: idx.textIndex,
		Results:   idx.results,
	}

	enc := gob.NewEncoder(w)
	err := enc.Encode(edx)
	if err != nil {
		return err
	}

	return nil
}

func NewIndexFromBytes(r io.Reader) (*Index, error) {
	var edx encodableTypeIndex

	dec := gob.NewDecoder(r)
	err := dec.Decode(&edx)
	if err != nil {
		return nil, err
	}

	idx := &Index{
		textIndex: edx.TextIndex,
		results:   edx.Results,
	}

	return idx, nil
}
