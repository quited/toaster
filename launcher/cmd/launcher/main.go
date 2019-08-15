package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"launcher/service"
	"log"
	"os"
)

func FatalWhileErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	config, err := loadConfig()
	FatalWhileErr(err)

	if *installService {
		FatalWhileErr(service.InstallService(config.Service.Name, config.Service.Description, config.Service.ProgramFile))
		return
	}

	if *removeService {
		FatalWhileErr(service.RemoveService(config.Service.Name))
		return
	}

	if *outputServiceName {
		fmt.Println("toast_workload:" + config.Service.Name)
		return
	}

	if *startService != "" {
		config.Manager.ApiEndpoint = *startService
		srv, err := config.LoadService()
		FatalWhileErr(err)
		stat, err := srv.Status()
		FatalWhileErr(err)
		if stat.State != svc.Stopped {
			FatalWhileErr(srv.ControlService(svc.Stop, svc.Stopped))
		}
		FatalWhileErr(srv.StartService())
		fmt.Println(config.Service.ApiEndpoint)
		return
	}

	if *stopService {
		srv, err := config.LoadService()
		FatalWhileErr(err)
		FatalWhileErr(srv.ControlService(svc.Stop, svc.Stopped))
		return
	}

	if *showStatus {
		srv, err := config.LoadService()
		FatalWhileErr(err)
		stat, err := srv.Status()
		FatalWhileErr(err)
		switch stat.State {
		case svc.Running:
			fmt.Println("running")
		default:
			fmt.Println("stopped")
		}
		return
	}

	fmt.Println("Please provide one command.")
	flag.Usage()
	os.Exit(-1)
}