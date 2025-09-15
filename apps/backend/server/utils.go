package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/alexedwards/argon2id"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

// Argon2 helpers using alexedwards/argon2id
func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func VerifyPassword(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

// AES-GCM helpers (keyBase64 must be base64 of 32 bytes)
func AESGCMEncryptB64(keyBase64 string, plaintext []byte) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("aes key must be 32 bytes (base64)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func AESGCMDecryptB64(keyBase64 string, b64cipher string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("aes key must be 32 bytes (base64)")
	}
	ct, err := base64.StdEncoding.DecodeString(b64cipher)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ct) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ct[:nonceSize], ct[nonceSize:]
	pt, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

// RenderMarkdown converts markdown to sanitized HTML (goldmark + bluemonday)
func RenderMarkdown(md string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	// sanitize produced HTML
	policy := bluemonday.UGCPolicy()
	sanitized := policy.SanitizeBytes(buf.Bytes())
	return string(sanitized), nil
}
