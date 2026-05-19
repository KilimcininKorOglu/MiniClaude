package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type Principal struct {
	UserID      string    `json:"user_id"`
	WorkspaceID string    `json:"workspace_id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type SignupInput struct {
	Email         string
	Username      string
	PasswordHash  string
	WorkspaceName string
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) CreateUserWithWorkspace(ctx context.Context, input SignupInput) (Principal, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	username := strings.TrimSpace(input.Username)
	workspaceName := strings.TrimSpace(input.WorkspaceName)
	if workspaceName == "" {
		workspaceName = username
	}

	if email == "" || username == "" || input.PasswordHash == "" {
		return Principal{}, fmt.Errorf("email, username, and password hash are required")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Principal{}, fmt.Errorf("begin signup transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var userID string
	if err := tx.QueryRow(ctx, `
		insert into users (email, username, password_hash)
		values ($1, $2, $3)
		returning id
	`, email, username, input.PasswordHash).Scan(&userID); err != nil {
		return Principal{}, fmt.Errorf("create user: %w", err)
	}

	slug := slugify(workspaceName)
	if slug == "" {
		slug = "workspace"
	}

	var workspaceID string
	if err := tx.QueryRow(ctx, `
		insert into workspaces (name, slug)
		values ($1, $2)
		returning id
	`, workspaceName, slug).Scan(&workspaceID); err != nil {
		return Principal{}, fmt.Errorf("create workspace: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into workspace_members (workspace_id, user_id, role)
		values ($1, $2, 'owner')
	`, workspaceID, userID); err != nil {
		return Principal{}, fmt.Errorf("create workspace membership: %w", err)
	}

	document := json.RawMessage(`{}`)
	checksum := checksumJSON(document)
	if _, err := tx.Exec(ctx, `
		insert into settings_documents (workspace_id, document, version, checksum, updated_by)
		values ($1, $2, 1, $3, $4)
	`, workspaceID, document, checksum, userID); err != nil {
		return Principal{}, fmt.Errorf("create settings document: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Principal{}, fmt.Errorf("commit signup transaction: %w", err)
	}

	return Principal{UserID: userID, WorkspaceID: workspaceID, Email: email, Username: username, Role: "owner"}, nil
}

func (s *Store) FindUserByLogin(ctx context.Context, login string) (Principal, string, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	row := s.pool.QueryRow(ctx, `
		select u.id, wm.workspace_id, u.email, u.username, wm.role, u.password_hash
		from users u
		join workspace_members wm on wm.user_id = u.id
		where (lower(u.email) = $1 or lower(u.username) = $1) and u.status = 'active'
		order by wm.created_at asc
		limit 1
	`, login)

	principal := Principal{}
	var passwordHash string
	if err := row.Scan(&principal.UserID, &principal.WorkspaceID, &principal.Email, &principal.Username, &principal.Role, &passwordHash); err != nil {
		if err == pgx.ErrNoRows {
			return Principal{}, "", fmt.Errorf("invalid credentials")
		}
		return Principal{}, "", fmt.Errorf("find user: %w", err)
	}
	return principal, passwordHash, nil
}

func (s *Store) CreateWebSession(ctx context.Context, principal Principal, jti string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into web_sessions (user_id, workspace_id, jti, expires_at)
		values ($1, $2, $3, $4)
	`, principal.UserID, principal.WorkspaceID, jti, expiresAt)
	if err != nil {
		return fmt.Errorf("create web session: %w", err)
	}
	return nil
}

func (s *Store) PrincipalForSession(ctx context.Context, jti string) (Principal, error) {
	row := s.pool.QueryRow(ctx, `
		select u.id, ws.workspace_id, u.email, u.username, wm.role, ws.expires_at
		from web_sessions ws
		join users u on u.id = ws.user_id
		join workspace_members wm on wm.workspace_id = ws.workspace_id and wm.user_id = ws.user_id
		where ws.jti = $1 and ws.revoked_at is null and ws.expires_at > now() and u.status = 'active'
	`, jti)

	principal := Principal{}
	if err := row.Scan(&principal.UserID, &principal.WorkspaceID, &principal.Email, &principal.Username, &principal.Role, &principal.ExpiresAt); err != nil {
		if err == pgx.ErrNoRows {
			return Principal{}, fmt.Errorf("session is not active")
		}
		return Principal{}, fmt.Errorf("read web session: %w", err)
	}
	return principal, nil
}

func (s *Store) RevokeWebSession(ctx context.Context, jti string) error {
	_, err := s.pool.Exec(ctx, `
		update web_sessions
		set revoked_at = now()
		where jti = $1 and revoked_at is null
	`, jti)
	if err != nil {
		return fmt.Errorf("revoke web session: %w", err)
	}
	return nil
}

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = slugPattern.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func checksumJSON(document []byte) string {
	sum := sha256.Sum256(document)
	return hex.EncodeToString(sum[:])
}
