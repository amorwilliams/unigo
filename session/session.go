package session

import (
	"encoding/json"
	"fmt"
	"github.com/amorwilliams/unigo/service"
	"github.com/amorwilliams/unigo/utils"
	"time"
)

// type Session struct {
// 	Conn *service.Conn
// 	UUID string
// }

// SessionStore contains all data for one session process with specific id.
type SessionStore interface {
	Set(key, value interface{}) error      //set session value
	Get(key interface{}) interface{}       //get session value
	Delete(key interface{}) error          //delete session value
	SessionID() string                     //back current sessionID
	SessionRelease(conn service.Connector) // release the resource & save data to provider & return the data
	Flush() error                          //delete all data
}

type Provider interface {
	SessionInit(config string) error
	SessionRead(sid string) (SessionStore, error)
	SessionExist(sid string) bool
	// SessionRegenerate(oldsid, sid string) (SessionStore, error)
	SessionAdd(conn service.Connector) error
	SessionDestroy(sid string) error
	SessionAll() int
	SessionGC()
}

var providers = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	providers[name] = provider
}

type managerConfig struct {
	Gclifetime     int64  `json:"gclifetime"`
	Maxlifetime    int64  `json:"maxLifetime"`
	ProviderConfig string `json:"providerConfig"`
}

type Manager struct {
	provider Provider
	config   *managerConfig
}

func NewManager(providerName, config string) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}
	cf := new(managerConfig)
	err := json.Unmarshal([]byte(config), cf)
	if err != nil {
		return nil, err
	}
	//TODO:Add config

	err = provider.SessionInit(cf.ProviderConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{
		provider: provider,
		config:   cf,
	}, nil
}

func (m *Manager) RegistSession(conn service.Connector) (seesion SessionStore, err error) {
	if m.provider.SessionExist
}

func (m *Manager) UnregistSession(sid string) error {
	return m.provider.SessionDestroy(sid)
}

// Get SessionStore by its id.
func (m *Manager) GetSession(sid string) (sessions SessionStore, err error) {
	sessions, err = m.provider.SessionRead(sid)
	return
}

// Start session gc process.
// it can do gc in times after gc lifetime.
func (m *Manager) GC() {
	m.provider.SessionGC()
	time.AfterFunc(time.Duration(m.config.Gclifetime)*time.Second, func() { m.GC() })
}

// Get all active sessions count number.
func (m *Manager) GetActiveSession() int {
	return m.provider.SessionAll()
}

func (m *Manager) sessionId(conn service.Connector) (string, error) {
	return utils.NewUUID(), nil
}
