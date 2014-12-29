package session

import (
	"github.com/amorwilliams/unigo/connector"
	log "github.com/cihub/seelog"
)

type SessionService struct {
	sessions map[string]Session
	uidMap   map[string]Session
}

func NewSessionService() *SessionService {
	return &SessionService{
		sessions: make(map[string]Session),
		uidMap:   make(map[string]Session),
	}
}

func (ss *SessionService) CreateSession(sid, frontendId string, conn Connector.Connector) *Session {
	sess := NewSession(sid, frontendId, conn, ss)
	ss.sessions[sess.id] = sess

	return sess
}

func (ss *SessionService) Bind(sid, uid string) {

}

func (ss *SessionService) Unbind(sid, uid string) {

}

// func (ss *SessionService) NewSession(sid) *Session {
// return &
// }

type SessionState int8

const (
	ST_INITED SessionState = iota
	ST_CLOSED
)

type Session struct {
	sid        string
	frontendId string
	uid        string
	settings   map[string]interface{}
	conn       Connector.Connector
	service    *SessionService
}

func NewSession(sid, frontendId string, conn Connector.Connector, service *SessionService) *Session {
	return &Session{
		sid:        sid,
		frontendId: frontendId,
		uid:        "",
		settings:   make(map[string]interface{}),
		conn:       conn,
		service:    service,
	}
}

func (s *Session) Bind(uid string) {
	s.uid = uid
}

func (s *Session) Unbind(uid string) {
	s.uid = ""
}

func (s *Session) Set(k string, v interface{}) {
	s.settings[k] = v
}

func (s *Session) Get(k string) interface{} {
	return s.settings[k]
}

func (s *Session) Send(data []byte) error {
	return s.conn.Send(data)
}

func (s *Session) Close(reason string) {

}
