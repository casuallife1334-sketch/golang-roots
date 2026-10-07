package repository

import (
	"context"
	"genealogy-tree/internal/core/repository/postgres"

	"github.com/jackc/pgx/v5"
)

type RelationshipRepository interface {
	DeleteRelationshipTx(context.Context, pgx.Tx, string, string) error
}

type PersonRepository interface {
	DeletePersonTx(context.Context, pgx.Tx, string, string) error
}

type FamilyRepository interface {
	ReconcileTx(context.Context, pgx.Tx, string, []string) error
}

type Repository struct {
	db            *postgres.Pool
	relationships RelationshipRepository
	persons       PersonRepository
	families      FamilyRepository
}

func NewRepository(db *postgres.Pool, relationships RelationshipRepository, persons PersonRepository, families FamilyRepository) *Repository {
	return &Repository{db: db, relationships: relationships, persons: persons, families: families}
}
