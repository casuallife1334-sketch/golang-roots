package server

import (
	"genealogy-tree/internal/core/transport/http/middleware"
	"net/http"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []middleware.Middleware
}

func (r Route) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(r.Handler, r.Middleware...)
}
