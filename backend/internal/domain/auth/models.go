package auth

import "time"

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type User struct {
	UserID        string
	EmailAddress  string
	PasswordHash  string
	AssignedRole  Role
	CreatedAtTime time.Time
	UpdatedAtTime time.Time
}

type RefreshSession struct {
	SessionID            string
	UserID               string
	RefreshTokenHash     string
	RefreshTokenVersion  int64
	IssuedAtTime         time.Time
	ExpiresAtTime        time.Time
	RevokedAtTime        *time.Time
	ReplacementSessionID *string
}

type TokenPair struct {
	AccessTokenValue  string
	RefreshTokenValue string
	AccessExpiryTime  time.Time
	RefreshExpiryTime time.Time
}
