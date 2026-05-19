package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID      string `json:"user_id"`
	WorkspaceID string `json:"workspace_id"`
	jwt.RegisteredClaims
}

func NewTokenManager(secret string, ttl time.Duration) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT secret must be at least 32 bytes")
	}
	return &TokenManager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *TokenManager) Issue(userID, workspaceID, jti string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(m.ttl)
	claims := Claims{
		UserID:      userID,
		WorkspaceID: workspaceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}
	return signed, expiresAt, nil
}

func (m *TokenManager) Parse(tokenValue string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parse jwt: %w", err)
	}
	if !token.Valid {
		return Claims{}, fmt.Errorf("invalid jwt")
	}
	if claims.ID == "" || claims.UserID == "" || claims.WorkspaceID == "" {
		return Claims{}, fmt.Errorf("missing jwt claims")
	}
	return claims, nil
}
