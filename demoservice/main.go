package main

import (
	"fmt"
	"github.com/quited/toaster/launcher/service"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"net/http"
	"os"
)

var elog debug.Log

type WindowsService struct{}

func (m *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	job := func() {
		http.Handle("/", &Service{})
		if err := http.ListenAndServe("localhost:1000", http.DefaultServeMux); err != nil {
			log.Fatalln("[ERROR] [Serve] ", err)
		}
	}

	go job()
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

type Service struct {
	manager *service.ApiEndpoint
}

func (s *Service) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/manage":
		var err error

		if req.Method != http.MethodPost {
			log.Println("[ERROR] [/manage] invalid request method ", req.Method)
			res.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if err := req.ParseForm(); err != nil {
			log.Println("[ERROR] [/manage] parsing form: ", err)
			res.WriteHeader(http.StatusNotAcceptable)
			return
		}

		s.manager, err = service.NewApiEndpointFromJsonStream(req.Body)
		defer req.Body.Close()
		if err != nil {
			log.Println("[ERROR] [/manage]", err)
			res.WriteHeader(http.StatusNotAcceptable)
			return
		}
	default:
		log.Println("[ERROR] [" + req.URL.Path + "] invalid path")
		res.WriteHeader(http.StatusNotFound)
		return
	}
}

func runService(name string, isDebug bool) {
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

func main() {
	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed %v", "demo_service", err))
		os.Exit(-1)
	}
	if interactive {
		runService("demo_service", true)
	} else {
		runService("demo_service", false)
	}
}
