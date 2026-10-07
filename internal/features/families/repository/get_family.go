package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

const familySelect = `
	SELECT
		f.id,
		f.tree_id,
		f.name,
		f.metadata,
		f.created_at,
		f.updated_at,
		COALESCE(
			jsonb_agg(DISTINCT jsonb_build_object(
				'person_id', fm.person_id::text,
				'role', fm.role::text
			)) FILTER (WHERE fm.person_id IS NOT NULL),
			'[]'::jsonb
		) AS members,
		COALESCE(
			array_agg(DISTINCT fr.relationship_id::text) FILTER (WHERE fr.relationship_id IS NOT NULL),
			ARRAY[]::text[]
		) AS relationship_ids
	FROM families f
	LEFT JOIN family_members fm ON fm.family_id = f.id
	LEFT JOIN family_relationships fr ON fr.family_id = f.id
`

func (r *FamiliesRepository) GetFamily(ctx context.Context, treeID, id string) (domain.Family, error) {
	return scanFamily(r.db.QueryRow(ctx, familySelect+`WHERE f.tree_id = $1 AND f.id = $2 GROUP BY f.id`, treeID, id))
}

func (r *FamiliesRepository) GetFamilies(ctx context.Context, treeID string) ([]domain.Family, error) {
	rows, err := r.db.Query(ctx, familySelect+`WHERE f.tree_id = $1 GROUP BY f.id ORDER BY f.created_at, f.id`, treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []domain.Family{}
	for rows.Next() {
		family, err := scanFamily(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, family)
	}
	return items, rows.Err()
}

type familyRow interface {
	Scan(dest ...any) error
}

func scanFamily(row familyRow) (domain.Family, error) {
	var family domain.Family
	var metadata, members []byte
	var relationshipIDs []string
	if err := row.Scan(
		&family.ID,
		&family.TreeID,
		&family.Name,
		&metadata,
		&family.CreatedAt,
		&family.UpdatedAt,
		&members,
		&relationshipIDs,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Family{}, ErrNotFound
		}
		return domain.Family{}, fmt.Errorf("scan family: %w", err)
	}
	if err := json.Unmarshal(metadata, &family.Metadata); err != nil {
		return domain.Family{}, fmt.Errorf("decode family metadata: %w", err)
	}
	if err := json.Unmarshal(members, &family.Members); err != nil {
		return domain.Family{}, fmt.Errorf("decode family members: %w", err)
	}
	family.RelationshipIDs = relationshipIDs
	return family, nil
}
