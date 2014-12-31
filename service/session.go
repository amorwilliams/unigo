package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/amorwilliams/unigo/connector"
	log "github.com/cihub/seelog"
)

type SessionService struct {
	single   bool
	sessions map[string]*Session
	uidMap   map[string][]Session
	l        sync.RWMutex
}

func NewSessionService() *SessionService {
	return &SessionService{
		single:   false,
		sessions: make(map[string]*Session),
		uidMap:   make(map[string][]Session),
		l:        new(sync.RWMutex),
	}
}

func (ss *SessionService) CreateSession(sid, frontendId string, conn Connector.Connector) (s *Session) {
	s = NewSession(sid, frontendId, conn, ss)
	ss.sessions[s.sid] = s

	return
}

//Bind the session with a user id.
func (ss *SessionService) Bind(sid, uid string) error {
	ss.l.Lock()
	defer ss.l.Unlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return errors.New(fmt.Sprintf("session does not exist, sid: %s", sid))
	}

	if len(s.uid) > 0 {
		if s.uid == uid { // already bound with the same uid
			return nil
		}

		// already bound with other uid
		return errors.New(fmt.Sprintf("session has already bound with %s", s.uid))
	}

	ses, ok := ss.uidMap[uid]
	if !ok && !ss.single {
		return errors.New(fmt.Sprintf("singleSession is enabled, and session has already bound with uid: %s", s.uid))
	}

	if !ok {
		ses = make([]Session, 0)
		ss.uidMap[uid] = ses
	}

	for i, l := 0, len(ses); i < l; i++ {
		if ses[i].uid == s.uid { // session has binded with the uid
			return errors.New(fmt.Sprintf("session has already bound with %s", s.uid))
		}
	}

	ses = append(ses, *s)
	s.bind(uid)
	return nil
}

// Unbind a session with the user id.
func (ss *SessionService) Unbind(sid, uid string) error {
	ss.l.Lock()
	defer ss.l.Unlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return errors.New(fmt.Sprintf("session does not exist, sid: %s", sid))
	}

	if len(s.uid) <= 0 || s.uid != uid {
		return errors.New(fmt.Sprintf("session has already bound with %s", s.uid))
	}

	ses, ok := ss.uidMap[uid]
	if ok {
		for i, l := 0, len(ses); i < l; i++ {
			if ses[i].uid == sid {
				ses = ses[i : i+1]
				break
			}
		}
	}

	if len(ses) == 0 {
		delete(ss.uidMap, uid)
	}

	s.unbind(uid)

	return nil
}

func (ss *SessionService) Get(sid string) (s *Session, ok bool) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	s, ok = ss.sessions[sid]
	return
}

func (ss *SessionService) GetByUid(uid string) (ses []Session, ok bool) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	ses, ok = ss.uidMap[uid]
	return
}

func (ss *SessionService) Remove(sid string) {
	ss.l.Lock()
	defer ss.l.Unlock()

	s, ok := ss.sessions[sid]
	if ok {
		uid := s.uid
		delete(ss.sessions, s.sid)

		ses, ok := ss.uidMap[uid]
		if !ok {
			return
		}

		for i, l := 0, len(ses); i < l; i++ {
			if ses[i].uid == sid {
				ses = ses[i : i+1]
				if len(ses) == 0 {
					delete(ss.uidMap, uid)
				}
				break
			}
		}
	}
}

func (ss *SessionService) Import(sid, k string, v interface{}) error {
	ss.l.Lock()
	defer ss.l.Unlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return errors.New(fmt.Sprintf("session does not exist, sid: %s", sid))
	}

	s.Set(k, v)
	return nil
}

func (ss *SessionService) ImportAll(sid string, settings map[string]interface{}) error {
	ss.l.Lock()
	defer ss.l.Unlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return errors.New(fmt.Sprintf("session does not exist, sid: %s", sid))
	}

	for k, v := range settings {
		s.Set(k, v)
	}
	return nil
}

func (ss *SessionService) Kick(uid, reason string) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	if reason == "" {
		reason = "kick"
	}

	ses, ok := ss.uidMap[uid]
	if !ok {
		return
	}

	sids := make([]string, 0)
	for _, s := range ses {
		sids = append(sids, s.sid)
	}

	for _, sid := range sids {
		ss.sessions[sid].Close(reason)
	}
}

