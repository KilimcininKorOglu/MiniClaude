package device

import (
	"testing"
	"time"
)

func TestClientTokenManagerIssueParse(t *testing.T) {
	manager, err := NewClientTokenManager("01234567890123456789012345678901", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	issuedAt := time.Now().UTC()
	token, expiresAt, err := manager.Issue("client-id", "workspace-id", issuedAt)
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
	if claims.ClientID != "client-id" || claims.WorkspaceID != "workspace-id" || claims.Subject != "client-id" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestRandomUserCodeFormat(t *testing.T) {
	code, err := randomUserCode()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 9 || code[4] != '-' {
		t.Fatalf("unexpected user code format: %q", code)
	}
}
