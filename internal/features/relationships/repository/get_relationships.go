package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"genealogy-tree/internal/core/domain"
)

func (r *RelationshipsRepository) GetRelationships(ctx context.Context, treeID, personID string) ([]domain.Relationship, error) {
	query := `SELECT id,person1_id,person2_id,type,direction,metadata,created_at,updated_at FROM relationships WHERE tree_id=$1`
	args := []any{treeID}
	if personID != "" {
		query += ` AND (person1_id=$2 OR person2_id=$2)`
		args = append(args, personID)
	}
	query += ` ORDER BY id`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.Relationship{}
	for rows.Next() {
		var relationship domain.Relationship
		var metadata []byte
		if err := rows.Scan(
			&relationship.ID,
			&relationship.Person1ID,
			&relationship.Person2ID,
			&relationship.Type,
			&relationship.Direction,
			&metadata,
			&relationship.CreatedAt,
			&relationship.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadata, &relationship.Metadata); err != nil {
			return nil, fmt.Errorf("decode relationship metadata: %w", err)
		}
		items = append(items, relationship)
	}
	return items, rows.Err()
}
