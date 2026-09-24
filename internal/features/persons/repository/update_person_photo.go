package repository

import (
	"context"
	"encoding/json"
	"errors"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *PersonsRepository) UpdatePersonPhoto(ctx context.Context, treeID, id string, photoURL *string) (domain.Person, error) {
	var person domain.Person
	var metadata []byte
	err := r.db.QueryRow(ctx, `UPDATE persons SET photo_url = $3, updated_at = now() WHERE tree_id = $1 AND id = $2 RETURNING id,first_name,patronymic,last_name,birth_date,death_date,gender,photo_url,metadata,created_at,updated_at`, treeID, id, photoURL).Scan(&person.ID, &person.FirstName, &person.Patronymic, &person.LastName, &person.BirthDate, &person.DeathDate, &person.Gender, &person.PhotoURL, &metadata, &person.CreatedAt, &person.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Person{}, ErrNotFound
	}
	if err != nil {
		return domain.Person{}, err
	}
	_ = json.Unmarshal(metadata, &person.Metadata)
	return person, nil
}
