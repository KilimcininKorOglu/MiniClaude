package serverhttp

import (
	"net/http"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/device"
)

func (s *Server) deviceStart(w http.ResponseWriter, r *http.Request) {
	if s.deviceService == nil {
		writeError(w, http.StatusServiceUnavailable, "device service is not available")
		return
	}

	var request device.StartRequest
	if r.Body != http.NoBody {
		if err := decodeJSON(w, r, &request); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	response, err := s.deviceService.Start(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (s *Server) deviceApprove(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.authenticateRequest(w, r)
	if !ok {
		return
	}
	if s.deviceService == nil {
		writeError(w, http.StatusServiceUnavailable, "device service is not available")
		return
	}

	var request device.ApproveRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	clientID, err := s.deviceService.Approve(r.Context(), principal, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"client_id": clientID})
}

func (s *Server) devicePoll(w http.ResponseWriter, r *http.Request) {
	if s.deviceService == nil {
		writeError(w, http.StatusServiceUnavailable, "device service is not available")
		return
	}

	var request device.PollRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := s.deviceService.Poll(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, response)
}
