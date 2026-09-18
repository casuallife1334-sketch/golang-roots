package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func (r *RelationshipsRepository) CreateRelationship(ctx context.Context, treeID string, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	var rel domain.Relationship
	err := r.db.QueryRow(ctx, `
		INSERT INTO relationships(id,tree_id,person1_id,person2_id,type,direction)
		SELECT $1,$2,p1.id,p2.id,$5,$6
		FROM persons p1 JOIN persons p2 ON p2.id = $4
		WHERE p1.id = $3 AND p1.tree_id = $2 AND p2.tree_id = $2
		RETURNING id,person1_id,person2_id,type,direction,created_at
	`, ulid.Make().String(), treeID, in.Person1ID, in.Person2ID, in.Type, in.Direction).Scan(&rel.ID, &rel.Person1ID, &rel.Person2ID, &rel.Type, &rel.Direction, &rel.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return rel, ErrDuplicate
			}
			if pgErr.Code == "23503" {
				return rel, ErrPersonNotFound
			}
		}
		return rel, err
	}
	return rel, nil
}
