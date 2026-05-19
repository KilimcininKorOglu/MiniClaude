package crypto

import "testing"

func TestSecretBoxEncryptDecrypt(t *testing.T) {
	box, err := NewSecretBox("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}

	nonce, ciphertext, err := box.Encrypt([]byte("secret"), []byte("provider-id"))
	if err != nil {
		t.Fatal(err)
	}
	if string(ciphertext) == "secret" {
		t.Fatal("ciphertext must not equal plaintext")
	}

	plaintext, err := box.Decrypt(nonce, ciphertext, []byte("provider-id"))
	if err != nil {
		t.Fatal(err)
	}
	if string(plaintext) != "secret" {
		t.Fatalf("expected decrypted secret, got %q", plaintext)
	}
}
