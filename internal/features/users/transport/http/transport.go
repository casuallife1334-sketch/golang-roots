package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"net/http"
)

type UsersService interface {
	GetUser(context.Context, string) (domain.User, error)
}

type UsersHTTPHandler struct {
	usersService UsersService
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{usersService: usersService}
}

func (h *UsersHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodGet, Path: "/users/me", Handler: h.GetMe},
	}
}
