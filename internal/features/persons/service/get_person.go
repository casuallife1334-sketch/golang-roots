package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) GetPerson(ctx context.Context, userID, treeID, id string) (domain.Person, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return domain.Person{}, err
	}
	return s.personsRepository.GetPerson(ctx, treeID, id)
}
