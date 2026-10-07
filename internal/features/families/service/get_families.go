package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *FamiliesService) GetFamilies(ctx context.Context, userID, treeID string) ([]domain.Family, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	return s.familyRepository.GetFamilies(ctx, treeID)
}

func (s *FamiliesService) GetFamily(ctx context.Context, userID, treeID, id string) (domain.Family, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return domain.Family{}, err
	}
	return s.familyRepository.GetFamily(ctx, treeID, id)
}
