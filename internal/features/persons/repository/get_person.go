package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *PersonsRepository) GetPerson(ctx context.Context, id string) (domain.Person, error) {
	return r.scanPerson(r.db.QueryRow(ctx, `SELECT id,first_name,last_name,birth_date,death_date,gender,photo_url,metadata,created_at,updated_at FROM persons WHERE id=$1`, id))
}

func (r *PersonsRepository) scanPerson(row pgx.Row) (domain.Person, error) {
	var p domain.Person
	var raw []byte
	err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.BirthDate, &p.DeathDate, &p.Gender, &p.PhotoURL, &raw, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Person{}, ErrNotFound
	}
	if err != nil {
		return domain.Person{}, fmt.Errorf("scan person: %w", err)
	}
	_ = json.Unmarshal(raw, &p.Metadata)
	return p, nil
}
