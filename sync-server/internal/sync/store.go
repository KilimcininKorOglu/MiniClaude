package sync

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type Snapshot struct {
	WorkspaceID string          `json:"workspace_id"`
	Document    json.RawMessage `json:"document"`
	Version     int64           `json:"version"`
	Checksum    string          `json:"checksum"`
}

type PushResult struct {
	Accepted bool     `json:"accepted"`
	Snapshot Snapshot `json:"snapshot"`
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Snapshot(ctx context.Context, workspaceID string) (Snapshot, error) {
	snapshot := Snapshot{WorkspaceID: workspaceID}
	err := s.pool.QueryRow(ctx, `
		select document, version, checksum
		from settings_documents
		where workspace_id = $1
	`, workspaceID).Scan(&snapshot.Document, &snapshot.Version, &snapshot.Checksum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Snapshot{}, fmt.Errorf("settings document not found")
		}
		return Snapshot{}, fmt.Errorf("read settings snapshot: %w", err)
	}
	return snapshot, nil
}

func (s *Store) Push(ctx context.Context, workspaceID, userID string, baseVersion int64, document json.RawMessage) (PushResult, error) {
	if !json.Valid(document) {
		return PushResult{}, fmt.Errorf("settings document must be valid JSON")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return PushResult{}, fmt.Errorf("begin settings push transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	current := Snapshot{WorkspaceID: workspaceID}
	if err := tx.QueryRow(ctx, `
		select document, version, checksum
		from settings_documents
		where workspace_id = $1
		for update
	`, workspaceID).Scan(&current.Document, &current.Version, &current.Checksum); err != nil {
		if err == pgx.ErrNoRows {
			return PushResult{}, fmt.Errorf("settings document not found")
		}
		return PushResult{}, fmt.Errorf("read settings snapshot: %w", err)
	}

	if current.Version != baseVersion {
		return PushResult{Accepted: false, Snapshot: current}, nil
	}

	next := Snapshot{
		WorkspaceID: workspaceID,
		Document:    document,
		Version:     current.Version + 1,
		Checksum:    checksum(document),
	}
	if _, err := tx.Exec(ctx, `
		update settings_documents
		set document = $2, version = $3, checksum = $4, updated_by = $5, updated_at = now()
		where workspace_id = $1
	`, workspaceID, next.Document, next.Version, next.Checksum, userID); err != nil {
		return PushResult{}, fmt.Errorf("update settings document: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		insert into settings_events (workspace_id, version, event_type, payload, created_by)
		values ($1, $2, 'settings_push', $3, $4)
	`, workspaceID, next.Version, next.Document, userID); err != nil {
		return PushResult{}, fmt.Errorf("append settings event: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return PushResult{}, fmt.Errorf("commit settings push transaction: %w", err)
	}
	return PushResult{Accepted: true, Snapshot: next}, nil
}

func checksum(document []byte) string {
	sum := sha256.Sum256(document)
	return hex.EncodeToString(sum[:])
}
