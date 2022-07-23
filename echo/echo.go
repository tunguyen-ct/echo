package echo

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"sync"
)

type Echo struct {
	Server          *http.Server
	Listener        net.Listener
	startupMutex    sync.RWMutex
	ListenerNetwork string
}

type Context interface {
	// Request returns `*http.Request`.
	Request() *http.Request
}

type Route interface{}

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c Context) error

func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello worlddddddd")
}

func New() (e *Echo) {
	e = &Echo{
		Server:          new(http.Server),
		ListenerNetwork: "tcp",
	}
	e.Server.Handler = e
	return
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func newListener(address, network string) (*tcpKeepAliveListener, error) {
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, ErrInvalidListenerNetwork
	}
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}

var (
	ErrInvalidListenerNetwork = errors.New("invalid listener network")
)

func (e *Echo) configureServer(s *http.Server) error {
	l, err := newListener(s.Addr, e.ListenerNetwork)
	if err != nil {
		return err
	}

	e.Listener = l
	s.Handler = e

	return nil
}

func (e *Echo) Start(address string) error {
	e.startupMutex.Lock()
	e.Server.Addr = address
	if err := e.configureServer(e.Server); err != nil {
		e.startupMutex.Unlock()
		return err
	}
	e.startupMutex.Unlock()
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

func handlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

func (e *Echo) add(host, method, path string, handler HandlerFunc) *Route {
	name := handlerName(handler)
	fmt.Println(name)
	return nil
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (e *Echo) Add(method, path string, handler HandlerFunc) *Route {
	return e.add("", method, path, handler)
}

// GET registers a new GET route for a path with matching handler in the router
// with optional route-level middleware.
func (e *Echo) GET(path string, h HandlerFunc) *Route {
	return e.Add(http.MethodGet, path, h)
}
