package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

type PasswordHasher interface {
	HashPassword(passwordValue string) string
	MatchesHash(passwordValue string, passwordHash string) bool
}

type SHA256PasswordHasher struct{}

func NewSHA256PasswordHasher() *SHA256PasswordHasher {
	return &SHA256PasswordHasher{}
}

func (sha256PasswordHasher *SHA256PasswordHasher) HashPassword(passwordValue string) string {
	hashBytes := sha256.Sum256([]byte(passwordValue))
	return base64.RawURLEncoding.EncodeToString(hashBytes[:])
}

func (sha256PasswordHasher *SHA256PasswordHasher) MatchesHash(passwordValue string, passwordHash string) bool {
	candidateHash := sha256PasswordHasher.HashPassword(passwordValue)
	if len(candidateHash) != len(passwordHash) {
		return false
	}
	comparisonResult := subtle.ConstantTimeCompare([]byte(candidateHash), []byte(passwordHash))
	return comparisonResult == 1
}
