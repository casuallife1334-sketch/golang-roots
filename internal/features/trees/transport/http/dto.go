package http

import (
	"genealogy-tree/internal/core/domain"
	"time"
)

type CreateTreeRequest struct {
	Name string `json:"name" example:"Family Petrov"`
}

type PatchTreeRequest struct {
	Name *string `json:"name,omitempty" example:"Family Petrov Updated"`
}

type TreeResponse struct {
	ID        string          `json:"id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	OwnerID   string          `json:"owner_id" example:"01JQ2Q4K8Z8F6M2Z3N4P5R6S7U"`
	Role      domain.TreeRole `json:"role" example:"owner"`
	Name      string          `json:"name" example:"Family Petrov"`
	CreatedAt time.Time       `json:"created_at" example:"2026-02-26T10:30:00Z"`
	UpdatedAt time.Time       `json:"updated_at" example:"2026-02-26T10:30:00Z"`
}

type TreesResponse []TreeResponse
