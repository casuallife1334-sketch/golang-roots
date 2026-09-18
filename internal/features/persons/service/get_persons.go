package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) GetPersons(ctx context.Context, userID, treeID string) ([]domain.Person, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	return s.personsRepository.GetPersons(ctx, treeID)
}
