package repository

import (
	"context"
	"encoding/json"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r *PersonsRepository) PatchPerson(ctx context.Context, treeID, id string, input domain.PatchPersonInput) (domain.Person, error) {
	metadata, err := json.Marshal(input.Metadata)
	if input.Metadata == nil {
		metadata = nil
	}
	if err != nil {
		return domain.Person{}, err
	}
	var p domain.Person
	var raw []byte
	err = r.db.QueryRow(ctx, `UPDATE persons SET first_name=COALESCE($3,first_name), last_name=COALESCE($4,last_name), birth_date=CASE WHEN $5::boolean THEN $6 ELSE birth_date END, death_date=CASE WHEN $7::boolean THEN $8 ELSE death_date END, gender=CASE WHEN $9::boolean THEN $10 ELSE gender END, photo_url=CASE WHEN $11::boolean THEN $12 ELSE photo_url END, metadata=CASE WHEN $13::boolean THEN $14::jsonb ELSE metadata END, updated_at=now() WHERE tree_id=$1 AND id=$2 RETURNING id,first_name,last_name,birth_date,death_date,gender,photo_url,metadata,created_at,updated_at`, treeID, id, input.FirstName, input.LastName, input.BirthDate != nil, valueTime(input.BirthDate), input.DeathDate != nil, valueTime(input.DeathDate), input.Gender != nil, valueGender(input.Gender), input.PhotoURL != nil, valueString(input.PhotoURL), input.Metadata != nil, metadata).Scan(&p.ID, &p.FirstName, &p.LastName, &p.BirthDate, &p.DeathDate, &p.Gender, &p.PhotoURL, &raw, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Person{}, ErrNotFound
	}
	if err != nil {
		return domain.Person{}, err
	}
	_ = json.Unmarshal(raw, &p.Metadata)
	return p, nil
}

func valueTime(v **time.Time) any {
	if v == nil || *v == nil {
		return nil
	}
	return **v
}
func valueGender(v **domain.Gender) any {
	if v == nil || *v == nil {
		return nil
	}
	return **v
}
func valueString(v **string) any {
	if v == nil || *v == nil {
		return nil
	}
	return **v
}
