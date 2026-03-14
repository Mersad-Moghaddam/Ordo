package auth

import (
	"context"
	"time"

	domainauth "github.com/ordo/backend/internal/domain/auth"
)

type UserRepository interface {
	FindByEmailAddress(requestContext context.Context, emailAddress string) (domainauth.User, error)
	FindByUserID(requestContext context.Context, userID string) (domainauth.User, error)
	CreateUser(requestContext context.Context, user domainauth.User) error
}

type RefreshSessionRepository interface {
	CreateSession(requestContext context.Context, refreshSession domainauth.RefreshSession) error
	FindBySessionID(requestContext context.Context, sessionID string) (domainauth.RefreshSession, error)
	RevokeSession(requestContext context.Context, sessionID string, revokedAtTime time.Time, replacementSessionID *string) error
}
