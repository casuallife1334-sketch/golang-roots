package repository

import "context"

func (r *FamiliesRepository) DeleteFamily(ctx context.Context, treeID, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM families WHERE tree_id = $1 AND id = $2`, treeID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
