package web

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nichojovi/tauria-test/cmd/internal"
	"github.com/nichojovi/tauria-test/cmd/web/api"
	"github.com/nichojovi/tauria-test/internal/utils/auth"
	myrouter "github.com/nichojovi/tauria-test/internal/utils/router"
	"gopkg.in/tylerb/graceful.v1"
)

type Opts struct {
	ListenAddress string
	AuthService   auth.AuthService
	Service       *internal.Service
}

type Handler struct {
	options     *Opts
	listenErrCh chan error
}

func New(o *Opts) *Handler {
	handler := &Handler{options: o}

	api.New(&api.Options{
		Prefix:         "/api",
		DefaultTimeout: 15,
		AuthService:    o.AuthService,
		Service:        o.Service,
	}).Register()

	return handler
}

func (h *Handler) Run() {
	log.Printf("Listening on %s", h.options.ListenAddress)
	h.listenErrCh <- Serve(h.options.ListenAddress, myrouter.WrapperHandler())
}

func (h *Handler) ListenError() <-chan error {
	return h.listenErrCh
}

var listenPort string
var cfgtestFlag bool

func Serve(hport string, handler http.Handler) error {

	checkConfigTest()

	l, err := Listen(hport)
	if err != nil {
		log.Fatalln(err)
	}

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	log.Println("starting serve on ", hport)
	return srv.Serve(l)
}

func Listen(hport string) (net.Listener, error) {
	var l net.Listener

	fd := os.Getenv("EINHORN_FDS")
	if fd != "" {
		sock, err := strconv.Atoi(fd)
		if err == nil {
			hport = "socketmaster:" + fd
			log.Println("detected socketmaster, listening on", fd)
			file := os.NewFile(uintptr(sock), "listener")
			fl, err := net.FileListener(file)
			if err == nil {
				l = fl
			}
		}
	}

	if listenPort != "" {
		hport = ":" + listenPort
	}

	checkConfigTest()

	if l == nil {
		var err error
		l, err = net.Listen("tcp4", hport)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

func checkConfigTest() {
	if cfgtestFlag == true {
		log.Println("config test mode, exiting")
		os.Exit(0)
	}
}
