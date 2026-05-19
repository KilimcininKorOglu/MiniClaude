package serverhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/providers"
	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
)

type providerView struct {
	ID            string
	Name          string
	BaseURL       string
	Model         string
	DefaultSonnet string
	DefaultHaiku  string
	DefaultOpus   string
	Description   string
	HasSecret     bool
	UpdatedAt     string
}

type sessionView struct {
	ID           string
	ProcessID    string
	Status       string
	ConnectedAt  string
	LastSeenAt   string
	TerminatedAt string
}

type clientView struct {
	ID         string
	Name       string
	MachineID  string
	CreatedAt  string
	LastSeenAt string
	Revoked    bool
	Sessions   []sessionView
}

type auditEventView struct {
	CreatedAt string
	ActorType string
	EventType string
	Metadata  string
}

func (s *Server) providerSaveForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if s.providerStore == nil {
		s.render(w, http.StatusServiceUnavailable, "providers.html", pageData{Title: "Providers", Principal: principal, Error: "Provider store is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, "providers.html", pageData{Title: "Providers", Principal: principal, Error: "Invalid form submission."})
		return
	}

	providerID, err := s.providerStore.Upsert(r.Context(), providers.UpsertInput{
		WorkspaceID:   principal.WorkspaceID,
		Name:          r.FormValue("name"),
		BaseURL:       r.FormValue("base_url"),
		Model:         r.FormValue("model"),
		DefaultSonnet: r.FormValue("default_sonnet"),
		DefaultHaiku:  r.FormValue("default_haiku"),
		DefaultOpus:   r.FormValue("default_opus"),
		Description:   r.FormValue("description"),
		Secret:        r.FormValue("secret"),
	})
	if err != nil {
		s.renderProvidersError(w, r, principal, err)
		return
	}
	_ = s.bumpSettingsVersion(r.Context(), principal.WorkspaceID, principal.UserID, "provider_saved", map[string]string{"provider_id": providerID, "name": strings.TrimSpace(r.FormValue("name"))})
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "provider_saved", map[string]string{"provider_id": providerID, "name": strings.TrimSpace(r.FormValue("name"))})
	s.broadcastSnapshot(r.Context(), principal.WorkspaceID)
	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

func (s *Server) providerDeleteForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if s.providerStore == nil {
		s.render(w, http.StatusServiceUnavailable, "providers.html", pageData{Title: "Providers", Principal: principal, Error: "Provider store is not available."})
		return
	}
	if err := r.ParseForm(); err != nil {
		s.renderProvidersError(w, r, principal, fmt.Errorf("invalid form submission"))
		return
	}
	providerID := r.FormValue("provider_id")
	if err := s.providerStore.Delete(r.Context(), principal.WorkspaceID, providerID); err != nil {
		s.renderProvidersError(w, r, principal, err)
		return
	}
	_ = s.bumpSettingsVersion(r.Context(), principal.WorkspaceID, principal.UserID, "provider_deleted", map[string]string{"provider_id": providerID})
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "provider_deleted", map[string]string{"provider_id": providerID})
	s.broadcastSnapshot(r.Context(), principal.WorkspaceID)
	http.Redirect(w, r, "/providers", http.StatusSeeOther)
}

