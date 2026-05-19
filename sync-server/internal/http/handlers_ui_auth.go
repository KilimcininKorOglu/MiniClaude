package serverhttp

import (
	"net/http"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

func (s *Server) loginForm(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		s.render(w, http.StatusServiceUnavailable, "login.html", pageData{Title: "Login", Error: "Auth service is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, "login.html", pageData{Title: "Login", Error: "Invalid form submission."})
		return
	}

	next := safeNextPath(r.FormValue("next"))
	session, err := s.authService.Login(r.Context(), auth.LoginRequest{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		s.render(w, http.StatusUnauthorized, "login.html", pageData{Title: "Login", Error: "Invalid credentials.", Next: next})
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	if next == "" {
		next = "/dashboard"
	}
	http.Redirect(w, r, next, http.StatusSeeOther)
}

func (s *Server) signupForm(w http.ResponseWriter, r *http.Request) {
	if s.authService == nil {
		s.render(w, http.StatusServiceUnavailable, "signup.html", pageData{Title: "Create account", Error: "Auth service is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, "signup.html", pageData{Title: "Create account", Error: "Invalid form submission."})
		return
	}

	session, err := s.authService.Signup(r.Context(), auth.SignupRequest{
		Email:         r.FormValue("email"),
		Username:      r.FormValue("username"),
		Password:      r.FormValue("password"),
		WorkspaceName: r.FormValue("workspace_name"),
	})
	if err != nil {
		s.render(w, http.StatusBadRequest, "signup.html", pageData{Title: "Create account", Error: err.Error()})
		return
	}

	auth.SetSessionCookie(w, session.Token, session.ExpiresAt, s.cookies)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s *Server) logoutForm(w http.ResponseWriter, r *http.Request) {
	if s.authService != nil {
		if cookie, err := r.Cookie(auth.SessionCookieName); err == nil {
			if _, jti, err := s.authService.Authenticate(r.Context(), cookie.Value); err == nil {
				_ = s.authService.Logout(r.Context(), jti)
			}
		}
	}
	auth.ClearSessionCookie(w, s.cookies)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
