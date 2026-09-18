package repository

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *TreesRepository) GetTree(ctx context.Context, userID, treeID string) (domain.Tree, error) {
	var tree domain.Tree
	err := r.db.QueryRow(ctx, `
		SELECT t.id, t.owner_id, t.name, t.created_at, t.updated_at
		FROM trees t JOIN tree_members tm ON tm.tree_id = t.id
		WHERE t.id = $1 AND tm.user_id = $2
	`, treeID, userID).Scan(&tree.ID, &tree.OwnerID, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Tree{}, ErrNotFound
	}
	return tree, err
}

func (r *TreesRepository) GetMemberRole(ctx context.Context, userID, treeID string) (domain.TreeRole, error) {
	var role domain.TreeRole
	err := r.db.QueryRow(ctx, `SELECT role FROM tree_members WHERE tree_id = $1 AND user_id = $2`, treeID, userID).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrAccessDenied
	}
	return role, err
}
