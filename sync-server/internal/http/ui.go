package serverhttp

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

//go:embed web/templates/*.html web/static/*
var uiFiles embed.FS

type pageData struct {
	Title            string
	Principal        auth.Principal
	Error            string
	Notice           string
	Next             string
	UserCode         string
	ClientID         string
	Providers        []providerView
	Clients          []clientView
	SettingsDocument string
	SettingsVersion  int64
	SettingsChecksum string
	AuditEvents      []auditEventView
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
	next := safeNextPath(r.URL.Query().Get("next"))
	if _, ok := s.currentPrincipal(r); ok {
		if next == "" {
			next = "/dashboard"
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}
	s.render(w, http.StatusOK, "login.html", pageData{Title: "Login", Next: next})
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
	events, err := s.auditEvents(r.Context(), principal.WorkspaceID)
	if err != nil {
		s.render(w, http.StatusBadRequest, "dashboard.html", pageData{Title: "Dashboard", Principal: principal, Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, "dashboard.html", pageData{Title: "Dashboard", Principal: principal, AuditEvents: events})
}

func (s *Server) providersPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	providers, err := s.providerViews(r.Context(), principal.WorkspaceID)
	if err != nil {
		s.render(w, http.StatusBadRequest, "providers.html", pageData{Title: "Providers", Principal: principal, Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, "providers.html", pageData{Title: "Providers", Principal: principal, Providers: providers})
}

func (s *Server) clientsPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	clients, err := s.clientViews(r.Context(), principal.WorkspaceID)
	if err != nil {
		s.render(w, http.StatusBadRequest, "clients.html", pageData{Title: "Clients", Principal: principal, Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, "clients.html", pageData{Title: "Clients", Principal: principal, Clients: clients})
}

func (s *Server) devicePage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.currentPrincipal(r)
	if !ok {
		http.Redirect(w, r, "/login?next="+url.QueryEscape(r.URL.RequestURI()), http.StatusSeeOther)
		return
	}
	s.render(w, http.StatusOK, "device.html", pageData{Title: "Link device", Principal: principal, UserCode: r.URL.Query().Get("user_code")})
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

func safeNextPath(next string) string {
	if len(next) == 0 || next[0] != '/' || (len(next) > 1 && next[1] == '/') {
		return ""
	}
	return next
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
