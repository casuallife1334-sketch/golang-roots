package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *RelationshipsRepository) GetRelationship(ctx context.Context, id string) (domain.Relationship, error) {
	var rel domain.Relationship
	err := r.db.QueryRow(ctx, `SELECT id,person1_id,person2_id,type,direction,created_at FROM relationships WHERE id=$1`, id).Scan(&rel.ID, &rel.Person1ID, &rel.Person2ID, &rel.Type, &rel.Direction, &rel.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return rel, ErrNotFound
	}
	return rel, err
}
