package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	store  *Store
	tokens *TokenManager
	now    func() time.Time
}

type SignupRequest struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	WorkspaceName string `json:"workspace_name"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Session struct {
	Principal Principal
	Token     string
	JTI       string
	ExpiresAt time.Time
}

func NewService(store *Store, tokens *TokenManager) *Service {
	return &Service{store: store, tokens: tokens, now: time.Now}
}

func (s *Service) Signup(ctx context.Context, request SignupRequest) (Session, error) {
	if err := validateSignup(request); err != nil {
		return Session{}, err
	}

	passwordHash, err := HashPassword(request.Password)
	if err != nil {
		return Session{}, err
	}

	principal, err := s.store.CreateUserWithWorkspace(ctx, SignupInput{
		Email:         request.Email,
		Username:      request.Username,
		PasswordHash:  passwordHash,
		WorkspaceName: request.WorkspaceName,
	})
	if err != nil {
		return Session{}, err
	}

	return s.createSession(ctx, principal)
}

func (s *Service) Login(ctx context.Context, request LoginRequest) (Session, error) {
	if strings.TrimSpace(request.Login) == "" || request.Password == "" {
		return Session{}, fmt.Errorf("login and password are required")
	}

	principal, passwordHash, err := s.store.FindUserByLogin(ctx, request.Login)
	if err != nil {
		return Session{}, err
	}

	ok, err := VerifyPassword(request.Password, passwordHash)
	if err != nil {
		return Session{}, err
	}
	if !ok {
		return Session{}, fmt.Errorf("invalid credentials")
	}

	return s.createSession(ctx, principal)
}

func (s *Service) Authenticate(ctx context.Context, token string) (Principal, string, error) {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return Principal{}, "", err
	}
	principal, err := s.store.PrincipalForSession(ctx, claims.ID)
	if err != nil {
		return Principal{}, "", err
	}
	if principal.UserID != claims.UserID || principal.WorkspaceID != claims.WorkspaceID {
		return Principal{}, "", fmt.Errorf("session claims do not match")
	}
	return principal, claims.ID, nil
}

func (s *Service) Logout(ctx context.Context, jti string) error {
	return s.store.RevokeWebSession(ctx, jti)
}

func (s *Service) createSession(ctx context.Context, principal Principal) (Session, error) {
	jti, err := RandomHex(32)
	if err != nil {
		return Session{}, err
	}

	now := s.now().UTC()
	token, expiresAt, err := s.tokens.Issue(principal.UserID, principal.WorkspaceID, jti, now)
	if err != nil {
		return Session{}, err
	}
	if err := s.store.CreateWebSession(ctx, principal, jti, expiresAt); err != nil {
		return Session{}, err
	}

	principal.ExpiresAt = expiresAt
	return Session{Principal: principal, Token: token, JTI: jti, ExpiresAt: expiresAt}, nil
}

func validateSignup(request SignupRequest) error {
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(request.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(request.Password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}
	return nil
}
