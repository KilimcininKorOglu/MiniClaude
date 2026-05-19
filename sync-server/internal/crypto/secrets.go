package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type SecretBox struct {
	gcm cipher.AEAD
}

func NewSecretBox(masterKey string) (*SecretBox, error) {
	if len(masterKey) < 32 {
		return nil, fmt.Errorf("provider secret master key must be at least 32 bytes")
	}

	key := sha256.Sum256([]byte(masterKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm cipher: %w", err)
	}
	return &SecretBox{gcm: gcm}, nil
}

func (b *SecretBox) Encrypt(plaintext, aad []byte) ([]byte, []byte, error) {
	nonce := make([]byte, b.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("generate secret nonce: %w", err)
	}
	return nonce, b.gcm.Seal(nil, nonce, plaintext, aad), nil
}

func (b *SecretBox) Decrypt(nonce, ciphertext, aad []byte) ([]byte, error) {
	plaintext, err := b.gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("decrypt secret: %w", err)
	}
	return plaintext, nil
}
