package main

import (
	"fmt"
	"github.com/quited/toaster/launcher/service"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"net/http"
	"os"
)

var elog debug.Log

type ServiceJob func(config *service.Config)

type WindowsService struct {
	conf *service.Config
	job  ServiceJob
}

func (m *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go m.job(m.conf)
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func (w *WindowsService) SetJob(job ServiceJob) {
	w.job = job
}

func NewWindowsService(configFile string) (obj *WindowsService, err error) {
	obj = new(WindowsService)
	obj.conf, err = service.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}
	obj.job = func(config *service.Config) {
		if err := http.ListenAndServe(config.Service.ApiEndpoint, http.DefaultServeMux); err != nil {
			elog.Error(1, fmt.Sprintf("%s service failed: %v", config.Service.Name, err))
			os.Exit(-1)
		}
	}
	return
}

func (w *WindowsService) Run(isDebug bool) {
	name := w.conf.Service.Name
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &WindowsService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
