package domain

import "time"

type TreeRole string

const (
	TreeRoleOwner  TreeRole = "owner"
	TreeRoleEditor TreeRole = "editor"
	TreeRoleViewer TreeRole = "viewer"
)

type Tree struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Role      TreeRole  `json:"role"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTreeInput struct {
	Name string `json:"name"`
}

type PatchTreeInput struct {
	Name *string `json:"name"`
}
