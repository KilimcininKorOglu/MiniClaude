package serverhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	handler := New(nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}

func TestMeRequiresSessionCookie(t *testing.T) {
	handler := New(nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}
