package server

import (
	"context"
	"errors"
	"fmt"
	corelogger "genealogy-tree/internal/core/logger"
	"genealogy-tree/internal/core/transport/http/middleware"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	middleware []middleware.Middleware
	log        *corelogger.Logger
}

func NewHTTPServer(config Config, log *corelogger.Logger, middlewares ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{mux: http.NewServeMux(), config: config, log: log, middleware: middlewares}
}
func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}
func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		s.mux.Handle(fmt.Sprintf("%s %s", route.Method, route.Path), route.WithMiddleware())
	}
}
func (s *HTTPServer) Run(ctx context.Context) error {
	httpServer := &http.Server{Addr: s.config.Addr, Handler: middleware.ChainMiddleware(s.mux, s.middleware...)}
	errorsChannel := make(chan error, 1)
	go func() {
		s.log.Debug("start HTTP server", zap.String("addr", s.config.Addr))
		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errorsChannel <- err
		}
	}()
	select {
	case err := <-errorsChannel:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		s.log.Debug("HTTP server stopped")
	}
	return nil
}
