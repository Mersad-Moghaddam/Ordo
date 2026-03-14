package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenClaims struct {
	SubjectUserID      string `json:"sub"`
	SessionID          string `json:"sid,omitempty"`
	AssignedRole       string `json:"role"`
	TokenType          string `json:"typ"`
	ExpiresAtUnix      int64  `json:"exp"`
	IssuedAtUnix       int64  `json:"iat"`
	RefreshTokenNumber int64  `json:"ver,omitempty"`
}

type TokenService interface {
	IssueToken(claims TokenClaims) (string, error)
	VerifyToken(tokenValue string) (TokenClaims, error)
	GenerateRandomToken() (string, error)
	HashToken(tokenValue string) string
}

type HMACTokenService struct {
	tokenSecret []byte
}

type TokenServiceOption func(tokenService *HMACTokenService) error

func NewHMACTokenService(options ...TokenServiceOption) (*HMACTokenService, error) {
	tokenService := &HMACTokenService{tokenSecret: []byte("phase1-default-secret")}
	for _, option := range options {
		if option == nil {
			continue
		}
		if optionError := option(tokenService); optionError != nil {
			return nil, optionError
		}
	}
	return tokenService, nil
}

func WithTokenSecret(secretValue string) TokenServiceOption {
	return func(tokenService *HMACTokenService) error {
		if secretValue == "" {
			return fmt.Errorf("token secret must not be empty")
		}
		tokenService.tokenSecret = []byte(secretValue)
		return nil
	}
}

func (hmacTokenService *HMACTokenService) IssueToken(claims TokenClaims) (string, error) {
	payloadBytes, marshalError := json.Marshal(claims)
	if marshalError != nil {
		return "", fmt.Errorf("marshal token claims failure: %w", marshalError)
	}
	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signatureSegment := hmacTokenService.sign(payloadSegment)
	return payloadSegment + "." + signatureSegment, nil
}

func (hmacTokenService *HMACTokenService) VerifyToken(tokenValue string) (TokenClaims, error) {
	tokenSegments := strings.Split(tokenValue, ".")
	if len(tokenSegments) != 2 {
		return TokenClaims{}, fmt.Errorf("malformed token")
	}
	expectedSignature := hmacTokenService.sign(tokenSegments[0])
	if !hmac.Equal([]byte(expectedSignature), []byte(tokenSegments[1])) {
		return TokenClaims{}, fmt.Errorf("signature mismatch")
	}
	payloadBytes, decodeError := base64.RawURLEncoding.DecodeString(tokenSegments[0])
	if decodeError != nil {
		return TokenClaims{}, fmt.Errorf("decode payload failure: %w", decodeError)
	}
	var tokenClaims TokenClaims
	if unmarshalError := json.Unmarshal(payloadBytes, &tokenClaims); unmarshalError != nil {
		return TokenClaims{}, fmt.Errorf("unmarshal claims failure: %w", unmarshalError)
	}
	if time.Now().Unix() >= tokenClaims.ExpiresAtUnix {
		return TokenClaims{}, fmt.Errorf("token expired")
	}
	return tokenClaims, nil
}

func (hmacTokenService *HMACTokenService) GenerateRandomToken() (string, error) {
	randomBytes := make([]byte, 32)
	if _, randomReadError := rand.Read(randomBytes); randomReadError != nil {
		return "", fmt.Errorf("random token generation failure: %w", randomReadError)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func (hmacTokenService *HMACTokenService) HashToken(tokenValue string) string {
	hashValue := sha256.Sum256([]byte(tokenValue))
	return base64.RawURLEncoding.EncodeToString(hashValue[:])
}

func (hmacTokenService *HMACTokenService) sign(payloadSegment string) string {
	hmacSigner := hmac.New(sha256.New, hmacTokenService.tokenSecret)
	_, _ = hmacSigner.Write([]byte(payloadSegment))
	signatureBytes := hmacSigner.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signatureBytes)
}
