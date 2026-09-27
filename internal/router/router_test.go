package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-practice/internal/router"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.New().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	expectedBody := "{\"status\":\"ok\"}\n"
	if recorder.Body.String() != expectedBody {
		t.Fatalf("expected body %q, got %q", expectedBody, recorder.Body.String())
	}
}
