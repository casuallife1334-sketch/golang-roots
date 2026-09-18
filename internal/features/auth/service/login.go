package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/security"
)

type LoginResult struct {
	User        domain.User
	AccessToken string
	ExpiresIn   int64
}

func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return LoginResult{}, ErrUnauthorized
	}
	user, passwordHash, err := s.usersRepository.GetUserByEmail(ctx, normalizedEmail)
	if err != nil || security.ComparePassword(passwordHash, password) != nil {
		return LoginResult{}, ErrUnauthorized
	}
	token, err := s.tokens.CreateAccessToken(user.ID)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{User: user, AccessToken: token, ExpiresIn: int64(s.tokens.TTL().Seconds())}, nil
}
