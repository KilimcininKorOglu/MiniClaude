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
	_ = s.audit(r.Context(), "", "", "client", "device_start", map[string]string{"client_name": request.ClientName})
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
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "device_approved", map[string]string{"client_id": clientID, "surface": "api"})
	writeJSON(w, http.StatusOK, map[string]string{"client_id": clientID})
}

func (s *Server) deviceApproveForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if s.deviceService == nil {
		s.render(w, r, http.StatusServiceUnavailable, "device.html", pageData{Title: "Link device", Principal: principal, Error: "Device service is not available."})
		return
	}
	if !s.requireCSRF(w, r) {
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, r, http.StatusBadRequest, "device.html", pageData{Title: "Link device", Principal: principal, Error: "Invalid form submission."})
		return
	}

	clientID, err := s.deviceService.Approve(r.Context(), principal, device.ApproveRequest{UserCode: r.FormValue("user_code"), ClientName: r.FormValue("client_name")})
	if err != nil {
		s.render(w, r, http.StatusBadRequest, "device.html", pageData{Title: "Link device", Principal: principal, Error: err.Error(), UserCode: r.FormValue("user_code")})
		return
	}
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "device_approved", map[string]string{"client_id": clientID, "surface": "web"})
	s.render(w, r, http.StatusOK, "device.html", pageData{Title: "Link device", Principal: principal, ClientID: clientID})
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
	if response.Status == "approved" {
		_ = s.audit(r.Context(), response.WorkspaceID, "", "client", "device_poll_consumed", map[string]string{"client_id": response.ClientID})
	}
	writeJSON(w, http.StatusOK, response)
}
