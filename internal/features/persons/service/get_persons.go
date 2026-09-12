package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) GetPersons(ctx context.Context) ([]domain.Person, error) {
	return s.personsRepository.GetPersons(ctx)
}
