package serverhttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}

func TestMeRequiresSessionCookie(t *testing.T) {
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestLoginPageRenders(t *testing.T) {
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
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
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
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

func TestClientRenameRedirectsWithoutSession(t *testing.T) {
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodPost, "/ui/clients/rename", strings.NewReader("client_id=client-1&name=Desk"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", response.Code)
	}
	if response.Header().Get("Location") != "/login" {
		t.Fatalf("expected redirect to /login, got %q", response.Header().Get("Location"))
	}
}

func TestSecureHeadersAreApplied(t *testing.T) {
	handler := New(nil, nil, nil, nil, nil, nil, auth.CookieConfig{})
	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("expected nosniff header")
	}
	if response.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatal("expected frame denial header")
	}
	if response.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("expected content security policy")
	}
}

func TestRequireCSRFRejectsMissingToken(t *testing.T) {
	server := &Server{}
	request := httptest.NewRequest(http.MethodPost, "/ui/providers", nil)
	response := httptest.NewRecorder()

	if server.requireCSRF(response, request) {
		t.Fatal("expected missing CSRF token to be rejected")
	}
	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", response.Code)
	}
}

func TestRequireCSRFAcceptsMatchingToken(t *testing.T) {
	server := &Server{}
	request := httptest.NewRequest(http.MethodPost, "/ui/providers", strings.NewReader("csrf_token=token-value"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(&http.Cookie{Name: csrfCookieName, Value: "token-value"})
	response := httptest.NewRecorder()

	if !server.requireCSRF(response, request) {
		t.Fatal("expected matching CSRF token to be accepted")
	}
}
