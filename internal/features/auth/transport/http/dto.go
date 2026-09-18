package http

import "time"

type RegisterRequest struct {
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"strong-password"`
}

type LoginRequest RegisterRequest

type AuthUserResponse struct {
	ID        string    `json:"id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Email     string    `json:"email" example:"ivan@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-02-26T10:30:00Z"`
}

type RegisterResponse AuthUserResponse

type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int64  `json:"expires_in" example:"900"`
}
