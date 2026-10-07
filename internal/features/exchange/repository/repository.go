package repository

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrInvalidImport = fmt.Errorf("%w: invalid GEDCOM import data", coreerrors.ErrInvalidArgument)

type ExchangeRepository struct {
	db *postgres.Pool
}

func NewExchangeRepository(db *postgres.Pool) *ExchangeRepository {
	return &ExchangeRepository{db: db}
}
