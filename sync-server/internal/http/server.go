package serverhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/device"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/providers"
	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool          *pgxpool.Pool
	mux           *http.ServeMux
	authService   *auth.Service
	deviceService *device.Service
	syncService   *settingssync.Service
	providerStore *providers.Store
	hub           *syncws.Hub
	cookies       auth.CookieConfig
}

func New(pool *pgxpool.Pool, authService *auth.Service, deviceService *device.Service, syncService *settingssync.Service, providerStore *providers.Store, hub *syncws.Hub, cookies auth.CookieConfig) http.Handler {
	server := &Server{pool: pool, mux: http.NewServeMux(), authService: authService, deviceService: deviceService, syncService: syncService, providerStore: providerStore, hub: hub, cookies: cookies}
	server.routes()
	return server.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /", s.homePage)
	s.mux.HandleFunc("GET /favicon.ico", s.favicon)
	s.mux.HandleFunc("GET /login", s.loginPage)
	s.mux.HandleFunc("GET /signup", s.signupPage)
	s.mux.HandleFunc("GET /dashboard", s.dashboardPage)
	s.mux.HandleFunc("GET /providers", s.providersPage)
	s.mux.HandleFunc("POST /ui/providers", s.providerSaveForm)
	s.mux.HandleFunc("POST /ui/providers/delete", s.providerDeleteForm)
	s.mux.HandleFunc("GET /settings", s.settingsPage)
	s.mux.HandleFunc("POST /ui/settings", s.settingsSaveForm)
	s.mux.HandleFunc("GET /clients", s.clientsPage)
	s.mux.HandleFunc("POST /ui/clients/revoke", s.clientRevokeForm)
	s.mux.HandleFunc("POST /ui/client-sessions/terminate", s.clientSessionTerminateForm)
	s.mux.HandleFunc("GET /device", s.devicePage)
	s.mux.HandleFunc("POST /device", s.deviceApproveForm)
	s.mux.HandleFunc("GET /web/static/", s.staticFile)
	s.mux.HandleFunc("POST /ui/login", s.loginForm)
	s.mux.HandleFunc("POST /ui/signup", s.signupForm)
	s.mux.HandleFunc("POST /ui/logout", s.logoutForm)
	s.mux.HandleFunc("POST /auth/signup", s.signup)
	s.mux.HandleFunc("POST /auth/login", s.login)
	s.mux.HandleFunc("POST /auth/logout", s.logout)
	s.mux.HandleFunc("GET /me", s.me)
	s.mux.HandleFunc("POST /api/device/start", s.deviceStart)
	s.mux.HandleFunc("POST /api/device/approve", s.deviceApprove)
	s.mux.HandleFunc("POST /api/device/poll", s.devicePoll)
	s.mux.HandleFunc("GET /api/settings/snapshot", s.settingsSnapshot)
	s.mux.HandleFunc("POST /api/settings/push", s.settingsPush)
	s.mux.HandleFunc("GET /api/sync/ws", s.syncWebSocket)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unhealthy"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
