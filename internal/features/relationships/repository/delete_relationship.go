package repository

import "context"

func (r *RelationshipsRepository) DeleteRelationship(ctx context.Context, treeID, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM relationships WHERE tree_id=$1 AND id=$2`, treeID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
