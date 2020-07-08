package index

import (
	"bytes"
	"encoding/gob"
)

var _ gob.GobEncoder = &Index{}
var _ gob.GobDecoder = &Index{}

func init() {
	gob.Register(Index{})
}

func (idx *Index) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	err := encoder.Encode(idx.mappings)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (idx *Index) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	err := decoder.Decode(&idx.mappings)
	if err != nil {
		return err
	}

	return nil
}
