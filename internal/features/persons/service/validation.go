package service

import (
	"genealogy-tree/internal/core/domain"
	"strings"
	"time"
	"unicode/utf8"
)

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
	if err := validatePersonDates(birth, death); err != nil {
		return ErrInvalid
	}
	if gender != nil && *gender != domain.GenderMale && *gender != domain.GenderFemale && *gender != domain.GenderOther {
		return ErrInvalid
	}
	return nil
}

func validatePersonDates(birth, death *time.Time) error {
	if birth != nil && death != nil && death.Before(*birth) {
		return ErrInvalid
	}
	return nil
}
