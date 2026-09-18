package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: tree not found", coreerrors.ErrNotFound)
var ErrAccessDenied = fmt.Errorf("%w: tree access denied", coreerrors.ErrForbidden)

type TreesRepository struct {
	db *postgres.Pool
}

func NewTreesRepository(db *postgres.Pool) *TreesRepository {
	return &TreesRepository{db: db}
}
