package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/nullable"
	"time"
)

type CreateFamilyRequest struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata" swaggertype:"object"`
}

type PatchFamilyRequest struct {
	Name     *string                        `json:"name"`
	Metadata nullable.Value[map[string]any] `json:"metadata" swaggertype:"object"`
}

type AddFamilyMemberRequest struct {
	PersonID string                  `json:"person_id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Role     domain.FamilyMemberRole `json:"role" example:"parent"`
}

type AttachRelationshipRequest struct {
	RelationshipID string `json:"relationship_id" example:"01JQ2Q4K9Y8F6M2Z3N4P5R6S7V"`
}

type FamilyResponse struct {
	ID              string                `json:"id" example:"01JQ2Q4K9Y8F6M2Z3N4P5R6S7V"`
	TreeID          string                `json:"tree_id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Name            string                `json:"name"`
	Metadata        map[string]any        `json:"metadata" swaggertype:"object"`
	Members         []domain.FamilyMember `json:"members"`
	RelationshipIDs []string              `json:"relationship_ids"`
	CreatedAt       time.Time             `json:"created_at" example:"2026-02-26T10:30:00Z"`
	UpdatedAt       time.Time             `json:"updated_at" example:"2026-02-26T10:35:00Z"`
}

type FamiliesResponse []FamilyResponse
