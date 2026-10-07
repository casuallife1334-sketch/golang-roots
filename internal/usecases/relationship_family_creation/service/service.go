package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

type Repository interface {
	CreateRelationshipWithFamily(context.Context, string, domain.CreateRelationshipInput, *string) (domain.Relationship, error)
}

type TreeAccess interface {
	CanWriteTree(context.Context, string, string) error
}

type RelationshipFamilyCreationService struct {
	repository Repository
	treeAccess TreeAccess
}

func NewRelationshipFamilyCreationService(repository Repository, treeAccess TreeAccess) *RelationshipFamilyCreationService {
	return &RelationshipFamilyCreationService{repository: repository, treeAccess: treeAccess}
}
