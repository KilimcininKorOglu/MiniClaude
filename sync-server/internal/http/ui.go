package serverhttp

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

//go:embed web/templates/*.html web/static/*
var uiFiles embed.FS

type pageData struct {
	Title     string
	Principal auth.Principal
	Error     string
}

func (s *Server) templates() (*template.Template, error) {
	return template.ParseFS(uiFiles, "web/templates/*.html")
}

func (s *Server) render(w http.ResponseWriter, status int, name string, data pageData) {
	templates, err := s.templates()
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if _, ok := s.currentPrincipal(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.currentPrincipal(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.render(w, http.StatusOK, "login.html", pageData{Title: "Login"})
}

func (s *Server) signupPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.currentPrincipal(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.render(w, http.StatusOK, "signup.html", pageData{Title: "Create account"})
}

func (s *Server) dashboardPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	s.render(w, http.StatusOK, "dashboard.html", pageData{Title: "Dashboard", Principal: principal})
}

func (s *Server) providersPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	s.render(w, http.StatusOK, "providers.html", pageData{Title: "Providers", Principal: principal})
}

func (s *Server) clientsPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	s.render(w, http.StatusOK, "clients.html", pageData{Title: "Clients", Principal: principal})
}

func (s *Server) favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) staticFile(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(uiFiles)).ServeHTTP(w, r)
}

func (s *Server) requireHTMLAuth(w http.ResponseWriter, r *http.Request) (auth.Principal, bool) {
	principal, ok := s.currentPrincipal(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return auth.Principal{}, false
	}
	return principal, true
}

func (s *Server) currentPrincipal(r *http.Request) (auth.Principal, bool) {
	if s.authService == nil {
		return auth.Principal{}, false
	}
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil || cookie.Value == "" {
		return auth.Principal{}, false
	}
	principal, _, err := s.authService.Authenticate(r.Context(), cookie.Value)
	if err != nil {
		return auth.Principal{}, false
	}
	return principal, true
}
