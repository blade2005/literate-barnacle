package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	service "github.com/blade2005/literate-barnacle/internal/service"
)

// Main should be so simple that it cannot have bugs, as it will generally be a
// function that is never run by tests.
func main() {

	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service := service.NewService(service.Addr(":8080"))

	server := &http.Server{
		Handler: service.Handler(),
		Addr:    *service.Addr(),
	}

	run(ctx, server, service)
}

type serverI interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func run(
	ctx context.Context,
	server serverI,
	service *service.Service) {

	service.AddRoute("/", hello)
	service.AddRoute("/bye", goodbye)

	service.Logger().Info("server starting")
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
			service.Logger().Error(err.Error())
		}
	case <-ctx.Done():
		service.Logger().Info("server stopping")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			service.Logger().Error(err.Error())
		}

		wg.Done()
		service.Logger().Info("Server shutdown complete. Exitting.")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("goodbye"))
}
