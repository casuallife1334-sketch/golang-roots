package repository

import (
	"context"
	"encoding/json"
	"genealogy-tree/internal/core/domain"
	"github.com/oklog/ulid/v2"
)

func (r *FamiliesRepository) CreateFamily(ctx context.Context, treeID string, input domain.CreateFamilyInput) (domain.Family, error) {
	metadata, err := json.Marshal(input.Metadata)
	if err != nil {
		return domain.Family{}, err
	}
	var family domain.Family
	err = r.db.QueryRow(ctx, `
		INSERT INTO families (id, tree_id, name, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tree_id, name, metadata, created_at, updated_at
	`, ulid.Make().String(), treeID, input.Name, metadata).Scan(
		&family.ID,
		&family.TreeID,
		&family.Name,
		&metadata,
		&family.CreatedAt,
		&family.UpdatedAt,
	)
	if err != nil {
		return domain.Family{}, err
	}
	if err := json.Unmarshal(metadata, &family.Metadata); err != nil {
		return domain.Family{}, err
	}
	family.Members = []domain.FamilyMember{}
	family.RelationshipIDs = []string{}
	return family, nil
}
