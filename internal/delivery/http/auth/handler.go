package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	domainauth "github.com/ordo/backend/internal/domain/auth"
	usecaseauth "github.com/ordo/backend/internal/usecase/auth"
)

type Handler struct {
	authService *usecaseauth.Service
}

func NewHandler(authService *usecaseauth.Service) *Handler {
	return &Handler{authService: authService}
}

func (handler *Handler) RegisterRoutes(fiberApplication fiber.Router) {
	fiberApplication.Post("/register", handler.register)
	fiberApplication.Post("/login", handler.login)
	fiberApplication.Post("/refresh", handler.refresh)
}

func (handler *Handler) register(fiberContext *fiber.Ctx) error {
	var registerRequest RegisterRequest
	if parserError := fiberContext.BodyParser(&registerRequest); parserError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	tokenPair, registerError := handler.authService.Register(fiberContext.UserContext(), usecaseauth.RegisterInput{EmailAddress: registerRequest.EmailAddress, Password: registerRequest.Password, AssignedRole: domainauth.Role(registerRequest.AssignedRole)})
	if registerError != nil {
		return respondWithAuthError(fiberContext, registerError)
	}
	return fiberContext.Status(fiber.StatusCreated).JSON(TokenResponse{AccessTokenValue: tokenPair.AccessTokenValue, RefreshTokenValue: tokenPair.RefreshTokenValue})
}

func (handler *Handler) login(fiberContext *fiber.Ctx) error {
	var loginRequest LoginRequest
	if parserError := fiberContext.BodyParser(&loginRequest); parserError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	tokenPair, loginError := handler.authService.Login(fiberContext.UserContext(), usecaseauth.LoginInput{EmailAddress: loginRequest.EmailAddress, Password: loginRequest.Password})
	if loginError != nil {
		return respondWithAuthError(fiberContext, loginError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(TokenResponse{AccessTokenValue: tokenPair.AccessTokenValue, RefreshTokenValue: tokenPair.RefreshTokenValue})
}

func (handler *Handler) refresh(fiberContext *fiber.Ctx) error {
	var refreshRequest RefreshRequest
	if parserError := fiberContext.BodyParser(&refreshRequest); parserError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	tokenPair, refreshError := handler.authService.Refresh(fiberContext.UserContext(), usecaseauth.RefreshInput{RefreshTokenValue: refreshRequest.RefreshTokenValue})
	if refreshError != nil {
		return respondWithAuthError(fiberContext, refreshError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(TokenResponse{AccessTokenValue: tokenPair.AccessTokenValue, RefreshTokenValue: tokenPair.RefreshTokenValue})
}

func respondWithAuthError(fiberContext *fiber.Ctx, responseError error) error {
	if errors.Is(responseError, domainauth.ErrInvalidCredentials) || errors.Is(responseError, domainauth.ErrInvalidToken) {
		return fiberContext.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domainauth.ErrEmailAlreadyRegistered) {
		return fiberContext.Status(fiber.StatusConflict).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domainauth.ErrForbiddenRole) {
		return fiberContext.Status(fiber.StatusForbidden).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	return fiberContext.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "internal error"})
}
