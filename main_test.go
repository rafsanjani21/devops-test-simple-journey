package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"fmt"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, DevOps! version = %s\n", version)
	})
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Hello, DevOps!") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}