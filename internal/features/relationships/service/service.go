package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrInvalid = fmt.Errorf("%w: invalid relationship input", coreerrors.ErrInvalidArgument)

type RelationshipsRepository interface {
	CreateRelationship(context.Context, string, domain.CreateRelationshipInput) (domain.Relationship, error)
	GetRelationship(context.Context, string, string) (domain.Relationship, error)
	GetRelationships(context.Context, string, string) ([]domain.Relationship, error)
	PatchRelationship(context.Context, string, string, domain.PatchRelationshipInput) (domain.Relationship, error)
}

type RelationshipDeleter interface {
	DeleteRelationship(context.Context, string, string) error
}

type TreeAccess interface {
	CanReadTree(context.Context, string, string) error
	CanWriteTree(context.Context, string, string) error
}

type RelationshipsService struct {
	relationshipsRepository RelationshipsRepository
	relationshipDeleter     RelationshipDeleter
	treeAccess              TreeAccess
}

func NewRelationshipsService(relationshipsRepository RelationshipsRepository, relationshipDeleter RelationshipDeleter, treeAccess TreeAccess) *RelationshipsService {
	return &RelationshipsService{relationshipsRepository: relationshipsRepository, relationshipDeleter: relationshipDeleter, treeAccess: treeAccess}
}
