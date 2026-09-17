package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

type UsersRepository interface {
	GetUser(context.Context, string) (domain.User, error)
}

type UsersService struct {
	usersRepository UsersRepository
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{usersRepository: usersRepository}
}

func (s *UsersService) GetUser(ctx context.Context, id string) (domain.User, error) {
	return s.usersRepository.GetUser(ctx, id)
}
