package repository

import (
	"context"
	"encoding/json"
	"genealogy-tree/internal/core/domain"
	"github.com/oklog/ulid/v2"
)

func (r *PersonsRepository) CreatePerson(ctx context.Context, input domain.CreatePersonInput) (domain.Person, error) {
	metadata, err := json.Marshal(input.Metadata)
	if err != nil {
		return domain.Person{}, err
	}
	var p domain.Person
	err = r.db.QueryRow(ctx, `INSERT INTO persons (id,first_name,last_name,birth_date,death_date,gender,photo_url,metadata) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,first_name,last_name,birth_date,death_date,gender,photo_url,metadata,created_at,updated_at`, ulid.Make().String(), input.FirstName, input.LastName, input.BirthDate, input.DeathDate, input.Gender, input.PhotoURL, metadata).Scan(&p.ID, &p.FirstName, &p.LastName, &p.BirthDate, &p.DeathDate, &p.Gender, &p.PhotoURL, &metadata, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return domain.Person{}, err
	}
	_ = json.Unmarshal(metadata, &p.Metadata)
	return p, nil
}
