package auth

type RegisterRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	AssignedRole string `json:"role"`
}

type LoginRequest struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

type RefreshRequest struct {
	RefreshTokenValue string `json:"refreshToken"`
}

type TokenResponse struct {
	AccessTokenValue  string `json:"accessToken"`
	RefreshTokenValue string `json:"refreshToken"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
