package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauth "github.com/ordo/backend/internal/domain/auth"
	infrastructuresecurity "github.com/ordo/backend/internal/infrastructure/security"
)

type mockTransactionManager struct{}

func (mockTransactionManager mockTransactionManager) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	return transactionWorkload(requestContext)
}

type mockUserRepository struct {
	findByEmailFunction func(emailAddress string) (domainauth.User, error)
	findByIDFunction    func(userID string) (domainauth.User, error)
	createUserFunction  func(user domainauth.User) error
}

func (mockRepository mockUserRepository) FindByEmailAddress(requestContext context.Context, emailAddress string) (domainauth.User, error) {
	return mockRepository.findByEmailFunction(emailAddress)
}

func (mockRepository mockUserRepository) FindByUserID(requestContext context.Context, userID string) (domainauth.User, error) {
	return mockRepository.findByIDFunction(userID)
}

func (mockRepository mockUserRepository) CreateUser(requestContext context.Context, user domainauth.User) error {
	return mockRepository.createUserFunction(user)
}

type mockRefreshSessionRepository struct {
	createSessionFunction func(refreshSession domainauth.RefreshSession) error
	findSessionFunction   func(sessionID string) (domainauth.RefreshSession, error)
	revokeSessionFunction func(sessionID string, replacementSessionID *string) error
}

func (mockRepository mockRefreshSessionRepository) CreateSession(requestContext context.Context, refreshSession domainauth.RefreshSession) error {
	return mockRepository.createSessionFunction(refreshSession)
}

func (mockRepository mockRefreshSessionRepository) FindBySessionID(requestContext context.Context, sessionID string) (domainauth.RefreshSession, error) {
	return mockRepository.findSessionFunction(sessionID)
}

func (mockRepository mockRefreshSessionRepository) RevokeSession(requestContext context.Context, sessionID string, revokedAtTime time.Time, replacementSessionID *string) error {
	return mockRepository.revokeSessionFunction(sessionID, replacementSessionID)
}

func TestRegisterAndLogin(testingSuite *testing.T) {
	nowValue := time.Unix(1700000000, 0)
	tokenService, serviceError := infrastructuresecurity.NewHMACTokenService(infrastructuresecurity.WithTokenSecret("phase1-secret"))
	if serviceError != nil {
		testingSuite.Fatalf("token service creation failure: %v", serviceError)
	}

	persistedUsers := map[string]domainauth.User{}
	persistedUsersByID := map[string]domainauth.User{}
	persistedSessions := map[string]domainauth.RefreshSession{}
	authService := NewService(
		mockUserRepository{
			findByEmailFunction: func(emailAddress string) (domainauth.User, error) {
				user, hasUser := persistedUsers[emailAddress]
				if !hasUser {
					return domainauth.User{}, domainauth.ErrUserNotFound
				}
				return user, nil
			},
			findByIDFunction: func(userID string) (domainauth.User, error) {
				user, hasUser := persistedUsersByID[userID]
				if !hasUser {
					return domainauth.User{}, domainauth.ErrUserNotFound
				}
				return user, nil
			},
			createUserFunction: func(user domainauth.User) error {
				persistedUsers[user.EmailAddress] = user
				persistedUsersByID[user.UserID] = user
				return nil
			},
		},
		mockRefreshSessionRepository{
			createSessionFunction: func(refreshSession domainauth.RefreshSession) error {
				persistedSessions[refreshSession.SessionID] = refreshSession
				return nil
			},
			findSessionFunction: func(sessionID string) (domainauth.RefreshSession, error) {
				refreshSession, hasSession := persistedSessions[sessionID]
				if !hasSession {
					return domainauth.RefreshSession{}, domainauth.ErrSessionNotFound
				}
				return refreshSession, nil
			},
			revokeSessionFunction: func(sessionID string, replacementSessionID *string) error {
				refreshSession, hasSession := persistedSessions[sessionID]
				if !hasSession {
					return domainauth.ErrSessionNotFound
				}
				revokedAtTime := nowValue
				refreshSession.RevokedAtTime = &revokedAtTime
				refreshSession.ReplacementSessionID = replacementSessionID
				persistedSessions[sessionID] = refreshSession
				return nil
			},
		},
		mockTransactionManager{},
		infrastructuresecurity.NewSHA256PasswordHasher(),
		tokenService,
		WithNowFunction(func() time.Time { return nowValue }),
	)
	authService.identifierFunction = func() (string, error) { return "identifier-value", nil }

	registerResult, registerError := authService.Register(context.Background(), RegisterInput{EmailAddress: "owner@ordo.dev", Password: "test-password", AssignedRole: domainauth.RoleOwner})
	if registerError != nil {
		testingSuite.Fatalf("register failure: %v", registerError)
	}
	if registerResult.AccessTokenValue == "" || registerResult.RefreshTokenValue == "" {
		testingSuite.Fatalf("tokens should not be empty")
	}

	_, duplicateError := authService.Register(context.Background(), RegisterInput{EmailAddress: "owner@ordo.dev", Password: "test-password", AssignedRole: domainauth.RoleOwner})
	if !errors.Is(duplicateError, domainauth.ErrEmailAlreadyRegistered) {
		testingSuite.Fatalf("expected duplicate registration error, got: %v", duplicateError)
	}

	loginResult, loginError := authService.Login(context.Background(), LoginInput{EmailAddress: "owner@ordo.dev", Password: "test-password"})
	if loginError != nil {
		testingSuite.Fatalf("login failure: %v", loginError)
	}
	if loginResult.AccessTokenValue == "" || loginResult.RefreshTokenValue == "" {
		testingSuite.Fatalf("login tokens should not be empty")
	}
}

func TestAuthorizeRole(testingSuite *testing.T) {
	tokenService, serviceError := infrastructuresecurity.NewHMACTokenService(infrastructuresecurity.WithTokenSecret("phase1-secret"))
	if serviceError != nil {
		testingSuite.Fatalf("token service creation failure: %v", serviceError)
	}
	authService := NewService(
		mockUserRepository{},
		mockRefreshSessionRepository{},
		mockTransactionManager{},
		infrastructuresecurity.NewSHA256PasswordHasher(),
		tokenService,
	)

	accessTokenValue, tokenError := tokenService.IssueToken(infrastructuresecurity.TokenClaims{SubjectUserID: "identifier", AssignedRole: string(domainauth.RoleAdmin), TokenType: "access", IssuedAtUnix: time.Now().Unix(), ExpiresAtUnix: time.Now().Add(5 * time.Minute).Unix()})
	if tokenError != nil {
		testingSuite.Fatalf("token issue failure: %v", tokenError)
	}
	if authorizationError := authService.AuthorizeRole(accessTokenValue, domainauth.RoleOwner, domainauth.RoleAdmin); authorizationError != nil {
		testingSuite.Fatalf("expected authorization success, got: %v", authorizationError)
	}
	if authorizationError := authService.AuthorizeRole(accessTokenValue, domainauth.RoleOwner); !errors.Is(authorizationError, domainauth.ErrForbiddenRole) {
		testingSuite.Fatalf("expected forbidden role error, got: %v", authorizationError)
	}
}
