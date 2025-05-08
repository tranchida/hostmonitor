package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestHostEndpoint(t *testing.T) {
	e := newEngine(zap.NewNop())

	req := httptest.NewRequest(http.MethodGet, "/host", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, rec.Code)
	}

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatalf("failed to get hostname: %v", err)
	}

	body := rec.Body.String()
	if !strings.Contains(body, hostname) {
		t.Errorf("response body does not contain hostname, got: %s", body)
	}
}

func TestStaticFileNotFound(t *testing.T) {
	e := newEngine(zap.NewNop())

	// Try to access a non-existent static file
	req := httptest.NewRequest(http.MethodGet, "/static/nonexistentfile.txt", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, rec.Code)
	}
}
