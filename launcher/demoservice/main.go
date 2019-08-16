package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"os"
)

func main() {
	interactive, err := svc.IsAnInteractiveSession()
	if err != nil {
		elog.Error(1, fmt.Sprintf("failed to detect session type"))
		os.Exit(-1)
	}
	ws, err := NewWindowsService("config.json")
	if err != nil {
		elog.Error(1, fmt.Sprintf("failed to create windows service"))
		os.Exit(-1)
	}
	ws.Run(interactive)
}
