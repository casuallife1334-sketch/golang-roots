package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: family not found", coreerrors.ErrNotFound)
var ErrDuplicate = fmt.Errorf("%w: family member or relationship already exists", coreerrors.ErrConflict)
var ErrPersonNotFound = fmt.Errorf("%w: person not found", coreerrors.ErrNotFound)

type FamiliesRepository struct {
	db *postgres.Pool
}

func NewFamiliesRepository(db *postgres.Pool) *FamiliesRepository {
	return &FamiliesRepository{db: db}
}
