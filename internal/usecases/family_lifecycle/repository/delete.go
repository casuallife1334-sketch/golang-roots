package repository

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) DeleteRelationship(ctx context.Context, treeID, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	familyIDs, err := affectedFamilies(ctx, tx, `
		SELECT fr.family_id FROM family_relationships fr
		JOIN relationships rel ON rel.id = fr.relationship_id
		WHERE rel.tree_id = $1 AND rel.id = $2
	`, treeID, id)
	if err != nil {
		return err
	}
	if err := r.relationships.DeleteRelationshipTx(ctx, tx, treeID, id); err != nil {
		return err
	}
	if err := r.families.ReconcileTx(ctx, tx, treeID, familyIDs); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) DeletePerson(ctx context.Context, treeID, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	familyIDs, err := affectedFamilies(ctx, tx, `
		SELECT DISTINCT fr.family_id FROM family_relationships fr
		JOIN relationships rel ON rel.id = fr.relationship_id
		WHERE rel.tree_id = $1 AND (rel.person1_id = $2 OR rel.person2_id = $2)
		UNION
		SELECT fm.family_id FROM family_members fm
		JOIN families f ON f.id = fm.family_id
		WHERE f.tree_id = $1 AND fm.person_id = $2
	`, treeID, id)
	if err != nil {
		return err
	}
	if err := r.persons.DeletePersonTx(ctx, tx, treeID, id); err != nil {
		return err
	}
	if err := r.families.ReconcileTx(ctx, tx, treeID, familyIDs); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func affectedFamilies(ctx context.Context, tx pgx.Tx, query, treeID, id string) ([]string, error) {
	rows, err := tx.Query(ctx, query, treeID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	slices.Sort(ids) // Lock affected families in a stable order across transactions.
	return ids, nil
}
