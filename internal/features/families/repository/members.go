package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *FamiliesRepository) AddMember(ctx context.Context, treeID, familyID string, input domain.AddFamilyMemberInput) error {
	err := r.db.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO family_members (tree_id, family_id, person_id, role)
			SELECT f.tree_id, f.id, p.id, $3
			FROM families f
			JOIN persons p ON p.tree_id = f.tree_id
			WHERE f.id = $1 AND f.tree_id = $2 AND p.id = $4
			RETURNING family_id
		)
		UPDATE families f
		SET auto_created = false, updated_at = now()
		FROM inserted
		WHERE f.id = inserted.family_id
		RETURNING f.id
	`, familyID, treeID, input.Role, input.PersonID).Scan(new(string))
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

func (r *FamiliesRepository) RemoveMember(ctx context.Context, treeID, familyID, personID string) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM family_members fm
		USING families f
		WHERE fm.family_id = f.id
		  AND f.tree_id = $1
		  AND f.id = $2
		  AND fm.person_id = $3
	`, treeID, familyID, personID)
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
