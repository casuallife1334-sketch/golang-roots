package repository

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (r *RelationshipsRepository) GetRelationships(ctx context.Context, treeID, personID string) ([]domain.Relationship, error) {
	query := `SELECT id,person1_id,person2_id,type,direction,created_at FROM relationships WHERE tree_id=$1`
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
		if err := rows.Scan(
			&relationship.ID,
			&relationship.Person1ID,
			&relationship.Person2ID,
			&relationship.Type,
			&relationship.Direction,
			&relationship.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, relationship)
	}
	return items, rows.Err()
}
