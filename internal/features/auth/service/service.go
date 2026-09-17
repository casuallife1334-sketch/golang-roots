package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
)

var (
	ErrInvalid      = fmt.Errorf("%w: invalid auth input", coreerrors.ErrInvalidArgument)
	ErrUnauthorized = coreerrors.ErrUnauthorized
)

type UsersRepository interface {
	CreateUser(context.Context, domain.CreateUserInput) (domain.User, error)
	GetUserByEmail(context.Context, string) (domain.User, string, error)
}

type AuthService struct {
	usersRepository UsersRepository
	tokens          *security.TokenManager
}

func NewAuthService(usersRepository UsersRepository, tokens *security.TokenManager) *AuthService {
	return &AuthService{usersRepository: usersRepository, tokens: tokens}
}
