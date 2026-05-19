package device

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClientTokenManager struct {
	secret []byte
	ttl    time.Duration
}

type ClientClaims struct {
	ClientID    string `json:"client_id"`
	WorkspaceID string `json:"workspace_id"`
	jwt.RegisteredClaims
}

func NewClientTokenManager(secret string, ttl time.Duration) (*ClientTokenManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("client token secret must be at least 32 bytes")
	}
	return &ClientTokenManager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *ClientTokenManager) Issue(clientID, workspaceID string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(m.ttl)
	claims := ClientClaims{
		ClientID:    clientID,
		WorkspaceID: workspaceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   clientID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign client jwt: %w", err)
	}
	return signed, expiresAt, nil
}

func (m *ClientTokenManager) Parse(tokenValue string) (ClientClaims, error) {
	claims := ClientClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return ClientClaims{}, fmt.Errorf("parse client jwt: %w", err)
	}
	if !token.Valid {
		return ClientClaims{}, fmt.Errorf("invalid client jwt")
	}
	if claims.ClientID == "" || claims.WorkspaceID == "" {
		return ClientClaims{}, fmt.Errorf("missing client jwt claims")
	}
	return claims, nil
}