func (s *Server) settingsPage(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	data, err := s.settingsPageData(r.Context(), principal, "", "")
	if err != nil {
		s.render(w, http.StatusBadRequest, "settings.html", pageData{Title: "Settings", Principal: principal, Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, "settings.html", data)
}

func (s *Server) settingsSaveForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, "settings.html", pageData{Title: "Settings", Principal: principal, Error: "Invalid form submission."})
		return
	}
	document := json.RawMessage(strings.TrimSpace(r.FormValue("document")))
	if !json.Valid(document) {
		data, _ := s.settingsPageData(r.Context(), principal, "Settings document must be valid JSON.", "")
		data.SettingsDocument = string(document)
		s.render(w, http.StatusBadRequest, "settings.html", data)
		return
	}
	var baseVersion int64
	if _, err := fmt.Sscanf(r.FormValue("base_version"), "%d", &baseVersion); err != nil {
		s.render(w, http.StatusBadRequest, "settings.html", pageData{Title: "Settings", Principal: principal, Error: "Settings version is invalid."})
		return
	}
	result, err := s.syncService.Push(r.Context(), principal.WorkspaceID, principal.UserID, settingssync.PushRequest{BaseVersion: baseVersion, Document: document})
	if err != nil {
		data, _ := s.settingsPageData(r.Context(), principal, err.Error(), "")
		s.render(w, http.StatusBadRequest, "settings.html", data)
		return
	}
	if !result.Accepted {
		data, _ := s.settingsPageData(r.Context(), principal, "Settings changed elsewhere. Review the latest version and save again.", "")
		s.render(w, http.StatusConflict, "settings.html", data)
		return
	}
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "settings_saved", map[string]any{"version": result.Snapshot.Version})
	s.broadcastSnapshot(r.Context(), principal.WorkspaceID)
	data, err := s.settingsPageData(r.Context(), principal, "", "Settings saved.")
	if err != nil {
		s.render(w, http.StatusBadRequest, "settings.html", pageData{Title: "Settings", Principal: principal, Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, "settings.html", data)
}

func (s *Server) clientRevokeForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		s.renderClientsError(w, r, principal, fmt.Errorf("invalid form submission"))
		return
	}
	clientID := r.FormValue("client_id")
	_, err := s.pool.Exec(r.Context(), `
		update clients set revoked_at = now()
		where workspace_id = $1 and id = $2 and revoked_at is null
	`, principal.WorkspaceID, clientID)
	if err != nil {
		s.renderClientsError(w, r, principal, fmt.Errorf("revoke client: %w", err))
		return
	}
	_, _ = s.pool.Exec(r.Context(), `
		update client_sessions set status = 'terminated', terminated_at = now()
		where workspace_id = $1 and client_id = $2 and status = 'active'
	`, principal.WorkspaceID, clientID)
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "client_revoked", map[string]string{"client_id": clientID})
	if s.hub != nil {
		s.hub.Broadcast(principal.WorkspaceID, syncws.Message{Type: "terminate_session", WorkspaceID: principal.WorkspaceID, Payload: map[string]string{"client_id": clientID}})
	}
	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

