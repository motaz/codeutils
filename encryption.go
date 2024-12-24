package codeutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"fmt"

	"strings"
)

// deriveKey generates a 256-bit key from a password using SHA-256.
func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// encryptChunk encrypts a single chunk of text using AES-GCM.
func encryptChunk(password, plaintext string) (string, error) {

	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	encrypted := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decryptChunk decrypts a single Base64-encoded chunk using AES-GCM.
func decryptChunk(password, encodedCiphertext string) (string, error) {

	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64 ciphertext: %w", err)
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:aesGCM.NonceSize()], ciphertext[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// encryptText encrypts the entire text in chunks and concatenates the results.
func EncryptText(password, plaintext string) (string, error) {

	parts := strings.SplitAfter(plaintext, "\n")
	var encryptedParts []string

	for _, part := range parts {
		encryptedChunk, err := encryptChunk(password, part)
		if err != nil {
			return "", err
		}
		encryptedParts = append(encryptedParts, encryptedChunk)
	}

	return strings.Join(encryptedParts, "|"), nil
}

// DecryptText decrypts concatenated chunks and reassembles the plaintext.
func DecryptText(password, encodedCiphertext string) (string, error) {

	parts := strings.Split(encodedCiphertext, "|")
	var decryptedParts []string

	for _, part := range parts {
		decryptedChunk, err := decryptChunk(password, part)
		if err != nil {
			return "", err
		}
		decryptedParts = append(decryptedParts, decryptedChunk)
	}

	return strings.Join(decryptedParts, ""), nil
}
