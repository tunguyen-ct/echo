package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

/*
Package echo implements high performance, minimalist Go web framework.

Example:

  package main

  import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
  )

  // Handler
  func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
  }

  func main() {
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", hello)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
  }

Learn more at https://echo.labstack.com
*/

type hotdog int

func (h hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func handlerName(h interface{}) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		fmt.Println("huhu")
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	fmt.Println("hihi")
	return t.String()
}

func sayHi() {
	fmt.Println("Hi")
}

type Echo struct {
	Server   *http.Server
	Listener net.Listener
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
	var h hotdog
	s.Handler = h

	return nil
}

func (e *Echo) Start() error {
	if err := e.configureServer(e.Server); err != nil {
		return err
	}
	return e.Server.Serve(e.Listener)
}

func NewServer() {
	name := handlerName(sayHi)
	fmt.Println(name)
	var h hotdog

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return
	}

	s := &http.Server{
		Addr:           ":8080",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.Serve(listener)
}

func main() {
	e := New()
	log.Fatal(e.Start())
}
