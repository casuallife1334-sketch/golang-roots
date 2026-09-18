package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *TreesRepository) PatchTree(ctx context.Context, ownerID, treeID string, input domain.PatchTreeInput) (domain.Tree, error) {
	var tree domain.Tree
	err := r.db.QueryRow(ctx, `
		UPDATE trees SET name = COALESCE($3, name), updated_at = now()
		WHERE id = $1 AND owner_id = $2
		RETURNING id, owner_id, name, created_at, updated_at
	`, treeID, ownerID, input.Name).Scan(&tree.ID, &tree.OwnerID, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Tree{}, ErrNotFound
	}
	return tree, err
}
