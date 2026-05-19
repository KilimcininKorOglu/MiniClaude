package serverhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type websocketMessage struct {
	Type        string          `json:"type"`
	BaseVersion int64           `json:"base_version,omitempty"`
	Document    json.RawMessage `json:"document,omitempty"`
}

func (s *Server) syncWebSocket(w http.ResponseWriter, r *http.Request) {
	if s.deviceService == nil || s.syncService == nil || s.hub == nil {
		writeError(w, http.StatusServiceUnavailable, "sync websocket is not available")
		return
	}

	principal, err := s.deviceService.AuthenticateClient(r.Context(), clientAccessToken(r))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid client token")
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sessionID, err := s.createClientSession(r.Context(), principal.WorkspaceID, principal.ClientID)
	if err != nil {
		_ = conn.WriteJSON(syncws.Message{Type: "error", Payload: err.Error()})
		return
	}
	defer s.closeClientSession(contextWithoutCancel(r.Context()), principal.WorkspaceID, sessionID)

	var writeMu sync.Mutex
	writeMessage := func(value any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return conn.WriteJSON(value)
	}

	if err := writeMessage(syncws.Message{Type: "session_started", WorkspaceID: principal.WorkspaceID, Payload: map[string]string{"client_id": principal.ClientID, "session_id": sessionID}}); err != nil {
		return
	}

	snapshot, err := s.syncService.Snapshot(r.Context(), principal.WorkspaceID)
	if err != nil {
		_ = writeMessage(syncws.Message{Type: "error", Payload: err.Error()})
		return
	}
	if err := writeMessage(syncws.Message{Type: "snapshot", WorkspaceID: principal.WorkspaceID, Version: snapshot.Version, Payload: snapshot}); err != nil {
		return
	}

	updates, unsubscribe := s.hub.Subscribe(principal.WorkspaceID)
	defer unsubscribe()
	closed := make(chan struct{})
	go func() {
		for {
			select {
			case update, ok := <-updates:
				if !ok {
					return
				}
				if err := writeMessage(update); err != nil {
					_ = conn.Close()
					return
				}
			case <-closed:
				return
			}
		}
	}()
	defer close(closed)

	for {
		var message websocketMessage
		if err := conn.ReadJSON(&message); err != nil {
			return
		}
		s.touchClientSession(r.Context(), principal.WorkspaceID, principal.ClientID, sessionID)
		switch message.Type {
		case "settings_push":
			result, err := s.syncService.Push(r.Context(), principal.WorkspaceID, principal.UserID, settingssync.PushRequest{BaseVersion: message.BaseVersion, Document: message.Document})
			if err != nil {
				_ = writeMessage(syncws.Message{Type: "error", Payload: err.Error()})
				continue
			}
			if !result.Accepted {
				_ = writeMessage(syncws.Message{Type: "version_reject", WorkspaceID: principal.WorkspaceID, Version: result.Snapshot.Version, Payload: result.Snapshot})
				continue
			}
			s.hub.Broadcast(principal.WorkspaceID, syncws.Message{Type: "settings_updated", WorkspaceID: principal.WorkspaceID, Version: result.Snapshot.Version, Payload: result.Snapshot})
		case "ping":
			_ = writeMessage(syncws.Message{Type: "pong", WorkspaceID: principal.WorkspaceID, Payload: map[string]string{"session_id": sessionID}})
		default:
			_ = writeMessage(syncws.Message{Type: "error", Payload: fmt.Sprintf("unknown message type %q", message.Type)})
		}
	}
}

func (s *Server) createClientSession(ctx context.Context, workspaceID, clientID string) (string, error) {
	var sessionID string
	err := s.pool.QueryRow(ctx, `
		insert into client_sessions (client_id, workspace_id, process_id)
		values ($1, $2, $3)
		returning id::text
	`, clientID, workspaceID, "websocket").Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("create client session: %w", err)
	}
	return sessionID, nil
}

func (s *Server) touchClientSession(ctx context.Context, workspaceID, clientID, sessionID string) {
	_, _ = s.pool.Exec(ctx, `
		update clients set last_seen_at = now()
		where workspace_id = $1 and id = $2
	`, workspaceID, clientID)
	_, _ = s.pool.Exec(ctx, `
		update client_sessions set last_seen_at = now()
		where workspace_id = $1 and id = $2 and status = 'active'
	`, workspaceID, sessionID)
}

func (s *Server) closeClientSession(ctx context.Context, workspaceID, sessionID string) {
	deadline, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, _ = s.pool.Exec(deadline, `
		update client_sessions set status = 'disconnected'
		where workspace_id = $1 and id = $2 and status = 'active'
	`, workspaceID, sessionID)
}

func contextWithoutCancel(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}

func clientAccessToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}
	return r.URL.Query().Get("access_token")
}
