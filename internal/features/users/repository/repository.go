package repository

import (
	"fmt"

	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
)

var ErrNotFound = fmt.Errorf("%w: user not found", coreerrors.ErrNotFound)
var ErrDuplicate = fmt.Errorf("%w: email already exists", coreerrors.ErrConflict)

type UsersRepository struct {
	db *postgres.Pool
}

func NewUsersRepository(db *postgres.Pool) *UsersRepository {
	return &UsersRepository{db: db}
}
