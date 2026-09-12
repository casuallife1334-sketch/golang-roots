package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func (r *RelationshipsRepository) CreateRelationship(ctx context.Context, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	var rel domain.Relationship
	err := r.db.QueryRow(ctx, `INSERT INTO relationships(id,person1_id,person2_id,type,direction) VALUES($1,$2,$3,$4,$5) RETURNING id,person1_id,person2_id,type,direction,created_at`, ulid.Make().String(), in.Person1ID, in.Person2ID, in.Type, in.Direction).Scan(&rel.ID, &rel.Person1ID, &rel.Person2ID, &rel.Type, &rel.Direction, &rel.CreatedAt)
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
