package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *FamiliesService) CreateFamily(ctx context.Context, userID, treeID string, input domain.CreateFamilyInput) (domain.Family, error) {
	var err error
	input, err = normalizeFamilyInput(input)
	if err != nil {
		return domain.Family{}, ErrInvalid
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Family{}, err
	}
	return s.familyRepository.CreateFamily(ctx, treeID, input)
}
