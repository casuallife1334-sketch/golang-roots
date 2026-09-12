package server

import (
	"fmt"
	"genealogy-tree/internal/core/transport/http/middleware"
	"net/http"
)

type ApiVersion string

const ApiVersion1 ApiVersion = "v1"

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []middleware.Middleware
}

func NewApiVersionRouter(version ApiVersion, middlewares ...middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{ServeMux: http.NewServeMux(), apiVersion: version, middleware: middlewares}
}
func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		r.Handle(fmt.Sprintf("%s %s", route.Method, route.Path), route.WithMiddleware())
	}
}
func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(r, r.middleware...)
}
