package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandlerGET(t *testing.T) {
	// Create a request to pass to the handler
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// Record the response
	rr := httptest.NewRecorder()

	HealthHandlerGET(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := "Running"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}
