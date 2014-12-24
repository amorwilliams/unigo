package bsonrpc

import (
	"io"
	"net/rpc"
)

type ServerCodec struct {
	conn    io.ReadWriteCloser
	Encoder *Encoder
	Decoder *Decoder
}

func NewServerCodec(conn io.ReadWriteCloser) (codec *ServerCodec) {
	codec = &ServerCodec{
		conn:    conn,
		Encoder: NewEncoder(conn),
		Decoder: NewDecoder(conn),
	}
	return
}

func (sc *ServerCodec) ReadRequestHeader(rq *rpc.Request) (err error) {
	err = sc.Decoder.Decode(rq)
	return
}

func (sc *ServerCodec) ReadRequestBody(v interface{}) (err error) {
	err = sc.Decoder.Decode(v)
	return
}

func (sc *ServerCodec) WriteResponse(rs *rpc.Response, v interface{}) (err error) {
	err = sc.Encoder.Encode(rs)
	if err != nil {
		return
	}
	err = sc.Encoder.Encode(v)
	if err != nil {
		return
	}
	return
}

func (sc *ServerCodec) Close() (err error) {
	err = sc.conn.Close()
	return
}

func ServeConn(conn io.ReadWriteCloser) (s *rpc.Server) {
	s = rpc.NewServer()
	s.ServeCodec(NewServerCodec(conn))
	return
}