func (s *Server) clientSessionTerminateForm(w http.ResponseWriter, r *http.Request) {
	principal, ok := s.requireHTMLAuth(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		s.renderClientsError(w, r, principal, fmt.Errorf("invalid form submission"))
		return
	}
	sessionID := r.FormValue("session_id")
	var clientID string
	err := s.pool.QueryRow(r.Context(), `
		update client_sessions set status = 'terminated', terminated_at = now()
		where workspace_id = $1 and id = $2 and status = 'active'
		returning client_id::text
	`, principal.WorkspaceID, sessionID).Scan(&clientID)
	if err != nil {
		s.renderClientsError(w, r, principal, fmt.Errorf("terminate session: %w", err))
		return
	}
	_ = s.audit(r.Context(), principal.WorkspaceID, principal.UserID, "user", "client_session_terminated", map[string]string{"client_id": clientID, "session_id": sessionID})
	if s.hub != nil {
		s.hub.Broadcast(principal.WorkspaceID, syncws.Message{Type: "terminate_session", WorkspaceID: principal.WorkspaceID, Payload: map[string]string{"client_id": clientID, "session_id": sessionID}})
	}
	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

func (s *Server) providerViews(ctx context.Context, workspaceID string) ([]providerView, error) {
	if s.providerStore == nil {
		return nil, fmt.Errorf("provider store is not available")
	}
	items, err := s.providerStore.List(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	views := make([]providerView, 0, len(items))
	for _, item := range items {
		views = append(views, providerView{
			ID: item.ID, Name: item.Name, BaseURL: item.BaseURL, Model: item.Model,
			DefaultSonnet: item.DefaultSonnet, DefaultHaiku: item.DefaultHaiku, DefaultOpus: item.DefaultOpus,
			Description: item.Description, HasSecret: item.HasSecret, UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}
	return views, nil
}

func (s *Server) clientViews(ctx context.Context, workspaceID string) ([]clientView, error) {
	rows, err := s.pool.Query(ctx, `
		select id::text, name, coalesce(machine_id, ''), to_char(created_at, 'YYYY-MM-DD HH24:MI'), coalesce(to_char(last_seen_at, 'YYYY-MM-DD HH24:MI'), ''), revoked_at is not null
		from clients
		where workspace_id = $1
		order by created_at desc
	`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}
	defer rows.Close()

	clients := []clientView{}
	for rows.Next() {
		client := clientView{}
		if err := rows.Scan(&client.ID, &client.Name, &client.MachineID, &client.CreatedAt, &client.LastSeenAt, &client.Revoked); err != nil {
			return nil, fmt.Errorf("scan client: %w", err)
		}
		client.Sessions, err = s.sessionViews(ctx, workspaceID, client.ID)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate clients: %w", err)
	}
	return clients, nil
}

func (s *Server) sessionViews(ctx context.Context, workspaceID, clientID string) ([]sessionView, error) {
	rows, err := s.pool.Query(ctx, `
		select id::text, coalesce(process_id, ''), status, to_char(connected_at, 'YYYY-MM-DD HH24:MI'), to_char(last_seen_at, 'YYYY-MM-DD HH24:MI'), coalesce(to_char(terminated_at, 'YYYY-MM-DD HH24:MI'), '')
		from client_sessions
		where workspace_id = $1 and client_id = $2
		order by connected_at desc
	`, workspaceID, clientID)
	if err != nil {
		return nil, fmt.Errorf("list client sessions: %w", err)
	}
	defer rows.Close()

	sessions := []sessionView{}
	for rows.Next() {
		session := sessionView{}
		if err := rows.Scan(&session.ID, &session.ProcessID, &session.Status, &session.ConnectedAt, &session.LastSeenAt, &session.TerminatedAt); err != nil {
			return nil, fmt.Errorf("scan client session: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate client sessions: %w", err)
	}
	return sessions, nil
}

func (s *Server) settingsPageData(ctx context.Context, principal auth.Principal, errorMessage, notice string) (pageData, error) {
	var document json.RawMessage
	var version int64
	var checksum string
	if err := s.pool.QueryRow(ctx, `
		select document, version, checksum
		from settings_documents
		where workspace_id = $1
	`, principal.WorkspaceID).Scan(&document, &version, &checksum); err != nil {
		return pageData{}, fmt.Errorf("read settings document: %w", err)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, document, "", "  "); err != nil {
		return pageData{}, fmt.Errorf("format settings document: %w", err)
	}
	return pageData{Title: "Settings", Principal: principal, Error: errorMessage, Notice: notice, SettingsDocument: pretty.String(), SettingsVersion: version, SettingsChecksum: checksum}, nil
}

func (s *Server) auditEvents(ctx context.Context, workspaceID string) ([]auditEventView, error) {
	rows, err := s.pool.Query(ctx, `
		select to_char(created_at, 'YYYY-MM-DD HH24:MI'), actor_type, event_type, metadata::text
		from audit_events
		where workspace_id = $1
		order by created_at desc
		limit 20
	`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	defer rows.Close()

	events := []auditEventView{}
	for rows.Next() {
		event := auditEventView{}
		if err := rows.Scan(&event.CreatedAt, &event.ActorType, &event.EventType, &event.Metadata); err != nil {
			return nil, fmt.Errorf("scan audit event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit events: %w", err)
	}
	return events, nil
}

func (s *Server) bumpSettingsVersion(ctx context.Context, workspaceID, userID, eventType string, metadata any) error {
	payload, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin settings version transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var version int64
	if err := tx.QueryRow(ctx, `
		update settings_documents
		set version = version + 1, updated_by = $2, updated_at = now()
		where workspace_id = $1
		returning version
	`, workspaceID, userID).Scan(&version); err != nil {
		return fmt.Errorf("bump settings version: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		insert into settings_events (workspace_id, version, event_type, payload, created_by)
		values ($1, $2, $3, $4, $5)
	`, workspaceID, version, eventType, payload, userID); err != nil {
		return fmt.Errorf("append settings event: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settings version transaction: %w", err)
	}
	return nil
}

func (s *Server) audit(ctx context.Context, workspaceID, userID, actorType, eventType string, metadata any) error {
	payload, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `
		insert into audit_events (workspace_id, user_id, actor_type, event_type, metadata)
		values ($1, $2, $3, $4, $5)
	`, workspaceID, userID, actorType, eventType, payload)
	return err
}

func (s *Server) broadcastSnapshot(ctx context.Context, workspaceID string) {
	if s.hub == nil || s.syncService == nil {
		return
	}
	snapshot, err := s.syncService.Snapshot(ctx, workspaceID)
	if err != nil {
		return
	}
	s.hub.Broadcast(workspaceID, syncws.Message{Type: "settings_updated", WorkspaceID: workspaceID, Version: snapshot.Version, Payload: snapshot})
}

func (s *Server) renderProvidersError(w http.ResponseWriter, r *http.Request, principal auth.Principal, err error) {
	providers, _ := s.providerViews(r.Context(), principal.WorkspaceID)
	s.render(w, http.StatusBadRequest, "providers.html", pageData{Title: "Providers", Principal: principal, Providers: providers, Error: err.Error()})
}

func (s *Server) renderClientsError(w http.ResponseWriter, r *http.Request, principal auth.Principal, err error) {
	clients, _ := s.clientViews(r.Context(), principal.WorkspaceID)
	s.render(w, http.StatusBadRequest, "clients.html", pageData{Title: "Clients", Principal: principal, Clients: clients, Error: err.Error()})
}
