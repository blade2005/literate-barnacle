package service

import (
	"io"
	"log/slog"
	"net/http"
	"os"
)

type ServiceConfig func(*Service)

type httpHandler interface {
	http.Handler
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Service struct {
	logger  *slog.Logger
	handler httpHandler
	stdout  io.Writer
	stderr  io.Writer
	addr    *string
}

func NewService(cfgs ...ServiceConfig) *Service {
	var s *Service = &Service{}
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
	if s.handler == nil {
		s.handler = http.NewServeMux()
	}
	if s.addr == nil {
		// We cannot take the address of a constant; but we can take
		// the address of a variable which contains a constant.
		a := ":8080"
		s.addr = &a
	}

	return s
}

func Logger(l *slog.Logger) func(s *Service) {
	return func(s *Service) {
		s.logger = l
	}
}

func Handler(h httpHandler) func(s *Service) {
	return func(s *Service) {
		s.handler = h
	}
}

func Stdout(f io.Writer) func(s *Service) {
	return func(s *Service) {
		s.stdout = f
	}
}

func Stderr(f io.Writer) func(s *Service) {
	return func(s *Service) {
		s.stderr = f
	}
}

func Addr(a string) func(s *Service) {
	return func(s *Service) {
		s.addr = &a
	}
}

func (s *Service) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// By adding HandleFunc to an interface, and defining that our handler has this function, we don't have to typecast here anymore.

	// var mux *http.ServeMux = s.handler.(*http.ServeMux)
	s.handler.HandleFunc(pattern, handler)
}

func (s *Service) Handler() http.Handler {
	return s.handler
}

func (s *Service) Addr() *string {
	return s.addr
}

func (s *Service) Logger() *slog.Logger {
	return s.logger
}
