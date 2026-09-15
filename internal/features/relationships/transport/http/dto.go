package http

import (
	"genealogy-tree/internal/core/domain"
	"time"
)

type CreateRelationshipRequest struct {
	Person1ID string                        `json:"person1_id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Person2ID string                        `json:"person2_id" example:"01JQ2Q4K8Z8F6M2Z3N4P5R6S7U"`
	Type      domain.RelationshipType       `json:"type" example:"parent_child"`
	Direction *domain.RelationshipDirection `json:"direction,omitempty" example:"parent"`
}

type RelationshipResponse struct {
	ID        string                        `json:"id" example:"01JQ2Q4K9Y8F6M2Z3N4P5R6S7V"`
	Person1ID string                        `json:"person1_id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Person2ID string                        `json:"person2_id" example:"01JQ2Q4K8Z8F6M2Z3N4P5R6S7U"`
	Type      domain.RelationshipType       `json:"type" example:"parent_child"`
	Direction *domain.RelationshipDirection `json:"direction,omitempty" example:"parent"`
	CreatedAt time.Time                     `json:"created_at" example:"2026-02-26T10:30:00Z"`
}

type RelationshipsResponse []RelationshipResponse
