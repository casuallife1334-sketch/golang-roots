package repository

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (r *TreesRepository) GetTrees(ctx context.Context, userID string) ([]domain.Tree, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id, t.owner_id, tm.role, t.name, t.created_at, t.updated_at
		FROM trees t JOIN tree_members tm ON tm.tree_id = t.id
		WHERE tm.user_id = $1 ORDER BY t.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	trees := []domain.Tree{}
	for rows.Next() {
		var tree domain.Tree
		if err := rows.Scan(&tree.ID, &tree.OwnerID, &tree.Role, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt); err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}
	return trees, rows.Err()
}
