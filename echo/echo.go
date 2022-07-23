package echo

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type Echo struct {
	Server       *http.Server
	Listener     net.Listener
	startupMutex sync.RWMutex
}

func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func New() (e *Echo) {
	e = &Echo{
		Server: new(http.Server),
	}
	e.Server.Handler = e
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

// ListenerAddr returns net.Addr for Listener
func (e *Echo) ListenerAddr() net.Addr {
	e.startupMutex.RLock()
	defer e.startupMutex.RUnlock()
	if e.Listener == nil {
		return nil
	}
	return e.Listener.Addr()
}
