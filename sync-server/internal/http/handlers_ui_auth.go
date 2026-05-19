package serverhttp

import (
	"net/http"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		s.render(w, r, http.StatusServiceUnavailable, "login.html", pageData{Title: "Login", Error: "Auth service is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "login.html", pageData{Title: "Login", Error: "Invalid form submission."})
		return
	}

	next := safeNextPath(r.FormValue("next"))
	session, err := s.authService.Login(r.Context(), auth.LoginRequest{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		_ = s.audit(r.Context(), "", "", "user", "login_failed", map[string]string{"surface": "web"})
		s.render(w, r, http.StatusUnauthorized, "login.html", pageData{Title: "Login", Error: "Invalid credentials.", Next: next})
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	_ = s.audit(r.Context(), session.Principal.WorkspaceID, session.Principal.UserID, "user", "login_succeeded", map[string]string{"surface": "web"})
	if next == "" {
		next = "/dashboard"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (s *Server) signupForm(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		s.render(w, r, http.StatusServiceUnavailable, "signup.html", pageData{Title: "Create account", Error: "Auth service is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "signup.html", pageData{Title: "Create account", Error: "Invalid form submission."})
		return
	}

	session, err := s.authService.Signup(r.Context(), auth.SignupRequest{
		Email:         r.FormValue("email"),
		Username:      r.FormValue("username"),
		Password:      r.FormValue("password"),
		WorkspaceName: r.FormValue("workspace_name"),
	})
	if err != nil {
		s.render(w, r, http.StatusBadRequest, "signup.html", pageData{Title: "Create account", Error: err.Error()})
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	_ = s.audit(r.Context(), session.Principal.WorkspaceID, session.Principal.UserID, "user", "signup_succeeded", map[string]string{"surface": "web"})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) logoutForm(w http.ResponseWriter, r *http.Request) {
	if !s.requireCSRF(w, r) {
		return
	}
	if s.authService != nil {
		if cookie, err := r.Cookie(auth.SessionCookieName); err == nil {
			if principal, jti, err := s.authService.Authenticate(r.Context(), cookie.Value); err == nil {
				_ = s.authService.Logout(r.Context(), jti)
				_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "logout_succeeded", map[string]string{"surface": "web"})
			}
		}
	}
	auth.ClearSessionCookie(w, s.cookies)
	clearCSRFCookie(w, s.cookies)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
