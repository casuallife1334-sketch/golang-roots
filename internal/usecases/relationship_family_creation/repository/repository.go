package repository

import (
	"genealogy-tree/internal/core/repository/postgres"
)

type RelationshipFamilyRepository struct {
	db            *postgres.Pool
	relationships RelationshipRepository
	families      FamilyRepository
}

func NewRelationshipFamilyRepository(
	db *postgres.Pool,
	relationships RelationshipRepository,
	families FamilyRepository,
) *RelationshipFamilyRepository {
	return &RelationshipFamilyRepository{
		db:            db,
		relationships: relationships,
		families:      families,
	}
}
