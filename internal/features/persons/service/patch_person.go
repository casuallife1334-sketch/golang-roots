package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
)

func (s *PersonsService) PatchPerson(ctx context.Context, id string, input domain.PatchPersonInput) (domain.Person, error) {
	if input.FirstName != nil && strings.TrimSpace(*input.FirstName) == "" || input.LastName != nil && strings.TrimSpace(*input.LastName) == "" {
		return domain.Person{}, ErrInvalid
	}
	if input.BirthDate != nil && input.DeathDate != nil && *input.BirthDate != nil && *input.DeathDate != nil && (*input.DeathDate).Before(**input.BirthDate) {
		return domain.Person{}, ErrInvalid
	}
	return s.personsRepository.PatchPerson(ctx, id, input)
}
