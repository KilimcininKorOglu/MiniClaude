package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	passwordMemory      = 64 * 1024
	passwordIterations  = 3
	passwordParallelism = 2
	passwordSaltLength  = 16
	passwordKeyLength   = 32
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is required")
	}

	salt := make([]byte, passwordSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, passwordIterations, passwordMemory, passwordParallelism, passwordKeyLength)
	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		passwordMemory,
		passwordIterations,
		passwordParallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, fmt.Errorf("invalid password hash format")
	}

	params, err := parseParams(parts[3])
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("decode password salt: %w", err)
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("decode password hash: %w", err)
	}

	actual := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.parallelism, uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1, nil
}

type argonParams struct {
	memory      uint32
	time        uint32
	parallelism uint8
}

func parseParams(raw string) (argonParams, error) {
	values := map[string]string{}
	for _, part := range strings.Split(raw, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			return argonParams{}, fmt.Errorf("invalid argon2 parameter %q", part)
		}
		values[key] = value
	}

	memory, err := parseUint32(values["m"])
	if err != nil {
		return argonParams{}, fmt.Errorf("invalid argon2 memory: %w", err)
	}
	time, err := parseUint32(values["t"])
	if err != nil {
		return argonParams{}, fmt.Errorf("invalid argon2 time: %w", err)
	}
	parallelism, err := parseUint8(values["p"])
	if err != nil {
		return argonParams{}, fmt.Errorf("invalid argon2 parallelism: %w", err)
	}

	return argonParams{memory: memory, time: time, parallelism: parallelism}, nil
}

func parseUint32(raw string) (uint32, error) {
	value, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func parseUint8(raw string) (uint8, error) {
	value, err := strconv.ParseUint(raw, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(value), nil
}
