package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *DocumentsRepository) GetDocument(ctx context.Context, treeID, id string) (domain.Document, error) {
	document, err := scanDocument(r.db.QueryRow(ctx, `
		SELECT id,tree_id,person_id,relationship_id,file_name,storage_key,content_type,size_bytes,created_by,created_at
		FROM documents WHERE tree_id = $1 AND id = $2
	`, treeID, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Document{}, ErrNotFound
	}
	return document, err
}

type documentRow interface {
	Scan(...any) error
}

func scanDocument(row documentRow) (domain.Document, error) {
	var document domain.Document
	var personID, relationshipID *string
	if err := row.Scan(&document.ID, &document.TreeID, &personID, &relationshipID, &document.FileName, &document.StorageKey, &document.ContentType, &document.SizeBytes, &document.CreatedBy, &document.CreatedAt); err != nil {
		return domain.Document{}, err
	}
	switch {
	case personID != nil:
		document.Owner = domain.DocumentOwner{Type: domain.DocumentOwnerPerson, ID: *personID}
	case relationshipID != nil:
		document.Owner = domain.DocumentOwner{Type: domain.DocumentOwnerRelationship, ID: *relationshipID}
	}
	return document, nil
}
