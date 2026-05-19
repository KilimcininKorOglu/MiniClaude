package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	providercrypto "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool    *pgxpool.Pool
	secrets *providercrypto.SecretBox
}

type SecretMetadata struct {
	ProviderID string    `json:"provider_id"`
	KeyVersion int       `json:"key_version"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Provider struct {
	ID            string
	Name          string
	BaseURL       string
	Model         string
	DefaultSonnet string
	DefaultHaiku  string
	DefaultOpus   string
	Description   string
	HasSecret     bool
	UpdatedAt     time.Time
}

type UpsertInput struct {
	WorkspaceID   string
	Name          string
	BaseURL       string
	Model         string
	DefaultSonnet string
	DefaultHaiku  string
	DefaultOpus   string
	Description   string
	Secret        string
}

func NewStore(pool *pgxpool.Pool, secrets *providercrypto.SecretBox) *Store {
	return &Store{pool: pool, secrets: secrets}
}

func (s *Store) List(ctx context.Context, workspaceID string) ([]Provider, error) {
	rows, err := s.pool.Query(ctx, `
		select pc.id::text, pc.name, coalesce(pc.base_url, ''), coalesce(pc.model, ''),
			coalesce(pc.default_models->>'sonnet', ''), coalesce(pc.default_models->>'haiku', ''), coalesce(pc.default_models->>'opus', ''),
			coalesce(pc.metadata->>'description', ''), ps.provider_id is not null, pc.updated_at
		from provider_configs pc
		left join provider_secrets ps on ps.provider_id = pc.id
		where pc.workspace_id = $1
		order by pc.name asc
	`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list providers: %w", err)
	}
	defer rows.Close()

	providers := []Provider{}
	for rows.Next() {
		provider := Provider{}
		if err := rows.Scan(&provider.ID, &provider.Name, &provider.BaseURL, &provider.Model, &provider.DefaultSonnet, &provider.DefaultHaiku, &provider.DefaultOpus, &provider.Description, &provider.HasSecret, &provider.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan provider: %w", err)
		}
		providers = append(providers, provider)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate providers: %w", err)
	}
	return providers, nil
}

func (s *Store) Upsert(ctx context.Context, input UpsertInput) (string, error) {
	input.WorkspaceID = strings.TrimSpace(input.WorkspaceID)
	input.Name = strings.TrimSpace(input.Name)
	input.BaseURL = strings.TrimSpace(input.BaseURL)
	input.Model = strings.TrimSpace(input.Model)
	input.Secret = strings.TrimSpace(input.Secret)
	if input.WorkspaceID == "" || input.Name == "" {
		return "", fmt.Errorf("workspace and provider name are required")
	}

	defaultModels, err := json.Marshal(map[string]string{
		"sonnet": strings.TrimSpace(input.DefaultSonnet),
		"haiku":  strings.TrimSpace(input.DefaultHaiku),
		"opus":   strings.TrimSpace(input.DefaultOpus),
	})
	if err != nil {
		return "", err
	}
	metadata, err := json.Marshal(map[string]string{"description": strings.TrimSpace(input.Description)})
	if err != nil {
		return "", err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin provider transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var providerID string
	if err := tx.QueryRow(ctx, `
		insert into provider_configs (workspace_id, name, base_url, model, default_models, metadata)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (workspace_id, name) do update set
			base_url = excluded.base_url,
			model = excluded.model,
			default_models = excluded.default_models,
			metadata = excluded.metadata,
			updated_at = now()
		returning id::text
	`, input.WorkspaceID, input.Name, input.BaseURL, input.Model, defaultModels, metadata).Scan(&providerID); err != nil {
		return "", fmt.Errorf("upsert provider config: %w", err)
	}

	if input.Secret != "" {
		if s.secrets == nil {
			return "", fmt.Errorf("provider secret storage is not configured")
		}
		nonce, ciphertext, err := s.secrets.Encrypt([]byte(input.Secret), []byte(providerID))
		if err != nil {
			return "", err
		}
		if _, err := tx.Exec(ctx, `
			insert into provider_secrets (provider_id, key_version, nonce, ciphertext)
			values ($1, 1, $2, $3)
			on conflict (provider_id) do update set
				key_version = provider_secrets.key_version + 1,
				nonce = excluded.nonce,
				ciphertext = excluded.ciphertext,
				updated_at = now()
		`, providerID, nonce, ciphertext); err != nil {
			return "", fmt.Errorf("store provider secret: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit provider transaction: %w", err)
	}
	return providerID, nil
}

func (s *Store) Delete(ctx context.Context, workspaceID, providerID string) error {
	workspaceID = strings.TrimSpace(workspaceID)
	providerID = strings.TrimSpace(providerID)
	if workspaceID == "" || providerID == "" {
		return fmt.Errorf("workspace and provider id are required")
	}
	_, err := s.pool.Exec(ctx, `delete from provider_configs where workspace_id = $1 and id = $2`, workspaceID, providerID)
	if err != nil {
		return fmt.Errorf("delete provider: %w", err)
	}
	return nil
}

func (s *Store) SettingsProviders(ctx context.Context, workspaceID string) (map[string]any, error) {
	providers, err := s.List(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	settings := map[string]any{}
	for _, provider := range providers {
		env := map[string]string{}
		if provider.BaseURL != "" {
			env["ANTHROPIC_BASE_URL"] = provider.BaseURL
		}
		if secret, _, err := s.ReadSecret(ctx, provider.ID); err == nil && secret != "" {
			env["ANTHROPIC_AUTH_TOKEN"] = secret
		}
		defaultModels := map[string]string{}
		if provider.DefaultSonnet != "" {
			defaultModels["sonnet"] = provider.DefaultSonnet
		}
		if provider.DefaultHaiku != "" {
			defaultModels["haiku"] = provider.DefaultHaiku
		}
		if provider.DefaultOpus != "" {
			defaultModels["opus"] = provider.DefaultOpus
		}
		value := map[string]any{"env": env}
		if provider.Model != "" {
			value["model"] = provider.Model
		}
		if len(defaultModels) > 0 {
			value["defaultModels"] = defaultModels
		}
		if provider.Description != "" {
			value["description"] = provider.Description
		}
		settings[provider.Name] = value
	}
	return settings, nil
}

func (s *Store) UpsertSecret(ctx context.Context, providerID, secret string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" || secret == "" {
		return fmt.Errorf("provider id and secret are required")
	}
	if s.secrets == nil {
		return fmt.Errorf("provider secret storage is not configured")
	}

	nonce, ciphertext, err := s.secrets.Encrypt([]byte(secret), []byte(providerID))
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx, `
		insert into provider_secrets (provider_id, key_version, nonce, ciphertext)
		values ($1, 1, $2, $3)
		on conflict (provider_id) do update set
			key_version = provider_secrets.key_version + 1,
			nonce = excluded.nonce,
			ciphertext = excluded.ciphertext,
			updated_at = now()
	`, providerID, nonce, ciphertext)
	if err != nil {
		return fmt.Errorf("store provider secret: %w", err)
	}
	return nil
}

func (s *Store) ReadSecret(ctx context.Context, providerID string) (string, SecretMetadata, error) {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return "", SecretMetadata{}, fmt.Errorf("provider id is required")
	}
	if s.secrets == nil {
		return "", SecretMetadata{}, fmt.Errorf("provider secret storage is not configured")
	}

	metadata := SecretMetadata{ProviderID: providerID}
	var nonce []byte
	var ciphertext []byte
	err := s.pool.QueryRow(ctx, `
		select key_version, nonce, ciphertext, updated_at
		from provider_secrets
		where provider_id = $1
	`, providerID).Scan(&metadata.KeyVersion, &nonce, &ciphertext, &metadata.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", SecretMetadata{}, fmt.Errorf("provider secret not found")
		}
		return "", SecretMetadata{}, fmt.Errorf("read provider secret: %w", err)
	}

	plaintext, err := s.secrets.Decrypt(nonce, ciphertext, []byte(providerID))
	if err != nil {
		return "", SecretMetadata{}, err
	}
	return string(plaintext), metadata, nil
}

func (s *Store) DeleteSecret(ctx context.Context, providerID string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" {
		return fmt.Errorf("provider id is required")
	}

	_, err := s.pool.Exec(ctx, `
		delete from provider_secrets
		where provider_id = $1
	`, providerID)
	if err != nil {
		return fmt.Errorf("delete provider secret: %w", err)
	}
	return nil
}
