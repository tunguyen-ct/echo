package echo

import (
	"fmt"
	"net"
	"net/http"
)

type Echo struct {
	Server   *http.Server
	Listener net.Listener
}

func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func New() (e *Echo) {
	e = &Echo{
		Server: new(http.Server),
	}
	return
}

func (e *Echo) configureServer(s *http.Server) error {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return err
	}

	e.Listener = listener
	s.Handler = e

	return nil
}

func (e *Echo) Start() error {
	if err := e.configureServer(e.Server); err != nil {
		return err
	}
	return e.Server.Serve(e.Listener)
}