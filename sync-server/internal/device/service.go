package device

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
)

const userCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Service struct {
	store  *Store
	tokens *ClientTokenManager
	appURL string
	now    func() time.Time
}

type StartRequest struct {
	ClientName string `json:"client_name"`
}

type StartResponse struct {
	DeviceCode      string    `json:"device_code"`
	UserCode        string    `json:"user_code"`
	VerificationURL string    `json:"verification_url"`
	ExpiresAt       time.Time `json:"expires_at"`
	IntervalSeconds int       `json:"interval_seconds"`
}

type ApproveRequest struct {
	UserCode   string `json:"user_code"`
	ClientName string `json:"client_name"`
}

type PollRequest struct {
	DeviceCode string `json:"device_code"`
}

type PollResponse struct {
	Status      string    `json:"status"`
	ClientID    string    `json:"client_id,omitempty"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

func NewService(store *Store, tokens *ClientTokenManager, appURL string) *Service {
	return &Service{store: store, tokens: tokens, appURL: strings.TrimRight(appURL, "/"), now: time.Now}
}

func (s *Service) Start(ctx context.Context, request StartRequest) (StartResponse, error) {
	deviceCode, err := auth.RandomHex(32)
	if err != nil {
		return StartResponse{}, err
	}
	userCode, err := randomUserCode()
	if err != nil {
		return StartResponse{}, err
	}
	expiresAt := s.now().UTC().Add(10 * time.Minute)
	if err := s.store.CreateRequest(ctx, Request{DeviceCode: deviceCode, UserCode: userCode, ExpiresAt: expiresAt}); err != nil {
		return StartResponse{}, err
	}
	return StartResponse{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationURL: fmt.Sprintf("%s/device", s.appURL),
		ExpiresAt:       expiresAt,
		IntervalSeconds: 2,
	}, nil
}

func (s *Service) Approve(ctx context.Context, principal auth.Principal, request ApproveRequest) (string, error) {
	if strings.TrimSpace(request.UserCode) == "" {
		return "", fmt.Errorf("user code is required")
	}
	return s.store.Approve(ctx, request.UserCode, principal, request.ClientName)
}

func (s *Service) Poll(ctx context.Context, request PollRequest) (PollResponse, error) {
	result, err := s.store.Poll(ctx, request.DeviceCode)
	if err != nil {
		return PollResponse{}, err
	}
	response := PollResponse{Status: result.Status}
	if result.Status != "approved" {
		return response, nil
	}
	accessToken, expiresAt, err := s.tokens.Issue(result.ClientID, result.WorkspaceID, s.now().UTC())
	if err != nil {
		return PollResponse{}, err
	}
	response.ClientID = result.ClientID
	response.WorkspaceID = result.WorkspaceID
	response.AccessToken = accessToken
	response.ExpiresAt = expiresAt
	return response, nil
}

func (s *Service) AuthenticateClient(ctx context.Context, tokenValue string) (ClientPrincipal, error) {
	claims, err := s.tokens.Parse(tokenValue)
	if err != nil {
		return ClientPrincipal{}, err
	}
	return s.store.ClientPrincipal(ctx, claims.ClientID, claims.WorkspaceID)
}

func randomUserCode() (string, error) {
	var builder strings.Builder
	for i := 0; i < 8; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(userCodeAlphabet))))
		if err != nil {
			return "", fmt.Errorf("generate user code: %w", err)
		}
		if i == 4 {
			builder.WriteByte('-')
		}
		builder.WriteByte(userCodeAlphabet[index.Int64()])
	}
	return builder.String(), nil
}
