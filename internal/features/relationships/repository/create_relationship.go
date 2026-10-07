package repository

import (
	"context"
	"encoding/json"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func (r *RelationshipsRepository) CreateRelationship(ctx context.Context, treeID string, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Relationship{}, err
	}
	defer tx.Rollback(ctx)
	relationship, err := r.CreateRelationshipTx(ctx, tx, treeID, in)
	if err != nil {
		return domain.Relationship{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Relationship{}, err
	}
	return relationship, nil
}

func (r *RelationshipsRepository) CreateRelationshipTx(ctx context.Context, tx pgx.Tx, treeID string, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	metadata, err := json.Marshal(in.Metadata)
	if err != nil {
		return domain.Relationship{}, err
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO relationships(id,tree_id,person1_id,person2_id,type,direction,metadata)
		SELECT $1,$2,p1.id,p2.id,$5,$6,$7
		FROM persons p1 JOIN persons p2 ON p2.id = $4
		WHERE p1.id = $3 AND p1.tree_id = $2 AND p2.tree_id = $2
		RETURNING id,person1_id,person2_id,type,direction,metadata,created_at,updated_at
	`, ulid.Make().String(), treeID, in.Person1ID, in.Person2ID, in.Type, in.Direction, metadata)
	rel, err := r.scanRelationship(row)
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
