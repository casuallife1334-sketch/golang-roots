package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: person not found", coreerrors.ErrNotFound)

type PersonsRepository struct {
	db *postgres.Pool
}

func NewPersonsRepository(db *postgres.Pool) *PersonsRepository {
	return &PersonsRepository{db: db}
}
