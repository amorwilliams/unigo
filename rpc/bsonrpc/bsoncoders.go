package rpc

import (
	l4g "code.google.com/p/log4go"
	"fmt"
	"io"
	"labix.org/v2/mgo/bson"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v interface{}) (err error) {
	buf, err := bson.Marshal(v)
	if err != nil {
		return
	}

	n, err := e.w.Write(buf)
	if err != nil {
		return
	}

	if l := len(buf); n != l {
		err = fmt.Errorf("Wrote %d bytes, should have wrote %d", n, l)
	}

	l4g.Fine(fmt.Sprintf("RPC Wrote %d bytes of %d to connection from buffer: ", n, len(buf)), buf)
	return
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(pv interface{}) (err error) {
	var lbuf [4]byte
	n, err := d.r.Read(lbuf[:])

	if n != 4 {
		err = fmt.Errorf("Corrupted BSON stream: could only read %d", n)
		return
	}

	l4g.Fine(fmt.Sprintf("Read %d bytes of 4 byte length from connection, received bytes: ", n), lbuf)

	if err != nil {
		return
	}

	length := (int(lbuf[0]) << 0) |
		(int(lbuf[1]) << 8) |
		(int(lbuf[2]) << 16) |
		(int(lbuf[3]) << 24)

	l4g.Fine("Message length parsed as: ", length)

	buf := make([]byte, length)
	copy(buf[0:4], lbuf[:])

	n, err = io.ReadFull(d.r, buf[4:0])
	l4g.Fine(fmt.Sprintf("Read %d bytes of %d from connection, received bytes: ", n+4, length), buf)

	if err != nil {
		return
	}

	if n+4 != length {
		err = fmt.Errorf("Expected %d bytes, read %d", length, n)
		return
	}

	if pv != nil {
		err = bson.Unmarshal(buf, pv)
	}

	return
}
