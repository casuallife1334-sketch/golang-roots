package repository

import (
	"context"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

var ErrFamilyNotFound = fmt.Errorf("%w: family not found", coreerrors.ErrNotFound)

func (r *FamiliesRepository) EnsureRelationshipTx(ctx context.Context, tx pgx.Tx, treeID string, relationship domain.Relationship, requestedID *string) error {
	familyID, err := r.findOrCreateFamily(ctx, tx, treeID, relationship, requestedID)
	if err != nil {
		return err
	}
	if relationship.Type == domain.RelationshipSpouse {
		if err := addFamilyMember(ctx, tx, treeID, familyID, relationship.Person1ID, domain.FamilyMemberPartner); err != nil {
			return err
		}
		if err := addFamilyMember(ctx, tx, treeID, familyID, relationship.Person2ID, domain.FamilyMemberPartner); err != nil {
			return err
		}
	} else {
		if err := addFamilyMember(ctx, tx, treeID, familyID, relationship.Person1ID, domain.FamilyMemberParent); err != nil {
			return err
		}
		if err := addFamilyMember(ctx, tx, treeID, familyID, relationship.Person2ID, domain.FamilyMemberChild); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO family_relationships (tree_id, family_id, relationship_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, treeID, familyID, relationship.ID); err != nil {
		return err
	}
	return nil
}

func (r *FamiliesRepository) findOrCreateFamily(ctx context.Context, tx pgx.Tx, treeID string, relationship domain.Relationship, requestedID *string) (string, error) {
	if requestedID != nil {
		var familyID string
		if err := tx.QueryRow(ctx, `SELECT id FROM families WHERE tree_id = $1 AND id = $2`, treeID, *requestedID).Scan(&familyID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", ErrFamilyNotFound
			}
			return "", err
		}
		return familyID, nil
	}

	var familyID string
	var err error
	if relationship.Type == domain.RelationshipSpouse {
		// Prefer a family in which both people are already adults. Being
		// children in the same family must not make them spouses in that family.
		familyID, err = uniqueFamily(ctx, tx, `
			SELECT f.id
			FROM families f
			WHERE f.tree_id = $1
			  AND EXISTS (SELECT 1 FROM family_members m WHERE m.family_id = f.id AND m.person_id = $2 AND m.role IN ('partner', 'parent'))
			  AND EXISTS (SELECT 1 FROM family_members m WHERE m.family_id = f.id AND m.person_id = $3 AND m.role IN ('partner', 'parent'))
			LIMIT 2
		`, treeID, relationship.Person1ID, relationship.Person2ID)
		if err != nil || familyID != "" {
			return familyID, err
		}
		// A family started by a single parent and their children can become
		// the spouses' family, but not a family with another adult partner.
		familyID, err = uniqueFamily(ctx, tx, `
			SELECT f.id
			FROM families f
			WHERE f.tree_id = $1
			  AND EXISTS (SELECT 1 FROM family_members m WHERE m.family_id = f.id AND m.person_id IN ($2, $3) AND m.role IN ('partner', 'parent'))
			  AND NOT EXISTS (SELECT 1 FROM family_members m WHERE m.family_id = f.id AND m.role IN ('partner', 'parent') AND m.person_id NOT IN ($2, $3))
			  AND NOT EXISTS (SELECT 1 FROM family_members m WHERE m.family_id = f.id AND m.person_id IN ($2, $3) AND m.role = 'child')
			LIMIT 2
		`, treeID, relationship.Person1ID, relationship.Person2ID)
	} else {
		// A child's existing family identifies which of a parent's potentially
		// multiple families this relationship belongs to.
		familyID, err = uniqueFamily(ctx, tx, `
			SELECT DISTINCT f.id
			FROM families f
			JOIN family_members child ON child.family_id = f.id
			WHERE f.tree_id = $1
			  AND child.person_id = $2
			  AND child.role = 'child'
			LIMIT 2
		`, treeID, relationship.Person2ID)
		if err != nil || familyID != "" {
			return familyID, err
		}

		familyID, err = uniqueFamily(ctx, tx, `
			SELECT DISTINCT f.id
			FROM families f
			JOIN family_members parent ON parent.family_id = f.id
			WHERE f.tree_id = $1
			  AND parent.person_id = $2
			  AND parent.role IN ('partner', 'parent')
			LIMIT 2
		`, treeID, relationship.Person1ID)
	}
	if err != nil || familyID != "" {
		return familyID, err
	}

	familyID = ulid.Make().String()
	if _, err := tx.Exec(ctx, `INSERT INTO families (id, tree_id, auto_created) VALUES ($1, $2, true)`, familyID, treeID); err != nil {
		return "", err
	}
	return familyID, nil
}

// uniqueFamily reads at most two matches: the second match signals that an
// implicit choice would be arbitrary. The caller's transaction rolls back.
func uniqueFamily(ctx context.Context, tx pgx.Tx, query string, args ...any) (string, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		if id != "" {
			return "", coreerrors.ErrAmbiguousFamily
		}
		if err := rows.Scan(&id); err != nil {
			return "", err
		}
	}
	return id, rows.Err()
}

func addFamilyMember(ctx context.Context, tx pgx.Tx, treeID, familyID, personID string, role domain.FamilyMemberRole) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO family_members (tree_id, family_id, person_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING
	`, treeID, familyID, personID, role)
	return err
}
