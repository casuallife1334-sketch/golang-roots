package repository

import "context"

func (r *TreesRepository) DeleteTree(ctx context.Context, ownerID, treeID string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM trees WHERE id = $1 AND owner_id = $2`, treeID, ownerID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
