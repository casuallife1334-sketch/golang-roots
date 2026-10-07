package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *RelationshipsRepository) DeleteRelationshipTx(ctx context.Context, tx pgx.Tx, treeID, id string) error {
	result, err := tx.Exec(ctx, `DELETE FROM relationships WHERE tree_id = $1 AND id = $2`, treeID, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
