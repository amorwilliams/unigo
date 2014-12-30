package utils

import (
	"bytes"
	"encoding/gob"
)

func DeepCopy(src, dst interface{}) error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(src); err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewBuffer(buf.Bytes()))
	return dec.Decode(dst)
}
