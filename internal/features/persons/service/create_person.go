package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
	"time"
)

func (s *PersonsService) CreatePerson(ctx context.Context, input domain.CreatePersonInput) (domain.Person, error) {
	if err := validatePerson(input.FirstName, input.LastName, input.BirthDate, input.DeathDate, input.Gender); err != nil {
		return domain.Person{}, err
	}
	if input.Metadata == nil {
		input.Metadata = map[string]any{}
	}
	return s.personsRepository.CreatePerson(ctx, input)
}

func validatePerson(first, last string, birth, death *time.Time, gender *domain.Gender) error {
	if strings.TrimSpace(first) == "" || strings.TrimSpace(last) == "" {
		return ErrInvalid
	}
	if birth != nil && death != nil && death.Before(*birth) {
		return ErrInvalid
	}
	if gender != nil && *gender != domain.GenderMale && *gender != domain.GenderFemale && *gender != domain.GenderOther {
		return ErrInvalid
	}
	return nil
}
