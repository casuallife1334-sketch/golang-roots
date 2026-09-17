package middleware

import (
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	"genealogy-tree/internal/core/transport/http/response"
	"net/http"
	"strings"
)

func Auth(tokenManager *security.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerParts := strings.Fields(r.Header.Get("Authorization"))
			if len(headerParts) != 2 || !strings.EqualFold(headerParts[0], "Bearer") {
				response.NewHTTPResponseHandler(w).ErrorResponse(coreerrors.ErrUnauthorized, "authentication is required")
				return
			}

			userID, err := tokenManager.ParseAccessToken(headerParts[1])
			if err != nil {
				response.NewHTTPResponseHandler(w).ErrorResponse(coreerrors.ErrUnauthorized, "access token is invalid")
				return
			}
			next.ServeHTTP(w, r.WithContext(security.WithUserID(r.Context(), userID)))
		})
	}
}
