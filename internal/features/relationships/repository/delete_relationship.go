package repository

import "context"

func (r *RelationshipsRepository) DeleteRelationship(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM relationships WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
