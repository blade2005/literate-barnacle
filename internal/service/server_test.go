package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func TestHello(t *testing.T) {
	srv := NewService()
	srv.AddRoute("/", hello)
	server := httptest.NewServer(srv.Handler())
	defer server.Close()
	var err error

	res, err := http.Get(server.URL)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "hello", string(greeting))
}
