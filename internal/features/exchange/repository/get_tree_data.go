package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"genealogy-tree/internal/core/domain"
)

func (r *ExchangeRepository) GetTreeData(ctx context.Context, treeID string) ([]domain.Person, []domain.Relationship, error) {
	persons, err := r.getPersons(ctx, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("get tree persons: %w", err)
	}
	relationships, err := r.getRelationships(ctx, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("get tree relationships: %w", err)
	}
	return persons, relationships, nil
}

func (r *ExchangeRepository) getPersons(ctx context.Context, treeID string) ([]domain.Person, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, first_name, patronymic, last_name, birth_date, death_date,
		       gender, photo_url, metadata, created_at, updated_at
		FROM persons
		WHERE tree_id = $1
		ORDER BY id
	`, treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	persons := []domain.Person{}
	for rows.Next() {
		var person domain.Person
		var metadata []byte
		if err := rows.Scan(
			&person.ID,
			&person.FirstName,
			&person.Patronymic,
			&person.LastName,
			&person.BirthDate,
			&person.DeathDate,
			&person.Gender,
			&person.PhotoURL,
			&metadata,
			&person.CreatedAt,
			&person.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadata, &person.Metadata); err != nil {
			return nil, fmt.Errorf("decode person metadata: %w", err)
		}
		persons = append(persons, person)
	}
	return persons, rows.Err()
}

func (r *ExchangeRepository) getRelationships(ctx context.Context, treeID string) ([]domain.Relationship, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, person1_id, person2_id, type, direction, metadata, created_at, updated_at
		FROM relationships
		WHERE tree_id = $1
		ORDER BY id
	`, treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	relationships := []domain.Relationship{}
	for rows.Next() {
		var relationship domain.Relationship
		var metadata []byte
		if err := rows.Scan(
			&relationship.ID,
			&relationship.Person1ID,
			&relationship.Person2ID,
			&relationship.Type,
			&relationship.Direction,
			&metadata,
			&relationship.CreatedAt,
			&relationship.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadata, &relationship.Metadata); err != nil {
			return nil, fmt.Errorf("decode relationship metadata: %w", err)
		}
		relationships = append(relationships, relationship)
	}
	return relationships, rows.Err()
}
