package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domainauth "github.com/ordo/backend/internal/domain/auth"
	infrastructuresecurity "github.com/ordo/backend/internal/infrastructure/security"
	"github.com/ordo/backend/internal/repository"
	repositoryauth "github.com/ordo/backend/internal/repository/auth"
)

type Service struct {
	userRepository           repositoryauth.UserRepository
	refreshSessionRepository repositoryauth.RefreshSessionRepository
	transactionManager       repository.TransactionManager
	passwordHasher           infrastructuresecurity.PasswordHasher
	tokenService             infrastructuresecurity.TokenService
	accessTokenTTL           time.Duration
	refreshTokenTTL          time.Duration
	nowFunction              func() time.Time
	identifierFunction       func() (string, error)
}

type Option func(service *Service)

type RegisterInput struct {
	EmailAddress string
	Password     string
	AssignedRole domainauth.Role
}

type LoginInput struct {
	EmailAddress string
	Password     string
}

type RefreshInput struct {
	RefreshTokenValue string
}

func NewService(
	userRepository repositoryauth.UserRepository,
	refreshSessionRepository repositoryauth.RefreshSessionRepository,
	transactionManager repository.TransactionManager,
	passwordHasher infrastructuresecurity.PasswordHasher,
	tokenService infrastructuresecurity.TokenService,
	options ...Option,
) *Service {
	service := &Service{
		userRepository:           userRepository,
		refreshSessionRepository: refreshSessionRepository,
		transactionManager:       transactionManager,
		passwordHasher:           passwordHasher,
		tokenService:             tokenService,
		accessTokenTTL:           15 * time.Minute,
		refreshTokenTTL:          24 * time.Hour,
		nowFunction:              time.Now,
	}
	service.identifierFunction = service.generateIdentifier
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}
	return service
}

func WithAccessTokenTTL(accessTokenTTL time.Duration) Option {
	return func(service *Service) {
		if accessTokenTTL > 0 {
			service.accessTokenTTL = accessTokenTTL
		}
	}
}

func WithRefreshTokenTTL(refreshTokenTTL time.Duration) Option {
	return func(service *Service) {
		if refreshTokenTTL > 0 {
			service.refreshTokenTTL = refreshTokenTTL
		}
	}
}

func WithNowFunction(nowFunction func() time.Time) Option {
	return func(service *Service) {
		if nowFunction != nil {
			service.nowFunction = nowFunction
		}
	}
}

func (service *Service) Register(requestContext context.Context, registerInput RegisterInput) (domainauth.TokenPair, error) {
	if registerInput.AssignedRole == "" {
		registerInput.AssignedRole = domainauth.RoleMember
	}
	var tokenPair domainauth.TokenPair
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		issuedTokenPair, issuedError := service.registerWithinTransaction(transactionContext, registerInput)
		if issuedError != nil {
			return issuedError
		}
		tokenPair = issuedTokenPair
		return nil
	})
	if transactionError != nil {
		return domainauth.TokenPair{}, transactionError
	}
	return tokenPair, nil
}

func (service *Service) Login(requestContext context.Context, loginInput LoginInput) (domainauth.TokenPair, error) {
	user, findError := service.userRepository.FindByEmailAddress(requestContext, loginInput.EmailAddress)
	if findError != nil {
		if errors.Is(findError, domainauth.ErrUserNotFound) {
			return domainauth.TokenPair{}, domainauth.ErrInvalidCredentials
		}
		return domainauth.TokenPair{}, fmt.Errorf("find user by email failure: %w", findError)
	}
	if !service.passwordHasher.MatchesHash(loginInput.Password, user.PasswordHash) {
		return domainauth.TokenPair{}, domainauth.ErrInvalidCredentials
	}
	return service.issueSessionAndTokens(requestContext, user.UserID, user.AssignedRole, 1, service.nowFunction())
}

func (service *Service) Refresh(requestContext context.Context, refreshInput RefreshInput) (domainauth.TokenPair, error) {
	refreshClaims, refreshOpaqueValue, verificationError := service.verifyRefreshToken(refreshInput.RefreshTokenValue)
	if verificationError != nil {
		return domainauth.TokenPair{}, verificationError
	}
	refreshSession, sessionError := service.refreshSessionRepository.FindBySessionID(requestContext, refreshClaims.SessionID)
	if sessionError != nil {
		return domainauth.TokenPair{}, sessionError
	}
	if invalidSessionError := service.validateSession(refreshSession, refreshOpaqueValue, refreshClaims.RefreshTokenNumber); invalidSessionError != nil {
		return domainauth.TokenPair{}, invalidSessionError
	}
	user, userError := service.userRepository.FindByUserID(requestContext, refreshClaims.SubjectUserID)
	if userError != nil {
		return domainauth.TokenPair{}, userError
	}
	return service.rotateSession(requestContext, refreshSession, user)
}

func (service *Service) AuthorizeRole(tokenValue string, allowedRoles ...domainauth.Role) error {
	accessClaims, verificationError := service.tokenService.VerifyToken(tokenValue)
	if verificationError != nil || accessClaims.TokenType != "access" {
		return domainauth.ErrInvalidToken
	}
	for _, allowedRole := range allowedRoles {
		if allowedRole == domainauth.Role(accessClaims.AssignedRole) {
			return nil
		}
	}
	return domainauth.ErrForbiddenRole
}

func (service *Service) generateIdentifier() (string, error) {
	randomValue, randomError := service.tokenService.GenerateRandomToken()
	if randomError != nil {
		return "", randomError
	}
	normalizedValue := strings.ReplaceAll(randomValue, "-", "")
	if len(normalizedValue) > 36 {
		return normalizedValue[:36], nil
	}
	return normalizedValue, nil
}
