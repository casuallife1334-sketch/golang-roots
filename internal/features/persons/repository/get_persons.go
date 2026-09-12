package repository

import (
	"context"
	"encoding/json"
	"genealogy-tree/internal/core/domain"
)

func (r *PersonsRepository) GetPersons(ctx context.Context) ([]domain.Person, error) {
	rows, err := r.db.Query(ctx, `SELECT id,first_name,last_name,birth_date,death_date,gender,photo_url,metadata,created_at,updated_at FROM persons ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Person{}
	for rows.Next() {
		var p domain.Person
		var raw []byte
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.BirthDate, &p.DeathDate, &p.Gender, &p.PhotoURL, &raw, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(raw, &p.Metadata)
		items = append(items, p)
	}
	return items, rows.Err()
}
