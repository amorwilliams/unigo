package service

import (
	"fmt"
	"syscall"
)

type KillSignal struct {
	Signal syscall.Signal
}

func (ks KillSignal) String() string {
	return fmt.Sprintf("Got kill signal %q", ks.Signal)
}

type ServiceRegistered struct {
	ServiceInfo *ServiceInfo
}

func (sr ServiceRegistered) String() string {
	return fmt.Sprintf("Service %q registered", sr.ServiceInfo.Name)
}

type ServiceUnregistered struct {
	ServiceInfo *ServiceInfo
}

func (sr ServiceUnregistered) String() string {
	return fmt.Sprintf("Service %q unregistered", sr.ServiceInfo.Name)
}
