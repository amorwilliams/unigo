package session

import (
	"unigo/service"
)

type Session struct {
	ClientInfo string
	Conn       *service.Conn
}
