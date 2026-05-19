package device

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type Request struct {
	DeviceCode string
	UserCode   string
	ExpiresAt  time.Time
}

type PollResult struct {
	Status      string
	ClientID    string
	WorkspaceID string
}

type ClientPrincipal struct {
	ClientID    string
	WorkspaceID string
	UserID      string
	Name        string
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) CreateRequest(ctx context.Context, request Request) error {
	_, err := s.pool.Exec(ctx, `
		insert into device_login_requests (device_code, user_code, expires_at)
		values ($1, $2, $3)
	`, request.DeviceCode, request.UserCode, request.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create device login request: %w", err)
	}
	return nil
}

func (s *Store) Approve(ctx context.Context, userCode string, principal auth.Principal, clientName string) (string, error) {
	userCode = normalizeUserCode(userCode)
	clientName = strings.TrimSpace(clientName)
	if clientName == "" {
		clientName = "MiniClaude Client"
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin device approval transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var requestID string
	var status string
	var expiresAt time.Time
	if err := tx.QueryRow(ctx, `
		select id, status, expires_at
		from device_login_requests
		where user_code = $1
		for update
	`, userCode).Scan(&requestID, &status, &expiresAt); err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("device request not found")
		}
		return "", fmt.Errorf("read device request: %w", err)
	}
	if time.Now().UTC().After(expiresAt) {
		_, _ = tx.Exec(ctx, `update device_login_requests set status = 'expired' where id = $1`, requestID)
		return "", fmt.Errorf("device request expired")
	}
	if status != "pending" {
		return "", fmt.Errorf("device request is %s", status)
	}

	var clientID string
	if err := tx.QueryRow(ctx, `
		insert into clients (workspace_id, user_id, name)
		values ($1, $2, $3)
		returning id
	`, principal.WorkspaceID, principal.UserID, clientName).Scan(&clientID); err != nil {
		return "", fmt.Errorf("create client: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		update device_login_requests
		set workspace_id = $1, approved_by = $2, client_id = $3, status = 'approved', approved_at = now()
		where id = $4
	`, principal.WorkspaceID, principal.UserID, clientID, requestID); err != nil {
		return "", fmt.Errorf("approve device request: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit device approval transaction: %w", err)
	}
	return clientID, nil
}

func (s *Store) Poll(ctx context.Context, deviceCode string) (PollResult, error) {
	deviceCode = strings.TrimSpace(deviceCode)
	if deviceCode == "" {
		return PollResult{}, fmt.Errorf("device code is required")
	}

	result := PollResult{}
	var expiresAt time.Time
	err := s.pool.QueryRow(ctx, `
		select status, coalesce(client_id::text, ''), coalesce(workspace_id::text, ''), expires_at
		from device_login_requests
		where device_code = $1
	`, deviceCode).Scan(&result.Status, &result.ClientID, &result.WorkspaceID, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return PollResult{}, fmt.Errorf("device request not found")
		}
		return PollResult{}, fmt.Errorf("poll device request: %w", err)
	}
	if result.Status == "pending" && time.Now().UTC().After(expiresAt) {
		_, _ = s.pool.Exec(ctx, `update device_login_requests set status = 'expired' where device_code = $1`, deviceCode)
		result.Status = "expired"
	}
	return result, nil
}

func (s *Store) ClientPrincipal(ctx context.Context, clientID, workspaceID string) (ClientPrincipal, error) {
	principal := ClientPrincipal{}
	err := s.pool.QueryRow(ctx, `
		select id, workspace_id, user_id, name
		from clients
		where id = $1 and workspace_id = $2 and revoked_at is null
	`, clientID, workspaceID).Scan(&principal.ClientID, &principal.WorkspaceID, &principal.UserID, &principal.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ClientPrincipal{}, fmt.Errorf("client is not active")
		}
		return ClientPrincipal{}, fmt.Errorf("read client principal: %w", err)
	}
	return principal, nil
}

func normalizeUserCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "")
	if len(value) == 8 && !strings.Contains(value, "-") {
		return value[:4] + "-" + value[4:]
	}
	return value
}
