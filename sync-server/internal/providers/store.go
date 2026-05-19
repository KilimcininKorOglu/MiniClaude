package providers

import (
	"context"
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

func NewStore(pool *pgxpool.Pool, secrets *providercrypto.SecretBox) *Store {
	return &Store{pool: pool, secrets: secrets}
}

func (s *Store) UpsertSecret(ctx context.Context, providerID, secret string) error {
	providerID = strings.TrimSpace(providerID)
	if providerID == "" || secret == "" {
		return fmt.Errorf("provider id and secret are required")
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
