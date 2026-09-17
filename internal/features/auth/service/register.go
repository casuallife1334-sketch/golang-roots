package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/security"
	"net/mail"
	"strings"
)

func (s *AuthService) Register(ctx context.Context, email, password string) (domain.User, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	if len(password) < 8 || len(password) > 72 {
		return domain.User{}, ErrInvalid
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}
	return s.usersRepository.CreateUser(ctx, domain.CreateUserInput{Email: email, PasswordHash: hash})
}

func normalizeEmail(value string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", ErrInvalid
	}
	return email, nil
}
