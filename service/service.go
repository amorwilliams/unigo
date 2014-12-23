package service

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

import (
	l4g "code.google.com/p/log4go"
	"github.com/amorwilliams/unigo/utils"
)

type ServiceDelegate interface {
	Setup()
	Started()
	Stopped()
	Registered()
	Unregistered()
	CreatePeer(conn Connector) (Connector, error)
	NewMessageReceived(conn Connector, msg string)
	NewDataReceived(conn Connector, data []byte)
}

type ClientInfo struct {
	Address net.Addr
}

type Service struct {
	*ServiceInfo
	Delegate       ServiceDelegate
	Status         byte
	registeredChan chan bool
	shutdownChan   chan bool

	clientMutex sync.Mutex
	ClientInfo  map[string]ClientInfo

	// for sending the signal into mux()
	doneChan chan bool

	// for waiting for all shutdown operations
	doneGroup *sync.WaitGroup

	shuttingDown bool
}

func CreateService(sd ServiceDelegate, si *ServiceInfo) (s *Service) {
	s = &Service{
		Delegate:       sd,
		ServiceInfo:    si,
		registeredChan: make(chan bool),
		shutdownChan:   make(chan bool),
		ClientInfo:     make(map[string]ClientInfo),
		shuttingDown:   false,
	}

	if si.Name != "UnigoDaemon" {
		// Listen for admin requests
		// go s.serveAdminRequests()
	}

	return
}

// Notifies the cluster your service is ready to handle requests
func (s *Service) Register() {
	s.registeredChan <- true
}

func (s *Service) register() {
	// this version must be run from the mux() goroutine
	if s.Registered {
		return
	}

	// TODO: Register
	// err := skynet.GetServiceManager().Register(s.ServiceInfo.UUID)
	// if err != nil {
	// 	log.Println(log.ERROR, "Failed to register service: "+err.Error())
	// }

	s.Registered = true
	l4g.Info("%+v\n", ServiceRegistered{s.ServiceInfo})
	s.Delegate.Registered()
}

// Leave your service online, but notify the cluster it's not currently accepting new requests
func (s *Service) Unregister() {
	s.registeredChan <- false
}

func (s *Service) unregister() {
	// this version must be run from the mux() goroutine
	if !s.Registered {
		return
	}

	// err := skynet.GetServiceManager().Unregister(s.UUID)
	// if err != nil {
	// 	log.Println(log.ERROR, "Failed to unregister service: "+err.Error())
	// }

	s.Registered = false
	l4g.Info("%+v\n", ServiceUnregistered{s.ServiceInfo})
	s.Delegate.Unregistered() // Call user defined callback
}

func (s *Service) Start() (done *sync.WaitGroup) {

	go s.listen(":9000")

	// Watch signals for shutdown
	c := make(chan os.Signal)
	go watchSignals(c, s)

	s.doneChan = make(chan bool)

	s.doneGroup = &sync.WaitGroup{}
	s.doneGroup.Add(1)

	go func() {
		s.mux()
		s.doneGroup.Done()
	}()
	done = s.doneGroup

	// if r, err := config.Bool(s.Name, s.Version, "service.register"); err == nil {
	// 	s.Registered = r
	// }
	//s.Registered = true

	// err := skynet.GetServiceManager().Add(*s.ServiceInfo)
	// if err != nil {
	// 	log.Println(log.ERROR, "Failed to add service: "+err.Error())
	// }

	s.Register()

	s.Delegate.Started()

	return
}

func (s *Service) Shutdown() {
	if s.shuttingDown {
		return
	}

	s.registeredChan <- false
	s.shutdownChan <- true
}

func (s *Service) shutdown() {
	if s.shuttingDown {
		return
	}
	s.shuttingDown = true

	s.doneGroup.Add(1)

	close(s.doneChan)

	s.Delegate.Stopped()
	s.doneGroup.Done()
}

// TODO: Currently unimplemented
func (s *Service) IsTrusted(addr net.Addr) bool {
	return false
}

func (s *Service) getClientInfo(clientID string) (ci ClientInfo, ok bool) {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	ci, ok = s.ClientInfo[clientID]
	return
}

func (s *Service) listen(addr string) {
	// var laddr *net.TCPAddr
	// laddr, err := net.ResolveTCPAddr("tcp", addr)
	// if err != nil {
	// 	l4g.Error(err)
	// }

	http.HandleFunc("/", serveWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		l4g.Exit("ListenAndServe", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l4g.Error("Upgrade:", err)
		return
	}
	defer func() {
		ws.Close()
		l4g.Info("Client disconnected: ", ws.RemoteAddr())
	}()

	l4g.Info("Client connected: ", ws.RemoteAddr())
	conn := NewConn(ws)
	connectionChan <- conn

	go conn.writePump()
	conn.readPump()
}

func (s *Service) mux() {
	for {
		select {
		case conn := <-connectionChan:
			go func() {
				clientID := utils.NewUUID()

				s.clientMutex.Lock()
				s.ClientInfo[clientID] = ClientInfo{
					Address: conn.ws.RemoteAddr(),
				}
				s.clientMutex.Unlock()

				connector, err := s.Delegate.CreatePeer(conn)
				if err != nil {
					l4g.Error("CreatePeer error: ", err)
					return
				}

				for {
					select {
					case msg := <-conn.recMsg:
						s.Delegate.NewMessageReceived(connector, string(msg))
					case data := <-conn.recData:
						s.Delegate.NewDataReceived(connector, data)
					case <-conn.doneChan:
						break
					}
				}

			}()
		case register := <-s.registeredChan:
			if register {
				s.register()
			} else {
				s.unregister()
			}
		case <-s.shutdownChan:
			s.shutdown()
		case <-s.doneChan:
			break
		}
	}
}

func watchSignals(c chan os.Signal, s *Service) {
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSEGV /*syscall.SIGSTOP,*/, syscall.SIGTERM)

	for {
		select {
		case sig := <-c:
			switch sig.(syscall.Signal) {
			// Trap signals for clean shutdown
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT,
				syscall.SIGSEGV /*syscall.SIGSTOP,*/, syscall.SIGTERM:
				l4g.Info("%+v", KillSignal{sig.(syscall.Signal)})
				s.Shutdown()
				return
			}
		}
	}
}
