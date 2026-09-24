package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
)

func (s *PersonsService) PatchPerson(ctx context.Context, userID, treeID, id string, input domain.PatchPersonInput) (domain.Person, error) {
	if input.FirstName != nil && strings.TrimSpace(*input.FirstName) == "" || input.LastName != nil && strings.TrimSpace(*input.LastName) == "" {
		return domain.Person{}, ErrInvalid
	}
	if input.BirthDate != nil && input.DeathDate != nil && *input.BirthDate != nil && *input.DeathDate != nil && (*input.DeathDate).Before(**input.BirthDate) {
		return domain.Person{}, ErrInvalid
	}
	if input.Patronymic.Set {
		patronymic, err := normalizePatronymic(input.Patronymic.Value)
		if err != nil {
			return domain.Person{}, err
		}
		input.Patronymic.Value = patronymic
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Person{}, err
	}
	return s.personsRepository.PatchPerson(ctx, treeID, id, input)
}
