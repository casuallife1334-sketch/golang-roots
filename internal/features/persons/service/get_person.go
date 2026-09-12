package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) GetPerson(ctx context.Context, id string) (domain.Person, error) {
	return s.personsRepository.GetPerson(ctx, id)
}
