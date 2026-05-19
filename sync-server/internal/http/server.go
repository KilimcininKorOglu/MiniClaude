package serverhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool        *pgxpool.Pool
	mux         *http.ServeMux
	authService *auth.Service
	cookies     auth.CookieConfig
}

func New(pool *pgxpool.Pool, authService *auth.Service, cookies auth.CookieConfig) http.Handler {
	server := &Server{pool: pool, mux: http.NewServeMux(), authService: authService, cookies: cookies}
	server.routes()
	return server.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("POST /auth/signup", s.signup)
	s.mux.HandleFunc("POST /auth/login", s.login)
	s.mux.HandleFunc("POST /auth/logout", s.logout)
	s.mux.HandleFunc("GET /me", s.me)
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
