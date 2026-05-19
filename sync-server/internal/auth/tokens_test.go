package auth

import (
	"testing"
	"time"
)

func TestTokenManagerIssueParse(t *testing.T) {
	manager, err := NewTokenManager("01234567890123456789012345678901", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	issuedAt := time.Now().UTC()
	token, expiresAt, err := manager.Issue("user-id", "workspace-id", "session-id", issuedAt)
	if err != nil {
		t.Fatal(err)
	}
	if !expiresAt.After(issuedAt) {
		t.Fatal("expected token expiry after issue time")
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-id" || claims.WorkspaceID != "workspace-id" || claims.ID != "session-id" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
