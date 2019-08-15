// +build windows

package service

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"time"
)

type CommandCallback func(cmd svc.Cmd) interface{}

type Service struct {
	ServiceName string

	ManagerApi *ApiEndpoint
	ServiceApi *ApiEndpoint

	callbacks map[svc.Cmd]CommandCallback
}

func NewService(serviceName string, managerApi string, serviceApi string) (l *Service, err error) {
	l = new(Service)
	l.ServiceName = serviceName

	managerApiEndpoint, err1 := NewApiEndpoint(managerApi)
	serviceApiEndpoint, err2 := NewApiEndpoint(serviceApi)

	if err1 != nil || err2 != nil {
		return nil, err1
	}

	l.ManagerApi, l.ServiceApi = managerApiEndpoint, serviceApiEndpoint
	return
}

func (l *Service) On(cmd svc.Cmd, cb CommandCallback) {
	if cb == nil {
		delete(l.callbacks, cmd)
		return
	}
	l.callbacks[cmd] = cb
}

const SERVICE_CONTROL_START = svc.Cmd(0)

func (l *Service) StartService() (err error) {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(l.ServiceName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	if err = l.ServiceApi.SendNoResponse(l.ManagerApi); err != nil {
		return err
	}
	if cb, ok := l.callbacks[SERVICE_CONTROL_START]; ok {
		if err = l.ServiceApi.SendNoResponse(cb(SERVICE_CONTROL_START)); err != nil {
			return err
		}
	}
	return
}

func (l *Service) ControlService(c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(l.ServiceName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()

	if cb, ok := l.callbacks[c]; ok {
		if err := l.ServiceApi.SendNoResponse(cb(c)); err != nil {
			return err
		}
	}
	status, err := s.Control(c)

	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func (l *Service) Status() (svc.Status, error) {
	m, err := mgr.Connect()
	if err != nil {
		return svc.Status{}, err
	}
	defer m.Disconnect()
	s, err := m.OpenService(l.ServiceName)
	if err != nil {
		return svc.Status{}, fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	stat, err := s.Query()
	if err != nil {
		return svc.Status{}, err
	}
	return stat, err
}
