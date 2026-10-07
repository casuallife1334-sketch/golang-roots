package service

import (
	"context"
	"genealogy-tree/internal/core/validation"
)

func (s *FamiliesService) AttachRelationship(ctx context.Context, userID, treeID, familyID, relationshipID string) error {
	if err := validation.ValidateULID(relationshipID); err != nil {
		return ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.relationshipRepository.AttachRelationship(ctx, treeID, familyID, relationshipID)
}

func (s *FamiliesService) DetachRelationship(ctx context.Context, userID, treeID, familyID, relationshipID string) error {
	if err := validation.ValidateULID(relationshipID); err != nil {
		return ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.relationshipRepository.DetachRelationship(ctx, treeID, familyID, relationshipID)
}
