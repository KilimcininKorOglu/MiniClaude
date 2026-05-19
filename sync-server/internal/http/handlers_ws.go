package serverhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

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

	var writeMu sync.Mutex
	writeMessage := func(value any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return conn.WriteJSON(value)
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
			_ = writeMessage(syncws.Message{Type: "pong", WorkspaceID: principal.WorkspaceID})
		default:
			_ = writeMessage(syncws.Message{Type: "error", Payload: fmt.Sprintf("unknown message type %q", message.Type)})
		}
	}
}

func clientAccessToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}
	return r.URL.Query().Get("access_token")
}
