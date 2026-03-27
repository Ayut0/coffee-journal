package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ayut0/coffee-journal/api/builder"
	"github.com/Ayut0/coffee-journal/api/config"
	server "github.com/Ayut0/coffee-journal/api/internal/http"
	"github.com/stretchr/testify/assert"
)

func newTestDep() *builder.Dependency {
	return &builder.Dependency{
		Cfg: &config.Config{Port: "8080"},
	}
}

func TestHealthRoute(t *testing.T) {
	e := server.NewServer(newTestDep())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestStubRoutes(t *testing.T) {
	e := server.NewServer(newTestDep())
	routes := []string{"/api/beans", "/api/tastings", "/api/search"}
	for _, path := range routes {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotImplemented, rec.Code, "path: %s", path)
	}
}

func TestWrongMethod(t *testing.T) {
	e := server.NewServer(newTestDep())
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
