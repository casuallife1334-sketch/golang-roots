package repository

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (r *DocumentsRepository) GetDocuments(ctx context.Context, treeID string, owner domain.DocumentOwner) ([]domain.Document, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id,tree_id,person_id,relationship_id,file_name,storage_key,content_type,size_bytes,created_by,created_at
		FROM documents
		WHERE tree_id = $1
		  AND (($2 = 'person' AND person_id = $3) OR ($2 = 'relationship' AND relationship_id = $3))
		ORDER BY created_at, id
	`, treeID, owner.Type, owner.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.Document{}
	for rows.Next() {
		item, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
