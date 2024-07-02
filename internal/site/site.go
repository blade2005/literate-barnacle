package site

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
)

type SiteConfig func(*Site)

type httpServer interface {
	Serve(net.Listener) error
	ServeTLS(net.Listener, string, string) error
	Shutdown(ctx context.Context) error
}

type httpServeMux interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
	Handle(string, http.Handler)
	http.Handler
}

type Site struct {
	server   httpServer
	logger   *slog.Logger
	listener net.Listener
	smux     httpServeMux
	stdout   io.Writer
	stderr   io.Writer
	addr     *string
}

func NewSite(cfgs ...SiteConfig) *Site {
	var s *Site = &Site{}
	for _, f := range cfgs {
		f(s)
	}
	if s.stdout == nil {
		s.stdout = os.Stdout
	}
	if s.stderr == nil {
		s.stderr = os.Stderr
	}
	if s.logger == nil {
		s.logger = slog.New(slog.NewJSONHandler(s.stdout, nil))
	}
	if s.addr == nil {
		// We cannot take the address of a constant; but we can take
		// the address of a variable which contains a constant.
		a := ":8080"
		s.addr = &a
	}
	if s.smux == nil {
		s.smux = http.NewServeMux()
	}
	if s.listener == nil {
		ln, err := net.Listen("tcp", *s.addr)
		if err != nil {
			panic(err)
		}
		s.listener = ln
	}
	if s.server == nil {
		s.server = &http.Server{
			Handler: s.smux,
			Addr:    *s.addr,
		}
	}

	return s
}

func Server(h httpServer) func(s *Site) {
	return func(s *Site) {
		s.server = h
	}
}

func Logger(l *slog.Logger) func(s *Site) {
	return func(s *Site) {
		s.logger = l
	}
}

func Listener(l net.Listener) func(s *Site) {
	return func(s *Site) {
		s.listener = l
	}
}

func ServeMux(h httpServeMux) func(s *Site) {
	return func(s *Site) {
		s.smux = h
	}
}

func Stdout(f io.Writer) func(s *Site) {
	return func(s *Site) {
		s.stdout = f
	}
}

func Stderr(f io.Writer) func(s *Site) {
	return func(s *Site) {
		s.stderr = f
	}
}

func Addr(a string) func(s *Site) {
	return func(s *Site) {
		s.addr = &a
	}
}

func (s *Site) Handle(pattern string, handler http.Handler) {
	s.smux.Handle(pattern, handler)
}

func (s *Site) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.smux.HandleFunc(pattern, handler)
}

func (s *Site) Serve(l net.Listener) error {
	return s.server.Serve(l)
}

func (s *Site) ServeTLS(l net.Listener, certFile, keyFile string) error {
	return s.server.ServeTLS(l, certFile, keyFile)
}

func (s *Site) ListenAndServe() error {
	return s.Serve(s.listener)
}

func (s *Site) ListenAndServeTLS(certFile, keyFile string) error {
	return s.ServeTLS(s.listener, certFile, keyFile)
}

func (s *Site) Logger() *slog.Logger {
	return s.logger
}

func (s *Site) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
