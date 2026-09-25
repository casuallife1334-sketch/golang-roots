package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *RelationshipsRepository) GetRelationship(ctx context.Context, treeID, id string) (domain.Relationship, error) {
	return r.scanRelationship(r.db.QueryRow(ctx, `SELECT id,person1_id,person2_id,type,direction,metadata,created_at,updated_at FROM relationships WHERE tree_id=$1 AND id=$2`, treeID, id))
}

func (r *RelationshipsRepository) scanRelationship(row pgx.Row) (domain.Relationship, error) {
	var rel domain.Relationship
	var metadata []byte
	err := row.Scan(&rel.ID, &rel.Person1ID, &rel.Person2ID, &rel.Type, &rel.Direction, &metadata, &rel.CreatedAt, &rel.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Relationship{}, ErrNotFound
	}
	if err != nil {
		return domain.Relationship{}, fmt.Errorf("scan relationship: %w", err)
	}
	if err := json.Unmarshal(metadata, &rel.Metadata); err != nil {
		return domain.Relationship{}, fmt.Errorf("decode relationship metadata: %w", err)
	}
	return rel, nil
}
