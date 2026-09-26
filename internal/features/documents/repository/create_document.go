package repository

import (
	"context"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

var ErrOwnerNotFound = fmt.Errorf("%w: document owner not found", coreerrors.ErrNotFound)

func (r *DocumentsRepository) CreateDocument(ctx context.Context, input domain.CreateDocumentInput) (domain.Document, error) {
	var document domain.Document
	var personID, relationshipID *string
	err := r.db.QueryRow(ctx, `
		INSERT INTO documents (id,tree_id,person_id,relationship_id,file_name,storage_key,content_type,size_bytes,created_by)
		SELECT $1,$2,
			CASE WHEN $3 = 'person' THEN $4 ELSE NULL END,
			CASE WHEN $3 = 'relationship' THEN $4 ELSE NULL END,
			$5,$6,$7,$8,$9
		WHERE ($3 = 'person' AND EXISTS (SELECT 1 FROM persons WHERE id = $4 AND tree_id = $2))
		   OR ($3 = 'relationship' AND EXISTS (SELECT 1 FROM relationships WHERE id = $4 AND tree_id = $2))
		RETURNING id,tree_id,person_id,relationship_id,file_name,storage_key,content_type,size_bytes,created_by,created_at
	`, ulid.Make().String(), input.TreeID, input.Owner.Type, input.Owner.ID, input.FileName, input.StorageKey, input.ContentType, input.SizeBytes, input.CreatedBy).Scan(
		&document.ID,
		&document.TreeID,
		&personID,
		&relationshipID,
		&document.FileName,
		&document.StorageKey,
		&document.ContentType,
		&document.SizeBytes,
		&document.CreatedBy,
		&document.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Document{}, ErrOwnerNotFound
		}
		return domain.Document{}, err
	}
	if personID != nil {
		document.Owner = domain.DocumentOwner{Type: domain.DocumentOwnerPerson, ID: *personID}
	} else if relationshipID != nil {
		document.Owner = domain.DocumentOwner{Type: domain.DocumentOwnerRelationship, ID: *relationshipID}
	}
	return document, nil
}
