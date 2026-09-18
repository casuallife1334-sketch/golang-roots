package repository

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"github.com/oklog/ulid/v2"
)

func (r *TreesRepository) CreateTree(ctx context.Context, ownerID string, input domain.CreateTreeInput) (domain.Tree, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Tree{}, err
	}
	defer tx.Rollback(ctx)

	id := ulid.Make().String()
	var tree domain.Tree
	err = tx.QueryRow(ctx, `
		INSERT INTO trees (id, owner_id, name) VALUES ($1, $2, $3)
		RETURNING id, owner_id, name, created_at, updated_at
	`, id, ownerID, input.Name).Scan(&tree.ID, &tree.OwnerID, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt)
	if err != nil {
		return domain.Tree{}, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO tree_members (tree_id, user_id, role) VALUES ($1, $2, $3)`, id, ownerID, domain.TreeRoleOwner); err != nil {
		return domain.Tree{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Tree{}, err
	}
	return tree, nil
}
