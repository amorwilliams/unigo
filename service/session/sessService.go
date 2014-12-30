package session

import (
	"errors"
	"github.com/amorwilliams/unigo/connector"
	log "github.com/cihub/seelog"
)

type SessionService struct {
	sessions map[string]*Session
	uidMap   map[string][]Session
}

func NewSessionService() *SessionService {
	return &SessionService{
		sessions: make(map[string]*Session),
		uidMap:   make(map[string][]Session),
	}
}

func (ss *SessionService) CreateSession(sid, frontendId string, conn Connector.Connector) (sess *Session) {
	sess = NewSession(sid, frontendId, conn, ss)
	ss.sessions[sess.sid] = sess

	return
}

func (ss *SessionService) Bind(sid, uid string) {
	sess, ok := ss.sessions[sid]
	if !ok {
		return
	}

	if sess.uid != ""{
		if sess.uid == uid{
			return
		}
	}
	}
}

func (ss *SessionService) Unbind(sid, uid string) {

}

func (ss *SessionService) Get(sid string) *Session {
	return ss.sessions[sid]
}

func (ss *SessionService) GetByUid(uid string) []Session {
	return ss.uidMap[uid]
}

func (ss *SessionService) Remove(sid string) {
	if sess, ok := ss.sessions[sid]; ok {
		uid := sess.uid
		delete(ss.sessions, sess.uid)

		sessions, ok := ss.uidMap[uid]
		if !ok {
			return
		}
		for i, l := 0, len(sessions); i < l; i++ {
			if sessions[i].uid == sid {
				sessions = sessions[i : i+1]
				if len(sessions) == 0 {
					delete(ss.uidMap, uid)
				}
				break
			}
		}
	}
}

func (ss *SessionService) SendMessage(sid string, data []byte) error {
	sess, ok := ss.sessions[sid]
	if !ok {
		return errors.New("Fail to send message for non-existing session, sid: ", sid, " msg: ", msg)
	}
	return send(ss, sess, data)
}

func (ss *SessionService) SendMessageByUid(uid string, data []byte) error {
	sessions, ok := ss.uidMap[uid]
	if !ok {
		return errors.New(fmt.Sprintf("fail to send message by uid for non-existing session. uid: %s", uid))
	}

	for i, l := 0, len(sessions); i < l; i++ {
		send(ss, sessions[i], data)
	}
}

type SessCallback func(sess *Session)

func (ss *SessionService) ForeachSess(scb SessCallback) {
	for _, sess := range ss.sessions {
		scb(sess)
	}
}

func (ss *SessionService) ForeachBindedSess(scb SessCallback) {
	for _, sessions := range ss.uidMap {
		for i, l := 0, len(sessions); i < l; i++ {
			scb(sessions[i])
		}
	}
}

func (ss *SessionService) GetSessCount() int {
	return len(ss.sessions)
}

func send(service *SessionService, sess *Session, data []byte) error {
	return sess.send(data)
}

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
	state      SessionState
}

func NewSession(sid, frontendId string, conn Connector.Connector, service *SessionService) *Session {
	return &Session{
		sid:        sid,
		frontendId: frontendId,
		uid:        "",
		settings:   make(map[string]interface{}),
		conn:       conn,
		service:    service,
		state:      ST_INITED,
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

func (s *Session) send(data []byte) error {
	return s.conn.Send(data)
}

func (s *Session) Close(reason string) {
	log.Debugf("session on [%s] is closed with session id: %s", s.frontendId, s.sid)
	if s.state == ST_CLOSED {
		return
	}
	s.state = ST_CLOSED

}
