// crypto/encryptor_test.go
package crypto_test

import (
	"encoding/base64"
	"strings"
	"testing"
	"wn/internal/domain/services/crypto"
)

func TestNewEncryptor(t *testing.T) {
	t.Run("creates encryptor with valid key", func(t *testing.T) {
		masterKey := "my-secret-master-key"
		encryptor := crypto.NewEncryptor(masterKey)

		if encryptor == nil {
			t.Fatal("expected encryptor to be created, got nil")
		}
	})

	t.Run("same master key produces same derived key", func(t *testing.T) {
		e1 := crypto.NewEncryptor("same-key")
		e2 := crypto.NewEncryptor("same-key")

		// Проверяем, что шифрование одним ключом даёт одинаковый результат при расшифровке
		encrypted, err := e1.Encrypt("test")
		if err != nil {
			t.Fatalf("Encrypt() error = %v", err)
		}
		decrypted, err := e2.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt() error = %v", err)
		}
		if decrypted != "test" {
			t.Errorf("expected 'test', got %q", decrypted)
		}
	})

	t.Run("different master keys produce different results", func(t *testing.T) {
		e1 := crypto.NewEncryptor("key-one")
		e2 := crypto.NewEncryptor("key-two")

		encrypted, _ := e1.Encrypt("secret")
		_, err := e2.Decrypt(encrypted)
		if err == nil {
			t.Fatal("expected error when decrypting with wrong key")
		}
	})
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	testCases := []struct {
		name      string
		plainText string
		masterKey string
	}{
		{"simple text", "Hello, World!", "test-key"},
		{"unicode", "Привет, мир! 🌍", "unicode-key"},
		{"long text", strings.Repeat("A", 10000), "long-key"},
		{"special chars", "!@#$%^&*()_+-=[]{}", "special-key"},
		{"whitespace", "Line1\nLine2\tTab", "ws-key"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptor := crypto.NewEncryptor(tc.masterKey)

			encrypted, err := encryptor.Encrypt(tc.plainText)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}
			if encrypted == "" && tc.plainText != "" {
				t.Fatal("expected non-empty encrypted output")
			}

			decrypted, err := encryptor.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if decrypted != tc.plainText {
				t.Errorf("roundtrip failed:\nwant: %q\ngot:  %q", tc.plainText, decrypted)
			}
		})
	}
}

func TestEncrypt_EmptyString(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	result, err := encryptor.Encrypt("")
	if err != nil {
		t.Errorf("Encrypt(\"\") unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("Encrypt(\"\") = %q, want \"\"", result)
	}
}

func TestDecrypt_EmptyString(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	result, err := encryptor.Decrypt("")
	if err != nil {
		t.Errorf("Decrypt(\"\") unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("Decrypt(\"\") = %q, want \"\"", result)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	invalidInputs := []string{
		"not-valid-base64!!!",
		"!!!!",
		"@@@@",
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, err := encryptor.Decrypt(input)
			if err == nil {
				t.Error("expected error for invalid base64 input")
			}
			if !strings.Contains(err.Error(), "failed to decode base64") {
				t.Errorf("expected base64 decode error, got: %v", err)
			}
		})
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	original := "secret message"
	encrypted, err := encryptor.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Портим зашифрованный текст
	tampered := encrypted[:10] + "X" + encrypted[11:]

	_, err = encryptor.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected error when decrypting tampered ciphertext")
	}
	if !strings.Contains(err.Error(), "failed to decrypt") {
		t.Errorf("expected decryption error, got: %v", err)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	original := "confidential data"

	e1 := crypto.NewEncryptor("key-one")
	encrypted, err := e1.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	e2 := crypto.NewEncryptor("key-two")
	_, err = e2.Decrypt(encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
	if !strings.Contains(err.Error(), "failed to decrypt") {
		t.Errorf("expected authentication error, got: %v", err)
	}
}

func TestEncrypt_NonceRandomness(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")
	plainText := "same message"

	// Шифруем одно сообщение несколько раз
	encryptions := make([]string, 5)
	for i := range encryptions {
		enc, err := encryptor.Encrypt(plainText)
		if err != nil {
			t.Fatalf("Encrypt() error = %v", err)
		}
		encryptions[i] = enc
	}

	// Все результаты должны быть разными (из-за случайного nonce)
	for i := 0; i < len(encryptions); i++ {
		for j := i + 1; j < len(encryptions); j++ {
			if encryptions[i] == encryptions[j] {
				t.Errorf("encryption %d and %d produced same output", i, j)
			}
		}
	}

	// Но все должны корректно расшифровываться
	for i, enc := range encryptions {
		dec, err := encryptor.Decrypt(enc)
		if err != nil {
			t.Errorf("Decrypt() #%d error = %v", i, err)
			continue
		}
		if dec != plainText {
			t.Errorf("Decrypt() #%d = %q, want %q", i, dec, plainText)
		}
	}
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	// Валидный base64, но слишком короткий для AES-GCM
	tooShort := base64.StdEncoding.EncodeToString([]byte{1, 2, 3})

	_, err := encryptor.Decrypt(tooShort)
	if err == nil {
		t.Fatal("expected error for ciphertext too short")
	}
	if !strings.Contains(err.Error(), "ciphertext too short") {
		t.Errorf("expected 'ciphertext too short' error, got: %v", err)
	}
}

func TestEncrypt_OutputIsBase64(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-key")

	plainText := "test data"
	encrypted, err := encryptor.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Проверяем, что результат — валидный base64
	_, err = base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Errorf("encrypted output is not valid base64: %v", err)
	}

	// Проверяем отсутствие недопустимых символов
	invalidChars := []string{" ", "\n", "\r", "\t"}
	for _, ch := range invalidChars {
		if strings.Contains(encrypted, ch) {
			t.Errorf("encrypted output contains invalid character: %q", ch)
		}
	}
}
