package net

import (
	"net"
)

type Connector interface {
	SendMessage(msg string) error
	SendData(data []byte) error
}

type Service interface {
}

func NewService(serviceName, config string) {

}
