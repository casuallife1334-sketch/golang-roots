package repository

import (
	"context"
	"encoding/json"
	"genealogy-tree/internal/core/domain"
)

func (r *RelationshipsRepository) PatchRelationship(ctx context.Context, treeID, id string, input domain.PatchRelationshipInput) (domain.Relationship, error) {
	var metadata []byte
	var err error
	if input.Metadata.Set {
		metadata, err = json.Marshal(input.Metadata.Value)
		if err != nil {
			return domain.Relationship{}, err
		}
	}
	return r.scanRelationship(r.db.QueryRow(ctx, `
		UPDATE relationships
		SET metadata = CASE WHEN $3::boolean THEN $4::jsonb ELSE metadata END,
			updated_at = CASE WHEN $3::boolean THEN now() ELSE updated_at END
		WHERE tree_id = $1 AND id = $2
		RETURNING id,person1_id,person2_id,type,direction,metadata,created_at,updated_at
	`, treeID, id, input.Metadata.Set, metadata))
}
