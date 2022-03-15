package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	bindServer := NewServer(":1000", nil, nil)
	go bindServer.Run()
	defer bindServer.Shutdown()
	time.Sleep(1 * time.Second)

	t.Run("test success run server", func(t *testing.T) {
		s := NewServer(":2000", nil, nil)
		assert.NotNil(t, s)

		go func() {
			if err := s.Run(); err != nil {
				t.Error(err)
			}
		}()

		assert.NoError(t, s.Shutdown())
	})

	t.Run("test failed run with gracefully shutdown", func(t *testing.T) {
		s := NewServer(":1000", nil, nil)
		assert.NotNil(t, s)
		assert.NoError(t, s.Run())
		assert.NoError(t, s.Shutdown())
	})

	t.Run("test server ping handler", func(t *testing.T) {
		handler := NewServer(":2000", nil, nil).route()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Body.String(), `{"success": "ping"}`)
	})
}
