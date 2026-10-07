package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// ReconcileTx rebuilds the derived membership of auto-created families from
// their remaining relationships. Manually edited families own their members.
func (r *FamiliesRepository) ReconcileTx(ctx context.Context, tx pgx.Tx, treeID string, familyIDs []string) error {
	for _, familyID := range familyIDs {
		var autoCreated bool
		err := tx.QueryRow(ctx, `
			SELECT auto_created FROM families WHERE tree_id = $1 AND id = $2 FOR UPDATE
		`, treeID, familyID).Scan(&autoCreated)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return err
		}
		if !autoCreated {
			continue
		}
		result, err := tx.Exec(ctx, `
			DELETE FROM families f
			WHERE f.tree_id = $1 AND f.id = $2
			  AND NOT EXISTS (SELECT 1 FROM family_relationships fr WHERE fr.family_id = f.id)
		`, treeID, familyID)
		if err != nil {
			return err
		}
		if result.RowsAffected() != 0 {
			continue
		}
		if _, err := tx.Exec(ctx, `DELETE FROM family_members WHERE family_id = $1`, familyID); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO family_members (tree_id, family_id, person_id, role)
			SELECT DISTINCT $2::char(26), $1::char(26), member.person_id, member.role::family_member_role
			FROM family_relationships fr
			JOIN relationships rel ON rel.id = fr.relationship_id AND rel.tree_id = $2
			CROSS JOIN LATERAL (
				VALUES
				  (rel.person1_id, CASE WHEN rel.type = 'spouse' THEN 'partner' ELSE 'parent' END),
				  (rel.person2_id, CASE WHEN rel.type = 'spouse' THEN 'partner' ELSE 'child' END)
			) AS member(person_id, role)
			WHERE fr.family_id = $1
		`, familyID, treeID)
		if err != nil {
			return err
		}
	}
	return nil
}
