package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/features/exchange/gedcom"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func (r *ExchangeRepository) ImportSnapshot(ctx context.Context, treeID string, snapshot gedcom.Snapshot) (int, int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	idByExternalID := make(map[string]string, len(snapshot.Persons))
	for _, person := range snapshot.Persons {
		id := ulid.Make().String()
		metadata, err := json.Marshal(person.Input.Metadata)
		if err != nil {
			return 0, 0, fmt.Errorf("encode person metadata: %w", err)
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO persons (
				id, tree_id, first_name, patronymic, last_name, birth_date,
				death_date, gender, photo_url, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, id, treeID, person.Input.FirstName, person.Input.Patronymic, person.Input.LastName,
			person.Input.BirthDate, person.Input.DeathDate, person.Input.Gender,
			person.Input.PhotoURL, metadata)
		if err != nil {
			return 0, 0, fmt.Errorf("insert imported person: %w", err)
		}
		idByExternalID[person.ExternalID] = id
	}

	for _, relationship := range snapshot.Relationships {
		person1ID, ok1 := idByExternalID[relationship.Person1ExternalID]
		person2ID, ok2 := idByExternalID[relationship.Person2ExternalID]
		if !ok1 || !ok2 {
			return 0, 0, fmt.Errorf("%w: relationship references an unknown imported person", ErrInvalidImport)
		}
		metadata, err := json.Marshal(relationship.Input.Metadata)
		if err != nil {
			return 0, 0, fmt.Errorf("encode relationship metadata: %w", err)
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO relationships (
				id, tree_id, person1_id, person2_id, type, direction, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, ulid.Make().String(), treeID, person1ID, person2ID, relationship.Input.Type,
			relationship.Input.Direction, metadata)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return 0, 0, fmt.Errorf("%w: duplicate relationship in imported data", coreerrors.ErrConflict)
			}
			return 0, 0, fmt.Errorf("insert imported relationship: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, err
	}
	return len(snapshot.Persons), len(snapshot.Relationships), nil
}
