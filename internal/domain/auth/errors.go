package auth

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyRegistered = errors.New("email already registered")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session expired")
var ErrInvalidToken = errors.New("invalid token")
var ErrForbiddenRole = errors.New("forbidden role")
