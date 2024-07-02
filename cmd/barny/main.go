package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/blade2005/literate-barnacle/internal/site"
)

// Main should be so simple that it cannot have bugs, as it will generally be a
// function that is never run by tests.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run(ctx, site.NewSite())
}

type serverI interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Logger() *slog.Logger
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
}

func run(
	ctx context.Context,
	server serverI) {

	server.HandleFunc("/", hello)
	server.HandleFunc("/bye", goodbye)

	server.Logger().Info("server starting")
	serverErrors := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		serverErrors <- server.ListenAndServe()
	}()

	// This may seem a bit awkward, but imagine the case where `barny` runs
	// multiple ... what's smaller than a micro? But yeah, multiple servers
	// binding multiple ports.
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			server.Logger().Error(err.Error())
		}
	case <-ctx.Done():
		server.Logger().Info("server stopping")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Logger().Error(err.Error())
		}

		server.Logger().Info("Server shutdown complete. Exitting.")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("goodbye"))
}
