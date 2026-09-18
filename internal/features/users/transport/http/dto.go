package http

import "time"

type UserResponse struct {
	ID        string    `json:"id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Email     string    `json:"email" example:"ivan@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-26T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-02-26T10:30:00Z"`
}
