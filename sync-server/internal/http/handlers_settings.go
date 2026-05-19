package serverhttp

import (
	"net/http"

	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
)

func (s *Server) settingsSnapshot(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.authenticateRequest(w, r)
	if !ok {
		return
	}
	if s.syncService == nil {
		writeError(w, http.StatusServiceUnavailable, "sync service is not available")
		return
	}

	snapshot, err := s.syncService.Snapshot(r.Context(), principal.WorkspaceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) settingsPush(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.authenticateRequest(w, r)
	if !ok {
		return
	}
	if s.syncService == nil {
		writeError(w, http.StatusServiceUnavailable, "sync service is not available")
		return
	}

	var request settingssync.PushRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.syncService.Push(r.Context(), principal.WorkspaceID, principal.UserID, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !result.Accepted {
		writeJSON(w, http.StatusConflict, map[string]any{"type": "version_reject", "snapshot": result.Snapshot})
		return
	}
	if s.hub != nil {
		s.hub.Broadcast(principal.WorkspaceID, syncws.Message{Type: "settings_updated", WorkspaceID: principal.WorkspaceID, Version: result.Snapshot.Version, Payload: result.Snapshot})
	}
	writeJSON(w, http.StatusOK, result)
}
