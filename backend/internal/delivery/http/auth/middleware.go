package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	domainauth "github.com/ordo/backend/internal/domain/auth"
	usecaseauth "github.com/ordo/backend/internal/usecase/auth"
)

type Middleware struct {
	authService *usecaseauth.Service
}

func NewMiddleware(authService *usecaseauth.Service) *Middleware {
	return &Middleware{authService: authService}
}

func (middleware *Middleware) RequireRoles(allowedRoles ...domainauth.Role) fiber.Handler {
	return func(fiberContext *fiber.Ctx) error {
		authorizationHeader := fiberContext.Get("Authorization")
		tokenValue := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if tokenValue == "" || tokenValue == authorizationHeader {
			return fiberContext.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{ErrorMessage: "missing bearer token"})
		}
		authorizationError := middleware.authService.AuthorizeRole(tokenValue, allowedRoles...)
		if authorizationError != nil {
			return respondWithAuthError(fiberContext, authorizationError)
		}
		return fiberContext.Next()
	}
}
