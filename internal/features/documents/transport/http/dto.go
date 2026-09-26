package http

import (
	"genealogy-tree/internal/core/domain"
	"time"
)

type DocumentResponse struct {
	ID          string               `json:"id" example:"01JQ2Q4K9Y8F6M2Z3N4P5R6S7V"`
	TreeID      string               `json:"tree_id" example:"01JQ2Q4K7Y8F6M2Z3N4P5R6S7T"`
	Owner       domain.DocumentOwner `json:"owner"`
	FileName    string               `json:"file_name" example:"Marriage certificate.pdf"`
	ContentType string               `json:"content_type" example:"application/pdf"`
	SizeBytes   int64                `json:"size_bytes" example:"245760"`
	CreatedBy   string               `json:"created_by" example:"01JQ2Q4K8Z8F6M2Z3N4P5R6S7U"`
	CreatedAt   time.Time            `json:"created_at" example:"2026-02-26T10:30:00Z"`
}

type DocumentsResponse []DocumentResponse
