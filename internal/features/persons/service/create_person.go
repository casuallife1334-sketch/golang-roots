package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *PersonsService) CreatePerson(ctx context.Context, userID, treeID string, input domain.CreatePersonInput) (domain.Person, error) {
	if err := validatePerson(input.FirstName, input.LastName, input.BirthDate, input.DeathDate, input.Gender); err != nil {
		return domain.Person{}, err
	}
	if input.Metadata == nil {
		input.Metadata = map[string]any{}
	}
	patronymic, err := normalizePatronymic(input.Patronymic)
	if err != nil {
		return domain.Person{}, err
	}
	input.Patronymic = patronymic
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Person{}, err
	}
	return s.personsRepository.CreatePerson(ctx, treeID, input)
}
