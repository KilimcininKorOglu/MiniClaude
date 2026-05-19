package serverhttp

import (
	"encoding/json"
	"net/http"

	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
)

type settingsPushResponse struct {
	Type        string          `json:"type"`
	Accepted    bool            `json:"accepted"`
	WorkspaceID string          `json:"workspace_id"`
	Version     int64           `json:"version"`
	Checksum    string          `json:"checksum"`
	Document    json.RawMessage `json:"document"`
}

func settingsMessage(messageType string, snapshot settingssync.Snapshot) syncws.Message {
	return syncws.Message{Type: messageType, WorkspaceID: snapshot.WorkspaceID, Version: snapshot.Version, Payload: snapshot}
}

func settingsResponse(messageType string, accepted bool, snapshot settingssync.Snapshot) settingsPushResponse {
	return settingsPushResponse{Type: messageType, Accepted: accepted, WorkspaceID: snapshot.WorkspaceID, Version: snapshot.Version, Checksum: snapshot.Checksum, Document: snapshot.Document}
}

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
		writeJSON(w, http.StatusConflict, settingsResponse("version_reject", false, result.Snapshot))
		return
	}
	if s.hub != nil {
		s.hub.Broadcast(principal.WorkspaceID, settingsMessage("settings_updated", result.Snapshot))
	}
	writeJSON(w, http.StatusOK, settingsResponse("settings_applied", true, result.Snapshot))
}
