package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
	"time"
	"unicode/utf8"
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

func normalizePatronymic(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" || utf8.RuneCountInString(trimmed) > 200 {
		return nil, ErrInvalid
	}
	return &trimmed, nil
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
