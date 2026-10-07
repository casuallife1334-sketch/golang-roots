package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *FamiliesRepository) AttachRelationship(ctx context.Context, treeID, familyID, relationshipID string) error {
	err := r.db.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO family_relationships (tree_id, family_id, relationship_id)
			SELECT f.tree_id, f.id, r.id
			FROM families f
			JOIN relationships r ON r.tree_id = f.tree_id
			WHERE f.id = $1 AND f.tree_id = $2 AND r.id = $3
			RETURNING family_id
		)
		UPDATE families f
		SET auto_created = false, updated_at = now()
		FROM inserted
		WHERE f.id = inserted.family_id
		RETURNING f.id
	`, familyID, treeID, relationshipID).Scan(new(string))
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicate
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *FamiliesRepository) DetachRelationship(ctx context.Context, treeID, familyID, relationshipID string) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM family_relationships fr
		USING families f
		WHERE fr.family_id = f.id
		  AND f.tree_id = $1
		  AND f.id = $2
		  AND fr.relationship_id = $3
	`, treeID, familyID, relationshipID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	_, err = r.db.Exec(ctx, `
		UPDATE families
		SET auto_created = false, updated_at = now()
		WHERE tree_id = $1 AND id = $2
	`, treeID, familyID)
	return err
}
