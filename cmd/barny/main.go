package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/blade2005/literate-barnacle/internal/server"
)

func main() {
	s := server.NewServer()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Handler: s.Handler(),
		Addr:    ":8080",
	}

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		s.Logger().Info("server started")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger().Error(err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		<-done // Block on os signals
		s.Logger().Info("server stopping")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			s.Logger().Error(err.Error())
		}
	}()

	wg.Wait()
	s.Logger().Info("Server shutdown complete. Exitting.")

}
