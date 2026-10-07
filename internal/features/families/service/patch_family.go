package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *FamiliesService) PatchFamily(ctx context.Context, userID, treeID, id string, input domain.PatchFamilyInput) (domain.Family, error) {
	if input.Name != nil {
		value := *input.Name
		input.Name = &value
	}
	if input.Metadata != nil && *input.Metadata == nil {
		return domain.Family{}, ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Family{}, err
	}
	return s.familyRepository.PatchFamily(ctx, treeID, id, input)
}
