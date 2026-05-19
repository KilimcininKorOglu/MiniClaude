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

func TestLoginPageRenders(t *testing.T) {
	handler := New(nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Body.String() == "" {
		t.Fatal("expected login page body")
	}
}

func TestDashboardRedirectsWithoutSession(t *testing.T) {
	handler := New(nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", response.Code)
	}
	if response.Header().Get("Location") != "/login" {
		t.Fatalf("expected redirect to /login, got %q", response.Header().Get("Location"))
	}
}
