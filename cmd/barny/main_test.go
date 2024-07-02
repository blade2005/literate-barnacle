package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blade2005/literate-barnacle/internal/site"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	wg := sync.WaitGroup{}

	listener, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	server := site.NewSite(
		site.Listener(listener),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		run(ctx, server)
	}()

	t.Parallel() // all tests will run in parallel
	t.Run("hello", func(t *testing.T) {
		expected := []byte("hello")
		response, err := http.Get("http://127.0.0.1:" + fmt.Sprint(listener.Addr().(*net.TCPAddr).Port))
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		result, err := io.ReadAll(response.Body)
		response.Body.Close()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, expected, result)
	})

	t.Run("goodbye", func(t *testing.T) {
		expected := []byte("goodbye")

		response, err := http.Get("http://127.0.0.1:" + fmt.Sprint(listener.Addr().(*net.TCPAddr).Port) + "/bye")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		result, err := io.ReadAll(response.Body)
		response.Body.Close()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.Equal(t, expected, result)
	})

	cancel()
	wg.Wait()
}

func TestRunShutdown(t *testing.T) {
	// WARN: If this test if failing, it means we're exiting while threads
	// are still running / before shutdown is actually complete. That would
	// be a problem, and it's _really_ easy to get that wrong.
	wg := sync.WaitGroup{}
	listener, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var output bytes.Buffer

	server := site.NewSite(
		site.Stdout(&output),
		site.Listener(listener),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		run(ctx, server)
	}()
	cancel()
	wg.Wait()

	results := strings.Split(output.String(), "\n")
	assert.Contains(t, results[len(results)-2], "Server shutdown complete.")
}
