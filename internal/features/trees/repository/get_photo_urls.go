package repository

import "context"

func (r *TreesRepository) GetPhotoURLs(ctx context.Context, treeID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `SELECT photo_url FROM persons WHERE tree_id = $1 AND photo_url IS NOT NULL`, treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}
