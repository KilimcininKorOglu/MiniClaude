package serverhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	handler := New(nil)
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}
