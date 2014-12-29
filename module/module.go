package module

type Module interface {
	MonitorHandler()
	MasterHandler()
	ClientHandler()
	Start()
}
