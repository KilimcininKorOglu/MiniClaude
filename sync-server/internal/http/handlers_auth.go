package serverhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service is not available")
		return
	}

	var request auth.SignupRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := s.authService.Signup(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	writeJSON(w, http.StatusCreated, map[string]any{"user": session.Principal})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		writeError(w, http.StatusServiceUnavailable, "auth service is not available")
		return
	}

	var request auth.LoginRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := s.authService.Login(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	writeJSON(w, http.StatusOK, map[string]any{"user": session.Principal})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if s.authService != nil {
		if cookie, err := r.Cookie(auth.SessionCookieName); err == nil {
			if _, jti, err := s.authService.Authenticate(r.Context(), cookie.Value); err == nil {
				_ = s.authService.Logout(r.Context(), jti)
			}
		}
	}

	auth.ClearSessionCookie(w, s.cookies)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.authenticateRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": principal})
}

func (s *Server) authenticateRequest(w http.ResponseWriter, r *http.Request) (auth.Principal, bool) {
	if s.authService == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return auth.Principal{}, false
	}

	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return auth.Principal{}, false
	}

	principal, _, err := s.authService.Authenticate(r.Context(), cookie.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return auth.Principal{}, false
	}
	return principal, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body")
	}
	return nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
