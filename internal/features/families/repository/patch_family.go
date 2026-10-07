package repository

import (
	"context"
	"encoding/json"
	"genealogy-tree/internal/core/domain"
)

func (r *FamiliesRepository) PatchFamily(ctx context.Context, treeID, id string, input domain.PatchFamilyInput) (domain.Family, error) {
	var metadata []byte
	var err error
	if input.Metadata != nil {
		metadata, err = json.Marshal(*input.Metadata)
		if err != nil {
			return domain.Family{}, err
		}
	}
	_, err = r.db.Exec(ctx, `
		UPDATE families
		SET name = CASE WHEN $3::boolean THEN $4 ELSE name END,
			metadata = CASE WHEN $5::boolean THEN $6::jsonb ELSE metadata END,
			auto_created = CASE WHEN $3::boolean OR $5::boolean THEN false ELSE auto_created END,
			updated_at = CASE WHEN $3::boolean OR $5::boolean THEN now() ELSE updated_at END
		WHERE tree_id = $1 AND id = $2
	`, treeID, id, input.Name != nil, valueOrEmpty(input.Name), input.Metadata != nil, metadata)
	if err != nil {
		return domain.Family{}, err
	}
	return r.GetFamily(ctx, treeID, id)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
