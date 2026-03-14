package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domainauth "github.com/ordo/backend/internal/domain/auth"
	infrastructuresecurity "github.com/ordo/backend/internal/infrastructure/security"
)

func (service *Service) registerWithinTransaction(requestContext context.Context, registerInput RegisterInput) (domainauth.TokenPair, error) {
	findError := service.validateEmailAvailability(requestContext, registerInput.EmailAddress)
	if findError != nil {
		return domainauth.TokenPair{}, findError
	}
	userID, creationTime, creationError := service.createUser(requestContext, registerInput)
	if creationError != nil {
		return domainauth.TokenPair{}, creationError
	}
	return service.issueSessionAndTokens(requestContext, userID, registerInput.AssignedRole, 1, creationTime)
}

func (service *Service) validateEmailAvailability(requestContext context.Context, emailAddress string) error {
	_, findError := service.userRepository.FindByEmailAddress(requestContext, emailAddress)
	if findError == nil {
		return domainauth.ErrEmailAlreadyRegistered
	}
	if errors.Is(findError, domainauth.ErrUserNotFound) {
		return nil
	}
	return fmt.Errorf("find user precheck failure: %w", findError)
}

func (service *Service) createUser(requestContext context.Context, registerInput RegisterInput) (string, time.Time, error) {
	userIdentifier, identifierError := service.identifierFunction()
	if identifierError != nil {
		return "", time.Time{}, fmt.Errorf("create user identifier failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	user := domainauth.User{
		UserID:        userIdentifier,
		EmailAddress:  registerInput.EmailAddress,
		PasswordHash:  service.passwordHasher.HashPassword(registerInput.Password),
		AssignedRole:  registerInput.AssignedRole,
		CreatedAtTime: creationTime,
		UpdatedAtTime: creationTime,
	}
	if creationError := service.userRepository.CreateUser(requestContext, user); creationError != nil {
		return "", time.Time{}, fmt.Errorf("create user failure: %w", creationError)
	}
	return userIdentifier, creationTime, nil
}

func (service *Service) issueSessionAndTokens(
	requestContext context.Context,
	userID string,
	assignedRole domainauth.Role,
	refreshTokenNumber int64,
	issuedAtTime time.Time,
) (domainauth.TokenPair, error) {
	refreshSessionID, sessionError := service.identifierFunction()
	if sessionError != nil {
		return domainauth.TokenPair{}, fmt.Errorf("refresh session identifier failure: %w", sessionError)
	}
	accessTokenValue, accessExpiryTime, accessError := service.issueAccessToken(userID, assignedRole, issuedAtTime)
	if accessError != nil {
		return domainauth.TokenPair{}, accessError
	}
	refreshEnvelope, refreshExpiryTime, refreshError := service.issueRefreshEnvelope(userID, refreshSessionID, assignedRole, refreshTokenNumber, issuedAtTime)
	if refreshError != nil {
		return domainauth.TokenPair{}, refreshError
	}
	persistError := service.persistRefreshSession(requestContext, refreshSessionID, userID, refreshEnvelope, refreshTokenNumber, issuedAtTime, refreshExpiryTime)
	if persistError != nil {
		return domainauth.TokenPair{}, persistError
	}
	return domainauth.TokenPair{AccessTokenValue: accessTokenValue, RefreshTokenValue: refreshEnvelope, AccessExpiryTime: accessExpiryTime, RefreshExpiryTime: refreshExpiryTime}, nil
}

func (service *Service) issueAccessToken(userID string, assignedRole domainauth.Role, issuedAtTime time.Time) (string, time.Time, error) {
	accessExpiryTime := issuedAtTime.Add(service.accessTokenTTL)
	accessClaims := infrastructuresecurity.TokenClaims{
		SubjectUserID: userID,
		AssignedRole:  string(assignedRole),
		TokenType:     "access",
		IssuedAtUnix:  issuedAtTime.Unix(),
		ExpiresAtUnix: accessExpiryTime.Unix(),
	}
	accessTokenValue, accessError := service.tokenService.IssueToken(accessClaims)
	if accessError != nil {
		return "", time.Time{}, fmt.Errorf("issue access token failure: %w", accessError)
	}
	return accessTokenValue, accessExpiryTime, nil
}

func (service *Service) issueRefreshEnvelope(
	userID string,
	refreshSessionID string,
	assignedRole domainauth.Role,
	refreshTokenNumber int64,
	issuedAtTime time.Time,
) (string, time.Time, error) {
	refreshOpaqueValue, refreshTokenError := service.tokenService.GenerateRandomToken()
	if refreshTokenError != nil {
		return "", time.Time{}, fmt.Errorf("generate refresh token failure: %w", refreshTokenError)
	}
	refreshExpiryTime := issuedAtTime.Add(service.refreshTokenTTL)
	refreshClaims := infrastructuresecurity.TokenClaims{SubjectUserID: userID, SessionID: refreshSessionID, AssignedRole: string(assignedRole), TokenType: "refresh", IssuedAtUnix: issuedAtTime.Unix(), ExpiresAtUnix: refreshExpiryTime.Unix(), RefreshTokenNumber: refreshTokenNumber}
	refreshSignedValue, issueError := service.tokenService.IssueToken(refreshClaims)
	if issueError != nil {
		return "", time.Time{}, fmt.Errorf("issue refresh token failure: %w", issueError)
	}
	return refreshOpaqueValue + "." + refreshSignedValue, refreshExpiryTime, nil
}

func (service *Service) persistRefreshSession(
	requestContext context.Context,
	refreshSessionID string,
	userID string,
	refreshEnvelope string,
	refreshTokenNumber int64,
	issuedAtTime time.Time,
	expiresAtTime time.Time,
) error {
	refreshOpaqueValue, parseError := parseRefreshEnvelope(refreshEnvelope)
	if parseError != nil {
		return parseError
	}
	refreshSession := domainauth.RefreshSession{SessionID: refreshSessionID, UserID: userID, RefreshTokenHash: service.tokenService.HashToken(refreshOpaqueValue), RefreshTokenVersion: refreshTokenNumber, IssuedAtTime: issuedAtTime, ExpiresAtTime: expiresAtTime}
	if persistError := service.refreshSessionRepository.CreateSession(requestContext, refreshSession); persistError != nil {
		return fmt.Errorf("create refresh session failure: %w", persistError)
	}
	return nil
}

func (service *Service) verifyRefreshToken(refreshEnvelope string) (infrastructuresecurity.TokenClaims, string, error) {
	refreshOpaqueValue, refreshSignedValue, parseError := parseRefreshEnvelopeSigned(refreshEnvelope)
	if parseError != nil {
		return infrastructuresecurity.TokenClaims{}, "", domainauth.ErrInvalidToken
	}
	refreshClaims, verificationError := service.tokenService.VerifyToken(refreshSignedValue)
	if verificationError != nil || refreshClaims.TokenType != "refresh" {
		return infrastructuresecurity.TokenClaims{}, "", domainauth.ErrInvalidToken
	}
	return refreshClaims, refreshOpaqueValue, nil
}

func parseRefreshEnvelope(refreshEnvelope string) (string, error) {
	refreshOpaqueValue, _, parseError := parseRefreshEnvelopeSigned(refreshEnvelope)
	if parseError != nil {
		return "", parseError
	}
	return refreshOpaqueValue, nil
}

func parseRefreshEnvelopeSigned(refreshEnvelope string) (string, string, error) {
	tokenSegments := strings.SplitN(refreshEnvelope, ".", 2)
	if len(tokenSegments) != 2 || tokenSegments[0] == "" || tokenSegments[1] == "" {
		return "", "", fmt.Errorf("invalid refresh token envelope")
	}
	return tokenSegments[0], tokenSegments[1], nil
}

func (service *Service) validateSession(refreshSession domainauth.RefreshSession, refreshOpaqueValue string, refreshTokenNumber int64) error {
	if service.tokenService.HashToken(refreshOpaqueValue) != refreshSession.RefreshTokenHash {
		return domainauth.ErrInvalidToken
	}
	if refreshSession.RefreshTokenVersion != refreshTokenNumber {
		return domainauth.ErrInvalidToken
	}
	if refreshSession.RevokedAtTime != nil || service.nowFunction().After(refreshSession.ExpiresAtTime) {
		return domainauth.ErrSessionExpired
	}
	return nil
}

func (service *Service) rotateSession(requestContext context.Context, currentSession domainauth.RefreshSession, user domainauth.User) (domainauth.TokenPair, error) {
	newSessionID, sessionError := service.identifierFunction()
	if sessionError != nil {
		return domainauth.TokenPair{}, fmt.Errorf("new session identifier failure: %w", sessionError)
	}
	revokedTime := service.nowFunction()
	if revokeError := service.refreshSessionRepository.RevokeSession(requestContext, currentSession.SessionID, revokedTime, &newSessionID); revokeError != nil {
		return domainauth.TokenPair{}, fmt.Errorf("revoke existing session failure: %w", revokeError)
	}
	return service.issueSessionAndTokens(requestContext, user.UserID, user.AssignedRole, currentSession.RefreshTokenVersion+1, revokedTime)
}
