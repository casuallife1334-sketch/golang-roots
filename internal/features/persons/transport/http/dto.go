package http

import (
	"genealogy-tree/internal/core/domain"
	"time"
)

type MetadataExample struct {
	City       string `json:"city" example:"Moscow"`
	Occupation string `json:"occupation" example:"Engineer"`
}

type CreatePersonRequest struct {
	FirstName string          `json:"first_name" example:"Ivan"`
	LastName  string          `json:"last_name" example:"Petrov"`
	BirthDate *time.Time      `json:"birth_date,omitempty" example:"1980-01-02T00:00:00Z"`
	DeathDate *time.Time      `json:"death_date,omitempty" example:"2024-05-10T00:00:00Z"`
	Gender    *domain.Gender  `json:"gender,omitempty" example:"male"`
	PhotoURL  *string         `json:"photo_url,omitempty" example:"persons/01JQ2Q4K7Y8F6M2Z3N4P5R6S7T/photo"`
	Metadata  MetadataExample `json:"metadata"`
}

type PersonResponse struct {
	ID        string          `json:"id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	FirstName string          `json:"first_name" example:"Ivan"`
	LastName  string          `json:"last_name" example:"Petrov"`
	BirthDate *time.Time      `json:"birth_date,omitempty" example:"1980-01-02T00:00:00Z"`
	DeathDate *time.Time      `json:"death_date,omitempty" example:"2024-05-10T00:00:00Z"`
	Gender    *domain.Gender  `json:"gender,omitempty" example:"male"`
	PhotoURL  *string         `json:"photo_url,omitempty" example:"persons/01JQ2Q4K7Y8F6M2Z3N4P5R6S7T/photo"`
	Metadata  MetadataExample `json:"metadata"`
	CreatedAt time.Time       `json:"created_at" example:"2026-02-26T10:30:00Z"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" example:"2026-02-26T10:30:00Z"`
}

type PersonsResponse []PersonResponse
