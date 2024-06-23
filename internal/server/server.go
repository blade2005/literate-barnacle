package server

import (
	"log/slog"
	"net/http"
	"os"
)

type ServerConfig func(*Server)

type Server struct {
	logger  *slog.Logger
	handler http.Handler
	stdout  *os.File
	stderr  *os.File
}

func NewServer(cfgs ...ServerConfig) *Server {
	var s *Server = &Server{}
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

	s.addRoutes()

	return s
}

func (s *Server) addRoutes() {
	var mux *http.ServeMux = s.handler.(*http.ServeMux)
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/bye", goodbye)
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func (s *Server) Logger() *slog.Logger {
	return s.logger
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("goodbye"))
}
