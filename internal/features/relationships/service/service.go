package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrInvalid = fmt.Errorf("%w: invalid relationship input", coreerrors.ErrInvalidArgument)

type RelationshipsRepository interface {
	CreateRelationship(context.Context, domain.CreateRelationshipInput) (domain.Relationship, error)
	GetRelationship(context.Context, string) (domain.Relationship, error)
	DeleteRelationship(context.Context, string) error
}

type RelationshipsService struct {
	relationshipsRepository RelationshipsRepository
}

func NewRelationshipsService(relationshipsRepository RelationshipsRepository) *RelationshipsService {
	return &RelationshipsService{relationshipsRepository: relationshipsRepository}
}
