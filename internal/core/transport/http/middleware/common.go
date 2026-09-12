package middleware

import (
	corelogger "genealogy-tree/internal/core/logger"
	"genealogy-tree/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func CORS(allowedOrigins []string) Middleware {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				}
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = uuid.NewString()
			}
			r.Header.Set("X-Request-ID", id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *corelogger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLog := log.With(zap.String("request_id", r.Header.Get("X-Request-ID")), zap.String("url", r.URL.String()))
			next.ServeHTTP(w, r.WithContext(corelogger.ToContext(r.Context(), requestLog)))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := response.NewResponseWriter(w)
			started := time.Now()
			log := corelogger.FromContext(r.Context())
			log.Debug(">>> incoming HTTP request", zap.String("http_method", r.Method), zap.Time("time", started.UTC()))
			next.ServeHTTP(wrapped, r)
			log.Debug("<<< done HTTP request", zap.Int("status_code", wrapped.GetStatusCode()), zap.Duration("latency", time.Since(started)))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					corelogger.FromContext(r.Context()).Error("unexpected panic while handling HTTP request", zap.Any("panic", value))
					response.NewHTTPResponseHandler(w).PanicResponse(value, "during handle HTTP request got unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
