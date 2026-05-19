package sync

import (
	"context"
	"encoding/json"
)

type Service struct {
	store *Store
}

type PushRequest struct {
	BaseVersion int64           `json:"base_version"`
	Document    json.RawMessage `json:"document"`
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Snapshot(ctx context.Context, workspaceID string) (Snapshot, error) {
	return s.store.Snapshot(ctx, workspaceID)
}

func (s *Service) Push(ctx context.Context, workspaceID, userID string, request PushRequest) (PushResult, error) {
	return s.store.Push(ctx, workspaceID, userID, request.BaseVersion, request.Document)
}
