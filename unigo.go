package unigo

import (
	"github.com/amorwilliams/unigo/service"
	"github.com/amorwilliams/unigo/utils"
	// . "unigo/session"
)

const VERSION = "0.0.1"

type UniDefaultServer struct{}

func (s *UniDefaultServer) Setup()        {}
func (s *UniDefaultServer) Started()      {}
func (s *UniDefaultServer) Stopped()      {}
func (s *UniDefaultServer) Registered()   {}
func (s *UniDefaultServer) Unregistered() {}
func (s *UniDefaultServer) ConnectNew(conn service.Connector) (uuid string, err error) {
	return utils.NewUUID(), nil
}
func (s *UniDefaultServer) ConnectClose(conn service.Connector, uuid string)                   {}
func (s *UniDefaultServer) NewMessageReceived(conn service.Connector, msg string, uuid string) {}
func (s *UniDefaultServer) NewDataReceived(conn service.Connector, data []byte, uuid string)   {}

func NewSimpleServer() (s *service.Service) {
	sd := &UniDefaultServer{}

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
