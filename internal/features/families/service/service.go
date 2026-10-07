package service

import (
	"context"
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrInvalid = fmt.Errorf("%w: invalid family input", coreerrors.ErrInvalidArgument)

type TreeAccess interface {
	CanReadTree(context.Context, string, string) error
	CanWriteTree(context.Context, string, string) error
}

type FamiliesService struct {
	familyRepository       FamilyRepository
	membershipRepository   FamilyMembershipRepository
	relationshipRepository FamilyRelationshipRepository
	treeAccess             TreeAccess
}

func NewFamiliesService(
	familyRepository FamilyRepository,
	membershipRepository FamilyMembershipRepository,
	relationshipRepository FamilyRelationshipRepository,
	treeAccess TreeAccess,
) *FamiliesService {
	return &FamiliesService{
		familyRepository:       familyRepository,
		membershipRepository:   membershipRepository,
		relationshipRepository: relationshipRepository,
		treeAccess:             treeAccess,
	}
}