func (ss *SessionService) KickBySessionId(sid string) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return
	}

	s.Close("kick")
}

func (ss *SessionService) SendMessage(sid string, data []byte) error {
	ss.l.RLock()
	defer ss.l.RUnlock()

	s, ok := ss.sessions[sid]
	if !ok {
		return errors.New(fmt.Sprintf("Fail to send message for non-existing session, sid: %s msg: %v", sid, data))
	}
	return send(ss, s, data)
}

func (ss *SessionService) SendMessageByUid(uid string, data []byte) error {
	ss.l.RLock()
	defer ss.l.RUnlock()

	ses, ok := ss.uidMap[uid]
	if !ok {
		return errors.New(fmt.Sprintf("fail to send message by uid for non-existing session. uid: %s", uid))
	}

	var err error
	for _, s := range ses {
		err = send(ss, &s, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ss *SessionService) ForeachSess(scb func(s *Session)) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	for _, s := range ss.sessions {
		scb(s)
	}
}

func (ss *SessionService) ForeachBindedSess(scb func(s *Session)) {
	ss.l.RLock()
	defer ss.l.RUnlock()

	for _, ses := range ss.uidMap {
		for _, s := range ses {
			scb(&s)
		}
	}
}

func (ss *SessionService) GetSessCount() int {
	ss.l.RLock()
	defer ss.l.RUnlock()

	return len(ss.sessions)
}

func send(service *SessionService, s *Session, data []byte) error {
	return s.send(data)
}

// Session State
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
	l          sync.RWMutex
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
		l:          new(sync.RWMutex),
	}
}

func (s *Session) toFrontendSession() *FrontendSession {
	return NewFrontendSession(s)
}

func (s *Session) bind(uid string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.uid = uid
}

func (s *Session) unbind(uid string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.uid = ""
}

func (s *Session) Set(k string, v interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	s.settings[k] = v
}

func (s *Session) Get(k string) (v interface{}, ok bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	v, ok = s.settings[k]
	return
}

func (s *Session) send(data []byte) error {
	return s.conn.Send(data)
}

func (s *Session) Close(reason string) {
	s.l.Lock()
	defer s.l.Unlock()

	log.Debugf("session on [%s] is closed with session id: %s", s.frontendId, s.sid)
	if s.state == ST_CLOSED {
		return
	}

	s.state = ST_CLOSED
	s.service.Remove(s.sid)
	s.conn.Close()
}

type FrontendSession struct {
	sid        string
	frontendId string
	uid        string
	settings   map[string]interface{}
	session    *Session
	service    *SessionService
	l          sync.Mutex
}

func NewFrontendSession(s *Session) *FrontendSession {
	return &FrontendSession{
		sid:        s.sid,
		uid:        s.uid,
		frontendId: s.frontendId,
		settings:   dclone(s.settings),
		session:    s,
		service:    s.service,
	}
}

func (fs *FrontendSession) Bind(uid string) error {
	if err := fs.service.Bind(fs.sid, uid); err != nil {
		return err
	}

	fs.l.Lock()
	defer fs.l.Unlock()

	fs.uid = uid
	return nil
}

func (fs *FrontendSession) Unbind(uid string) error {
	if err := fs.service.Unbind(fs.sid, uid); err != nil {
		return err
	}

	fs.l.Lock()
	defer fs.l.Unlock()

	fs.uid = ""
	return nil
}

func (fs *FrontendSession) Get(k string) (v interface{}, ok bool) {
	fs.l.Lock()
	defer fs.l.Unlock()

	v, ok = fs.settings[k]
	return
}

func (fs *FrontendSession) Set(k string, v interface{}) {
	fs.l.Lock()
	defer fs.l.Unlock()

	fs.settings[k] = v
}

func (fs *FrontendSession) Push(k string) error {
	return fs.service.Import(fs.sid, k, fs.Get(k))
}

func (fs *FrontendSession) PushAll() error {
	return fs.service.ImportAll(fs.sid, fs.settings)
}

func dclone(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
