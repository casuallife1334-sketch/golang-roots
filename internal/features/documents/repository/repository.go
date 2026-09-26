package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: document not found", coreerrors.ErrNotFound)

type DocumentsRepository struct {
	db *postgres.Pool
}

func NewDocumentsRepository(db *postgres.Pool) *DocumentsRepository {
	return &DocumentsRepository{db: db}
}
