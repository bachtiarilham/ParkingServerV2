// package hash

// func Hash(password string) (string, error)

// func Compare(
// 	hash string,
// 	password string,
// ) error

// Package hash menyediakan fungsi hashing password.
// Implementasi ini menggunakan golang.org/x/crypto/bcrypt.
// Pastikan dependency sudah ada di go.mod saat run: go get golang.org/x/crypto
package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const saltLen = 16

// Hash menghasilkan hash password dengan random salt.
// Format output: "<salt>$<sha256(salt+password)>"
func Hash(password string) (string, error) {
	saltBytes := make([]byte, saltLen)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	salt := hex.EncodeToString(saltBytes)
	hash := sha256Hash(salt, password)
	return salt + "$" + hash, nil
}

// Compare memverifikasi password terhadap hash yang tersimpan.
func Compare(hash string, password string) error {
	parts := strings.SplitN(hash, "$", 2)
	if len(parts) != 2 {
		return errors.New("format hash tidak valid")
	}
	salt, storedHash := parts[0], parts[1]
	computed := sha256Hash(salt, password)
	if computed != storedHash {
		return errors.New("password salah")
	}
	return nil
}

func sha256Hash(salt, password string) string {
	h := sha256.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil))
}
