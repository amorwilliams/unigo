package unigo

import (
	"github.com/amorwilliams/unigo/service"
	// . "unigo/session"
)

const VERSION = "0.0.1"

type UniServer struct{}

func (s *UniServer) Setup()                                                       {}
func (s *UniServer) Started()                                                     {}
func (s *UniServer) Stopped()                                                     {}
func (s *UniServer) Registered()                                                  {}
func (s *UniServer) Unregistered()                                                {}
func (s *UniServer) CreatePeer(conn service.Connector) (service.Connector, error) { return conn, nil }
func (s *UniServer) NewMessageReceived(conn service.Connector, msg string)        {}
func (s *UniServer) NewDataReceived(conn service.Connector, data []byte)          {}

func NewSimpleServer() (s *service.Service) {
	sd := &UniServer{}

	si := &service.ServiceInfo{}
	si.Name = "Default Server"
	si.Version = "1"

	s = service.CreateService(sd, si)
	return
}

func CreateServer(sd service.ServiceDelegate, si *service.ServiceInfo) (s *service.Service) {
	s = service.CreateService(sd, si)
	return
}
