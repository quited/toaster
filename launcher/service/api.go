package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

type ApiEndpoint struct {
	Protocol string `json:"protocol"`
	Path     string `json:"address"`

	client *http.Client
}

func NewApiEndpoint(endpoint string) (obj *ApiEndpoint, err error) {
	obj = new(ApiEndpoint)
	splitRes := strings.SplitN(endpoint, "://", 2)

	obj.Protocol = splitRes[0]

	switch obj.Protocol {
	case "http", "unix", "tcp":
		obj.Path = splitRes[1]
		obj.MakeClient()
		return
	}

	return nil, errors.New("invalid endpoint")
}

func NewApiEndpointFromJsonStream(r io.Reader) (*ApiEndpoint, error) {
	obj := &ApiEndpoint{}
	err := json.NewDecoder(r).Decode(obj)
	if err != nil {
		return nil, err
	}
	switch obj.Protocol {
	case "http", "unix", "tcp":
		obj.MakeClient()
		return obj, nil
	}

	return nil, errors.New("invalid endpoint")
}

func (a *ApiEndpoint) MakeClient() {
	fakeClient := func(network string) *http.Client {
		fakeDial := func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			type resType struct {
				conn net.Conn
				err  error
			}
			done, resChan := func() (done chan struct{}, res chan resType) {
				done, res = make(chan struct{}), make(chan resType)
				go func() {
					conn, err := net.Dial(network, a.Path)
					res <- resType{conn, err}
					done <- struct{}{}
				}()
				return done, res
			}()

			select {
			case <-ctx.Done():
				return nil, errors.New("dial timeout")
			case <-done:
				res := <-resChan
				return res.conn, res.err
			}
		}
		tr := &http.Transport{
			DialContext: fakeDial,
		}
		client := &http.Client{Transport: tr}
		return client
	}
	switch a.Protocol {
	case "unix", "tcp":
		a.client = fakeClient(a.Protocol)
	case "http":
		a.client = http.DefaultClient
	}
}

func (a *ApiEndpoint) SendNoResponse(data interface{}) error {
	_, err := a.Send(data)
	return err
}

func (a *ApiEndpoint) Send(data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(jsonData)
	switch a.Protocol {
	case "unix", "tcp":
		return a.client.Post("http://dummy/manage", "application/json", buf)
	case "http":
		return a.client.Post("http://"+a.Path+"/manage", "application/json", buf)
	}
	panic("unreachable code")
}
