package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	s := &BotServer{}
	s.initDefaults()
	s.registerMetrics()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	s.engine.ServeHTTP(w, req)

	for _, metric := range []string{"go_gc", "go_goroutines", "go_memstats_alloc"} {
		assert.Contains(t, w.Body.String(), metric, "go metrics missing %s", metric)
	}
}

func TestDefaultGet(t *testing.T) {
	s := &BotServer{}
	s.initDefaults()
	s.registerDefaultHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	s.engine.ServeHTTP(w, req)

	assert.Equal(t, "Running Rivi", w.Body.String(), "default handler")
}

func TestDefaultGetInPath(t *testing.T) {
	s := &BotServer{Uri: "/example"}
	s.initDefaults()
	s.registerDefaultHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/example", nil)
	s.engine.ServeHTTP(w, req)

	assert.Equal(t, "Running Rivi", w.Body.String(), "default handler")
}
