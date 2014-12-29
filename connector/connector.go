package Connector

import (
	"errors"
	// log "github.com/cihub/seelog"
)

var (
	ErrSocketExist = errors.New("Connector: Socket is exist.")
)

var (
	instance Socket
)

type Socket interface {
	Start() error
	Stop() error
	Accept() chan Connector
}

type NewSocketProxy func(host, port, config string) (Socket, error)

func NewSocket(host, port, config string, fc NewSocketProxy) (socket Socket, err error) {
	if instance != nil {
		return nil, ErrSocketExist
	}

	socket, err = fc(host, port, config)
	if err != nil {
		return nil, err
	}

	instance = socket

	return
}

func GetInstance() (socket Socket) {
	return instance
}

type Connector interface {
	Recive() chan []byte
	Send(data []byte) error
	Close()
	Done() chan bool
}
