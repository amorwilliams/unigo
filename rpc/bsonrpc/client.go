package bsonrpc

import (
	"io"
	"net/rpc"
)

type ClientCodec struct {
	conn    io.ReadWriteCloser
	Encoder *Encoder
	Decoder *Decoder
}

func NewClientCodec(conn io.ReadWriteCloser) (codec *ClientCodec) {
	cc := &ClientCodec{
		conn:    conn,
		Encoder: NewEncoder(conn),
		Decoder: NewDecoder(conn),
	}
	codec = cc
	return
}

func (cc *ClientCodec) WriteRequest(req *rpc.Request, v interface{}) (err error) {
	err = cc.Encoder.Encode(req)
	if err != nil {
		cc.Close()
		return
	}

	err = cc.Encoder.Encode(v)
	if err != nil {
		return
	}

	return
}

func (cc *ClientCodec) ReadResponseHeader(res *rpc.Response) (err error) {
	err = cc.Decoder.Decode(res)
	return
}

func (cc *ClientCodec) ReadResponseBody(v interface{}) (err error) {
	err = cc.Decoder.Decode(v)
	return
}

func (cc *ClientCodec) Close() (err error) {
	err = cc.conn.Close()
	return
}

func NewClient(conn io.ReadWriteCloser) (c *rpc.Client) {
	cc := NewClientCodec(conn)
	c = rpc.NewClientWithCodec(cc)
	return
}
