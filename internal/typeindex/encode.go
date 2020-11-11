package typeindex

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/g-harel/gothrough/internal/stringindex"
	"github.com/g-harel/gothrough/internal/types"
)

// Register types that could be hidden behind "types.Type".
func init() {
	gob.Register(&types.Function{})
	gob.Register(&types.Interface{})
	gob.Register(&types.Value{})
}

type encodableTypeIndex struct {
	TextIndex *stringindex.Index
	Results   []*Result
}

// ToBytes encodes the index to bytes and writes them to the writer.
func (idx *Index) ToBytes(w io.Writer) error {
	edx := encodableTypeIndex{
		TextIndex: idx.textIndex,
		Results:   idx.results,
	}

	gw, err := gzip.NewWriterLevel(w, gzip.NoCompression)
	if err != nil {
		return fmt.Errorf("create gzip writer: %v", err)
	}

	enc := gob.NewEncoder(gw)
	err = enc.Encode(edx)
	if err != nil {
		return fmt.Errorf("encode index: %v", err)
	}

	err = gw.Close()
	if err != nil {
		return fmt.Errorf("close gzip writer: %v", err)
	}

	return nil
}

// NewIndexFromBytes decodes the reader's data into an index.
func NewIndexFromBytes(r io.Reader) (*Index, error) {
	var edx encodableTypeIndex

	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("create gzip reader: %v", err)
	}

	dec := gob.NewDecoder(gr)
	err = dec.Decode(&edx)
	if err != nil {
		return nil, fmt.Errorf("decode index: %v", err)
	}

	err = gr.Close()
	if err != nil {
		return nil, fmt.Errorf("close gzip reader: %v", err)
	}

	idx := &Index{
		textIndex: edx.TextIndex,
		results:   edx.Results,
	}

	return idx, nil
}
