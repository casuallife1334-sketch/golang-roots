package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: relationship not found", coreerrors.ErrNotFound)
var ErrDuplicate = fmt.Errorf("%w: relationship already exists", coreerrors.ErrConflict)
var ErrPersonNotFound = fmt.Errorf("%w: person not found", coreerrors.ErrNotFound)

type RelationshipsRepository struct {
	db *postgres.Pool
}

func NewRelationshipsRepository(db *postgres.Pool) *RelationshipsRepository {
	return &RelationshipsRepository{db: db}
}
