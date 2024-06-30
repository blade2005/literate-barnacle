package site

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func deadpool(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I love deadlocking servers!"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there deadpool"))
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("goodbye deadpond"))
}

func FuzzSite(f *testing.F) {
	wg := sync.WaitGroup{}

	// TODO: Probably we should not be using networking here...
	// https://github.com/golang/go/issues/14200
	listener, err := net.Listen("tcp", "localhost:0")
	if !assert.NoError(f, err) {
		f.FailNow()
	}

	server := NewSite(Listener(listener))
	server.HandleFunc("/", deadpool)
	server.HandleFunc("/bye", goodbye)
	server.HandleFunc("/hello", hello)

	testCases := []string{
		"",
		"bye",
		"hello",
	}

	responses := []string{
		"I love deadlocking servers!",
		"hello there deadpool",
		"goodbye deadpond",
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(listener)
	}()

	addr := "http://" + listener.Addr().String() + "/"

	for _, uri := range testCases {
		f.Add(uri, 5)
	}

	f.Fuzz(func(t *testing.T, uri string, multiplexity int) {
		if multiplexity == 0 {
			t.Skip()
		}

		if multiplexity < 0 {
			multiplexity = -multiplexity
		}

		multiplexity = multiplexity % 1000

		url := fmt.Sprintf("%s%s", addr, url.QueryEscape(uri))

		// Multithreaded test to ensure we don't have unexpected data race conditions.
		// rwlock ensures we can block until all goroutines have been created.
		wg := sync.WaitGroup{}

		for range multiplexity {
			wg.Add(1)
			go func() {
				defer wg.Done()
				testGet(t, &url, &responses)
			}()
		}

		wg.Wait()
		// If there is a connection tracking table on lo, we will fill
		// it. Slow down. If you comment this out, then run conntrack
		// -F in a loop in another terminal or something.
		//
		<-time.After(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if !assert.NoError(f, err) {
		f.FailNow()
	}
}

func testGet(t *testing.T, req *string, responses *[]string) {
	res, err := http.Get(*req)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Contains(t, *responses, string(greeting))
}
