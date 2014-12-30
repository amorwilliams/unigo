package service

import (
	"encoding/json"
	"sync"
)

const (
	MAX_CLIENTS = 1000
)

type ClientInfo map[string]string

type ClientMap map[string]ClientInfo

type ConnectionService struct {
	serverId     string
	connCount    int64
	loginedCount int64
	logined      ClientMap
	sync.Mutex
}

func NewConnectionService(serverId, config string) *ConnectionService {
	return &ConnectionService{
		serverId:     serverId,
		connCount:    0,
		loginedCount: 0,
		logined:      make(ClientMap, MAX_CLIENTS),
	}
}

//Add Client
func (cs *ConnectionService) Add(uid, info string) {
	cs.Lock()
	defer cs.Unlock()

	if _, ok := cs.logined[uid]; !ok {
		cs.loginedCount++
	}

	var cinfo ClientInfo
	json.Unmarshal([]byte(info), &cinfo)

	cs.logined[uid] = cinfo
}

func (cs *ConnectionService) Remove(uid string) {
	cs.Lock()
	defer cs.Unlock()

	if _, ok := cs.logined[uid]; !ok {
		cs.loginedCount--
	}

	delete(cs.logined, uid)
}

//Update client info
func (cs *ConnectionService) UpdateInfo(uid, info string) {
	cs.Lock()
	defer cs.Unlock()

	cinfo, ok := cs.logined[uid]
	if !ok {
		return
	}

	var newinfo ClientInfo
	json.Unmarshal([]byte(info), &newinfo)

	for k, v := range newinfo {
		cinfo[k] = v
	}
}

func (cs *ConnectionService) IncrCount() {
	cs.Lock()
	defer cs.Unlock()

	cs.connCount++
}

func (cs *ConnectionService) DecrCount(uid string) {
	cs.Lock()
	defer cs.Unlock()

	if cs.connCount > 0 {
		cs.connCount--
	}

	cs.Remove(uid)
}
