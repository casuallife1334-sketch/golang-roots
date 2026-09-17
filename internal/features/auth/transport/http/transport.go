package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"genealogy-tree/internal/features/auth/service"
	"net/http"
)

type AuthService interface {
	Register(context.Context, string, string) (domain.User, error)
	Login(context.Context, string, string) (service.LoginResult, error)
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: authService}
}

func (h *AuthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/auth/register", Handler: h.Register},
		{Method: http.MethodPost, Path: "/auth/login", Handler: h.Login},
	}
}
