package domain

import "time"

type FamilyMemberRole string

const (
	FamilyMemberPartner FamilyMemberRole = "partner"
	FamilyMemberParent  FamilyMemberRole = "parent"
	FamilyMemberChild   FamilyMemberRole = "child"
)

type FamilyMember struct {
	PersonID string           `json:"person_id"`
	Role     FamilyMemberRole `json:"role"`
}

type Family struct {
	ID              string         `json:"id"`
	TreeID          string         `json:"tree_id"`
	Name            string         `json:"name"`
	Metadata        map[string]any `json:"metadata"`
	Members         []FamilyMember `json:"members"`
	RelationshipIDs []string       `json:"relationship_ids"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type CreateFamilyInput struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`
}

type PatchFamilyInput struct {
	Name     *string         `json:"name"`
	Metadata *map[string]any `json:"metadata"`
}

type AddFamilyMemberInput struct {
	PersonID string           `json:"person_id"`
	Role     FamilyMemberRole `json:"role"`
}
