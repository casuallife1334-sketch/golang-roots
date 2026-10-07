package domain

import (
	"genealogy-tree/internal/core/nullable"
	"time"
)

type RelationshipType string
type RelationshipDirection string

const (
	RelationshipParentChild RelationshipType      = "parent_child"
	RelationshipSpouse      RelationshipType      = "spouse"
	DirectionParent         RelationshipDirection = "parent"
	DirectionChild          RelationshipDirection = "child"
)

type Relationship struct {
	ID string `json:"id"`
	// For parent_child, direction identifies which of Person1 and Person2 is the parent.
	// For spouse, both fields contain the two spouse IDs in normalized order.
	Person1ID string                 `json:"person1_id"`
	Person2ID string                 `json:"person2_id"`
	Type      RelationshipType       `json:"type"`
	Direction *RelationshipDirection `json:"direction,omitempty"`
	Metadata  map[string]any         `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at"`
}

type CreateRelationshipInput struct {
	Person1ID string                 `json:"person1_id"`
	Person2ID string                 `json:"person2_id"`
	Type      RelationshipType       `json:"type"`
	Direction *RelationshipDirection `json:"direction"`
	Metadata  map[string]any         `json:"metadata"`
}

type PatchRelationshipInput struct {
	Metadata nullable.Value[map[string]any] `json:"metadata" swaggertype:"object" extensions:"x-nullable"`
}
