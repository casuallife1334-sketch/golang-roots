package repository

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type RelationshipRepository interface {
	CreateRelationshipTx(context.Context, pgx.Tx, string, domain.CreateRelationshipInput) (domain.Relationship, error)
}

type FamilyRepository interface {
	EnsureRelationshipTx(context.Context, pgx.Tx, string, domain.Relationship, *string) error
}

func (r *RelationshipFamilyRepository) CreateRelationshipWithFamily(
	ctx context.Context,
	treeID string,
	input domain.CreateRelationshipInput,
	familyID *string,
) (domain.Relationship, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Relationship{}, err
	}
	defer tx.Rollback(ctx)

	// Serialize relationship/family creation within a tree. Otherwise concurrent
	// requests can both see no matching family and create two for the same child.
	if _, err := tx.Exec(ctx, `
		SELECT pg_advisory_xact_lock(hashtextextended('relationship-family:' || $1::text, 0))
	`, treeID); err != nil {
		return domain.Relationship{}, err
	}

	relationship, err := r.relationships.CreateRelationshipTx(ctx, tx, treeID, input)
	if err != nil {
		return domain.Relationship{}, err
	}
	if err := r.families.EnsureRelationshipTx(ctx, tx, treeID, relationship, familyID); err != nil {
		return domain.Relationship{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Relationship{}, err
	}
	return relationship, nil
}
