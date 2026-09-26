package domain

import "time"

type DocumentOwnerType string

const (
	DocumentOwnerPerson       DocumentOwnerType = "person"
	DocumentOwnerRelationship DocumentOwnerType = "relationship"
)

type DocumentOwner struct {
	Type DocumentOwnerType `json:"type"`
	ID   string            `json:"id"`
}

type Document struct {
	ID          string        `json:"id"`
	TreeID      string        `json:"tree_id"`
	Owner       DocumentOwner `json:"owner"`
	FileName    string        `json:"file_name"`
	ContentType string        `json:"content_type"`
	SizeBytes   int64         `json:"size_bytes"`
	CreatedBy   string        `json:"created_by"`
	StorageKey  string        `json:"-"`
	CreatedAt   time.Time     `json:"created_at"`
}

type CreateDocumentInput struct {
	TreeID      string
	Owner       DocumentOwner
	FileName    string
	StorageKey  string
	ContentType string
	SizeBytes   int64
	CreatedBy   string
}
