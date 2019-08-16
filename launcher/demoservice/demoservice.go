package main

import (
	"fmt"
	"github.com/quited/toaster/launcher/service"
	"log"
	"net/http"
	"os"
)

type DemoService struct {
	manager *service.ApiEndpoint
}

func (s *DemoService) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

func (s *DemoService) GetJob() ServiceJob {
	return func(config *service.Config) {
		http.Handle("/", s)
		conf, err := config.LoadService()
		if err != nil {
			elog.Error(1, fmt.Sprintf("%s service failed: %v", config.Service.Name, err))
			os.Exit(-1)
		}
		if err := http.ListenAndServe(conf.ServiceApi.Path, http.DefaultServeMux); err != nil {
			elog.Error(1, fmt.Sprintf("%s service failed: %v", config.Service.Name, err))
			os.Exit(-1)
		}
	}
}
