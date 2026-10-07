package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/validation"
)

func (s *FamiliesService) AddMember(ctx context.Context, userID, treeID, familyID string, input domain.AddFamilyMemberInput) error {
	if err := validation.ValidateULID(input.PersonID); err != nil || !validFamilyMemberRole(input.Role) {
		return ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.membershipRepository.AddMember(ctx, treeID, familyID, input)
}

func (s *FamiliesService) RemoveMember(ctx context.Context, userID, treeID, familyID, personID string) error {
	if err := validation.ValidateULID(personID); err != nil {
		return ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return err
	}
	return s.membershipRepository.RemoveMember(ctx, treeID, familyID, personID)
}
